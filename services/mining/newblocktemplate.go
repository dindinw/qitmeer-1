package mining

import (

	"fmt"
	"encoding/binary"
	"container/heap"
	"github.com/HalalChain/qitmeer-lib/core/protocol"
	"github.com/HalalChain/qitmeer/core/merkle"
	"github.com/HalalChain/qitmeer-lib/config"
	"github.com/HalalChain/qitmeer-lib/params"
	"github.com/HalalChain/qitmeer-lib/common/hash"
	"github.com/HalalChain/qitmeer/core/blockchain"
	"github.com/HalalChain/qitmeer-lib/engine/txscript"
	"github.com/HalalChain/qitmeer/services/blkmgr"
	"github.com/HalalChain/qitmeer/services/mempool"
	"github.com/HalalChain/qitmeer-lib/log"
	"github.com/HalalChain/qitmeer-lib/core/types"
	"github.com/HalalChain/qitmeer-lib/core/address"
	s "github.com/HalalChain/qitmeer-lib/core/serialization"
)


// NewBlockTemplate returns a new block template that is ready to be solved
// using the transactions from the passed transaction source pool and a coinbase
// that either pays to the passed address if it is not nil, or a coinbase that
// is redeemable by anyone if the passed address is nil.  The nil address
// functionality is useful since there are cases such as the getblocktemplate
// RPC where external mining software is responsible for creating their own
// coinbase which will replace the one generated for the block template.  Thus
// the need to have configured address can be avoided.
//
// The transactions selected and included are prioritized according to several
// factors.  First, each transaction has a priority calculated based on its
// value, age of inputs, and size.  Transactions which consist of larger
// amounts, older inputs, and small sizes have the highest priority.  Second, a
// fee per kilobyte is calculated for each transaction.  Transactions with a
// higher fee per kilobyte are preferred.  Finally, the block generation related
// policy settings are all taken into account.
//
// Transactions which only spend outputs from other transactions already in the
// block chain are immediately added to a priority queue which either
// prioritizes based on the priority (then fee per kilobyte) or the fee per
// kilobyte (then priority) depending on whether or not the BlockPrioritySize
// policy setting allots space for high-priority transactions.  Transactions
// which spend outputs from other transactions in the source pool are added to a
// dependency map so they can be added to the priority queue once the
// transactions they depend on have been included.
//
// Once the high-priority area (if configured) has been filled with
// transactions, or the priority falls below what is considered high-priority,
// the priority queue is updated to prioritize by fees per kilobyte (then
// priority).
//
// When the fees per kilobyte drop below the TxMinFreeFee policy setting, the
// transaction will be skipped unless the BlockMinSize policy setting is
// nonzero, in which case the block will be filled with the low-fee/free
// transactions until the block size reaches that minimum size.
//
// Any transactions which would cause the block to exceed the BlockMaxSize
// policy setting, exceed the maximum allowed signature operations per block, or
// otherwise cause the block to be invalid are skipped.
//
// Given the above, a block generated by this function is of the following form:
//
//   -----------------------------------  --  --
//  |      Coinbase Transaction         |   |   |
//  |-----------------------------------|   |   |
//  |                                   |   |   | ----- policy.BlockPrioritySize
//  |   High-priority Transactions      |   |   |
//  |                                   |   |   |
//  |-----------------------------------|   | --
//  |                                   |   |
//  |                                   |   |
//  |                                   |   |--- (policy.BlockMaxSize) / 2
//  |  Transactions prioritized by fee  |   |
//  |  until <= policy.TxMinFreeFee     |   |
//  |                                   |   |
//  |                                   |   |
//  |                                   |   |
//  |-----------------------------------|   |
//  |  Low-fee/Non high-priority (free) |   |
//  |  transactions (while block size   |   |
//  |  <= policy.BlockMinSize)          |   |
//   -----------------------------------  --
//
//  This function returns nil, nil if there are not enough voters on any of
//  the current top blocks to create a new block template.
// TODO, refactor NewBlockTemplate input dependencies
func NewBlockTemplate(policy *Policy,config *config.Config, params *params.Params,
	sigCache *txscript.SigCache, source TxSource, tsource blockchain.MedianTimeSource,
	blkMgr *blkmgr.BlockManager,  payToAddress types.Address,parents []*hash.Hash) (*types.BlockTemplate, error) {
	txSource := source
	blockManager := blkMgr
	timeSource := tsource
	subsidyCache := blockManager.GetChain().FetchSubsidyCache()


	// All transaction scripts are verified using the more strict standarad
	// flags.
	scriptFlags, err := policy.StandardVerifyFlags()
	if err != nil {
		return nil, err
	}

	// Extend the most recently known best block.
	// The most recently known best block is the top block that has the most
	// TODO,refactor the poolsize & finalstate
	best:=blockManager.GetChain().BestSnapshot()
	nextBlockHeight := uint64(blockManager.GetChain().BlockDAG().GetMainChainTip().GetHeight() + 1)
	nextBlockOrder:=best.Order+1
	// Get the current source transactions and create a priority queue to
	// hold the transactions which are ready for inclusion into a block
	// along with some priority related and fee metadata.  Reserve the same
	// number of items that are available for the priority queue.  Also,
	// choose the initial sort order for the priority queue based on whether
	// or not there is an area allocated for high-priority transactions.
	sourceTxns := txSource.MiningDescs()
	sortedByFee := policy.BlockPrioritySize == 0
	// TODO, impl more general priority func
	lessFunc := txPQByFee
	if sortedByFee {
		lessFunc = txPQByFee
	}
	priorityQueue := newTxPriorityQueue(len(sourceTxns), lessFunc)

	// Create a slice to hold the transactions to be included in the
	// generated block with reserved space.  Also create a utxo view to
	// house all of the input transactions so multiple lookups can be
	// avoided.
	blockTxns := make([]*types.Tx, 0, len(sourceTxns))
	blockUtxos := blockchain.NewUtxoViewpoint()

	// dependers is used to track transactions which depend on another
	// transaction in the source pool.  This, in conjunction with the
	// dependsOn map kept with each dependent transaction helps quickly
	// determine which dependent transactions are now eligible for inclusion
	// in the block once each transaction has been included.
	dependers := make(map[hash.Hash]map[hash.Hash]*txPrioItem)

	// Create slices to hold the fees and number of signature operations
	// for each of the selected transactions and add an entry for the
	// coinbase.  This allows the code below to simply append details about
	// a transaction as it is selected for inclusion in the final block.
	// However, since the total fees aren't known yet, use a dummy value for
	// the coinbase fee which will be updated later.
	txFees := make([]int64, 0, len(sourceTxns))
	txFeesMap := make(map[hash.Hash]int64)
	txSigOpCounts := make([]int64, 0, len(sourceTxns))
	txSigOpCountsMap := make(map[hash.Hash]int64)
	txFees = append(txFees, -1) // Updated once known

	log.Debug("Inclusion to new block", "transactions",len(sourceTxns))
mempoolLoop:
	for _, txDesc := range sourceTxns {
		// A block can't have more than one coinbase or contain
		// non-finalized transactions.
		tx := txDesc.Tx
		msgTx := tx.Transaction()
		if msgTx.IsCoinBaseTx() {
			log.Trace("Skipping coinbase tx %s", tx.Hash())
			continue
		}
		if !blockchain.IsFinalizedTransaction(tx, nextBlockHeight,tsource.AdjustedTime()) {

			log.Trace("Skipping non-finalized tx %s", tx.Hash())
			continue
		}

		// Fetch all of the utxos referenced by the this transaction.
		// NOTE: This intentionally does not fetch inputs from the
		// mempool since a transaction which depends on other
		// transactions in the mempool must come after those
		utxos, err := blockManager.GetChain().FetchUtxoView(tx)
		if err != nil {
			log.Warn("Unable to fetch utxo view for tx %s: "+
				"%v", tx.Hash(), err)
			continue
		}

		// Setup dependencies for any transactions which reference
		// other transactions in the mempool so they can be properly
		// ordered below.
		prioItem := &txPrioItem{tx: txDesc.Tx, txType: txDesc.Type}
		for _, txIn := range tx.Transaction().TxIn {

			originHash := &txIn.PreviousOut.Hash
			if blockManager.GetChain().IsInvalidTx(originHash) {
				log.Trace("Skipping tx %s because it "+
					"references bad output %s "+
					"which is not available",
					tx.Hash(), txIn.PreviousOut)
				continue mempoolLoop
			}

			originIndex := txIn.PreviousOut.OutIndex
			utxoEntry := utxos.LookupEntry(originHash)
			if utxoEntry == nil || utxoEntry.IsOutputSpent(originIndex) {
				if !txSource.HaveTransaction(originHash) {
					log.Trace("Skipping tx %s because "+
						"it references unspent output "+
						"%s which is not available",
						tx.Hash(), txIn.PreviousOut)
					continue mempoolLoop
				}

				// The transaction is referencing another
				// transaction in the source pool, so setup an
				// ordering dependency.
				deps, exists := dependers[*originHash]
				if !exists {
					deps = make(map[hash.Hash]*txPrioItem)
					dependers[*originHash] = deps
				}
				deps[*prioItem.tx.Hash()] = prioItem
				if prioItem.dependsOn == nil {
					prioItem.dependsOn = make(
						map[hash.Hash]struct{})
				}
				prioItem.dependsOn[*originHash] = struct{}{}

				// Skip the check below. We already know the
				// referenced transaction is available.
				continue
			}
		}

		// Calculate the final transaction priority using the input
		// value age sum as well as the adjusted transaction size.  The
		// formula is: sum(inputValue * inputAge) / adjustedTxSize
		prioItem.priority = mempool.CalcPriority(tx.Transaction(), utxos,
			nextBlockOrder)

		// Calculate the fee in Atoms/KB.
		// NOTE: This is a more precise value than the one calculated
		// during calcMinRelayFee which rounds up to the nearest full
		// kilobyte boundary.  This is beneficial since it provides an
		// incentive to create smaller transactions.
		txSize := tx.Transaction().SerializeSize()
		prioItem.feePerKB = (float64(txDesc.Fee) * float64(kilobyte)) /
			float64(txSize)
		prioItem.fee = txDesc.Fee

		// Add the transaction to the priority queue to mark it ready
		// for inclusion in the block unless it has dependencies.
		if prioItem.dependsOn == nil {
			heap.Push(priorityQueue, prioItem)
		}

		// Merge the referenced outputs from the input transactions to
		// this transaction into the block utxo view.  This allows the
		// code below to avoid a second lookup.
		mergeUtxoView(blockUtxos, utxos)
	}

	log.Trace("Priority queue","queue len", priorityQueue.Len(),
		"dependers len", len(dependers))

	// The starting block size is the size of the block header plus the max
	// possible transaction count size, plus the size of the coinbase
	// transaction.
	blockSize := uint32(blockHeaderOverhead)

	// Guesstimate for sigops based on valid txs in loop below. This number
	// tends to overestimate sigops because of the way the loop below is
	// coded and the fact that tx can sometimes be removed from the tx
	// trees if they fail one of the stake checks below the priorityQueue
	// pop loop. This is buggy, but not catastrophic behaviour. A future
	// release should fix it. TODO
	blockSigOps := int64(0)
	totalFees := int64(0)
	// Choose which transactions make it into the block.
	for priorityQueue.Len() > 0 {
		// Grab the highest priority (or highest fee per kilobyte
		// depending on the sort order) transaction.
		prioItem := heap.Pop(priorityQueue).(*txPrioItem)
		tx := prioItem.tx

		// Grab the list of transactions which depend on this one (if any).
		deps := dependers[*tx.Hash()]

		// Enforce maximum block size.  Also check for overflow.
		txSize := uint32(tx.Transaction().SerializeSize())
		blockPlusTxSize := blockSize + txSize
		if blockPlusTxSize < blockSize || blockPlusTxSize >= policy.BlockMaxSize {
			log.Trace(fmt.Sprintf("Skipping tx %s (size %v) because it "+
				"would exceed the max block size; cur block "+
				"size %v, cur num tx %v", tx.Hash(), txSize,
				blockSize, len(blockTxns)))
			logSkippedDeps(tx, deps)
			continue
		}

		// Enforce maximum signature operations per block.  Also check
		// for overflow.
		numSigOps := int64(blockchain.CountSigOps(tx, false))
		if blockSigOps+numSigOps < blockSigOps ||
			blockSigOps+numSigOps > blockchain.MaxSigOpsPerBlock {
			log.Trace("Skipping tx %s because it would "+
				"exceed the maximum sigops per block", tx.Hash())
			logSkippedDeps(tx, deps)
			continue
		}

		// This isn't very expensive, but we do this check a number of times.
		// Consider caching this in the mempool in the future.
		numP2SHSigOps, err := blockchain.CountP2SHSigOps(tx, false,
			blockUtxos)
		if err != nil {
			log.Trace("Skipping tx %s due to error in "+
				"CountP2SHSigOps: %v", tx.Hash(), err)
			logSkippedDeps(tx, deps)
			continue
		}
		numSigOps += int64(numP2SHSigOps)
		if blockSigOps+numSigOps < blockSigOps ||
			blockSigOps+numSigOps > blockchain.MaxSigOpsPerBlock {
			log.Trace("Skipping tx %s because it would "+
				"exceed the maximum sigops per block (p2sh)",
				tx.Hash())
			logSkippedDeps(tx, deps)
			continue
		}

		// Skip free transactions once the block is larger than the
		// minimum block size, except for stake transactions.
		if sortedByFee &&
			(prioItem.feePerKB < float64(policy.TxMinFreeFee)) &&
			(blockPlusTxSize >= policy.BlockMinSize) {

			log.Trace("Skipping tx %s with feePerKB %.2f "+
				"< TxMinFreeFee %d and block size %d >= "+
				"minBlockSize %d", tx.Hash(), prioItem.feePerKB,
				policy.TxMinFreeFee, blockPlusTxSize,
				policy.BlockMinSize)
			logSkippedDeps(tx, deps)
			continue
		}

		// Prioritize by fee per kilobyte once the block is larger than
		// the priority size or there are no more high-priority
		// transactions.
		if !sortedByFee && (blockPlusTxSize >= policy.BlockPrioritySize ||
			prioItem.priority <= mempool.MinHighPriority) {

			log.Trace("Switching to sort by fees per "+
				"kilobyte blockSize %d >= BlockPrioritySize "+
				"%d || priority %.2f <= minHighPriority %.2f",
				blockPlusTxSize, policy.BlockPrioritySize,
				prioItem.priority,mempool.MinHighPriority)

			sortedByFee = true
			priorityQueue.SetLessFunc(txPQByFee)  //TODO, revisit the PQ func

			// Put the transaction back into the priority queue and
			// skip it so it is re-priortized by fees if it won't
			// fit into the high-priority section or the priority is
			// too low.  Otherwise this transaction will be the
			// final one in the high-priority section, so just fall
			// though to the code below so it is added now.
			if blockPlusTxSize > policy.BlockPrioritySize ||
				prioItem.priority < mempool.MinHighPriority {

				heap.Push(priorityQueue, prioItem)
				continue
			}
		}

		// Ensure the transaction inputs pass all of the necessary
		// preconditions before allowing it to be added to the block.
		// The fraud proof is not checked because it will be filled in
		// by the miner.
		_, err = blockchain.CheckTransactionInputs(subsidyCache, tx,
			int64(nextBlockOrder), blockUtxos, false, params ) //TODO, remove the params dependence
		if err != nil {
			log.Trace("Skipping tx %s due to error in "+
				"CheckTransactionInputs: %v", tx.Hash(), err)
			logSkippedDeps(tx, deps)
			continue
		}
		err = blockchain.ValidateTransactionScripts(tx, blockUtxos,
			scriptFlags, sigCache)
		if err != nil {
			log.Trace("Skipping tx %s due to error in "+
				"ValidateTransactionScripts: %v", tx.Hash(), err)
			logSkippedDeps(tx, deps)
			continue
		}

		// Spend the transaction inputs in the block utxo view and add
		// an entry for it to ensure any transactions which reference
		// this one have it available as an input and can ensure they
		// aren't double spending.
		err = spendTransaction(blockUtxos, tx, int64(nextBlockOrder)) //TODO, remove type conversion
		if err != nil {
			log.Warn("Unable to spend transaction %v in the preliminary "+
				"UTXO view for the block template: %v",
				tx.Hash(), err)
		}

		// Add the transaction to the block, increment counters, and
		// save the fees and signature operation counts to the block
		// template.
		blockTxns = append(blockTxns, tx)
		blockSize += txSize
		blockSigOps += numSigOps

		txFeesMap[*tx.Hash()] = prioItem.fee
		txSigOpCountsMap[*tx.Hash()] = numSigOps

		log.Trace(fmt.Sprintf("Adding tx %s (priority %.2f, feePerKB %.2f)",
			prioItem.tx.Hash(), prioItem.priority, prioItem.feePerKB))

		// Add transactions which depend on this one (and also do not
		// have any other unsatisified dependencies) to the priority
		// queue.
		for _, item := range deps {
			// Add the transaction to the priority queue if there
			// are no more dependencies after this one.
			delete(item.dependsOn, *tx.Hash())
			if len(item.dependsOn) == 0 {
				heap.Push(priorityQueue, item)
			}
		}
	}

	// Create a standard coinbase transaction paying to the provided
	// address.  NOTE: The coinbase value will be updated to include the
	// fees from the selected transactions later after they have actually
	// been selected.  It is created here to detect any errors early
	// before potentially doing a lot of work below.  The extra nonce helps
	// ensure the transaction is not a duplicate transaction (paying the
	// same value to the same public key address would otherwise be an
	// identical transaction for block version 1).
	coinbaseScript := []byte{0x00, 0x00}
	coinbaseScript = append(coinbaseScript, []byte(coinbaseFlags)...)

	// Add a random coinbase nonce to ensure that tx prefix hash
	// so that our merkle root is unique for lookups needed for
	// getwork, etc.
	rand, err := s.RandomUint64()
	if err != nil {
		return nil, err
	}
	opReturnPkScript, err := standardCoinbaseOpReturn(uint32(nextBlockHeight),
		rand)
	if err != nil {
		return nil, err
	}
	voters := 0  //TODO remove voters
	coinbaseTx, err := createCoinbaseTx(subsidyCache,
		coinbaseScript,
		opReturnPkScript,
		int64(nextBlockHeight),    //TODO remove type conversion
		int64(nextBlockOrder),
		payToAddress,
		uint16(voters),
		params)
	if err != nil {
		return nil, err
	}

	numCoinbaseSigOps := int64(blockchain.CountSigOps(coinbaseTx, true))
	blockSize += uint32(coinbaseTx.Transaction().SerializeSize())
	blockSigOps += numCoinbaseSigOps
	txFeesMap[*coinbaseTx.Hash()] = 0
	txSigOpCountsMap[*coinbaseTx.Hash()] = numCoinbaseSigOps

	// Build tx lists for regular tx.
	blockTxnsRegular := make([]*types.Tx, 0, len(blockTxns)+1)

	// Append coinbase.
	blockTxnsRegular = append(blockTxnsRegular, coinbaseTx)

	// Append regular tx
	blockTxnsRegular = append(blockTxnsRegular, blockTxns...)

	for _, tx := range blockTxnsRegular {
		fee, ok := txFeesMap[*tx.Hash()]
		if !ok {
			return nil, fmt.Errorf("couldn't find fee for tx %v",
				*tx.Hash())
		}
		totalFees += fee
		txFees = append(txFees, fee)

		tsos, ok := txSigOpCountsMap[*tx.Hash()]
		if !ok {
			return nil, fmt.Errorf("couldn't find sig ops count for tx %v",
				*tx.Hash())
		}
		txSigOpCounts = append(txSigOpCounts, tsos)
	}


	txSigOpCounts = append(txSigOpCounts, numCoinbaseSigOps)

	// Now that the actual transactions have been selected, update the
	// block size for the real transaction count and coinbase value with
	// the total fees accordingly.
	if nextBlockOrder > 1 {
		blockSize -= s.MaxVarIntPayload -
			uint32(s.VarIntSerializeSize(uint64(len(blockTxnsRegular))))
		coinbaseTx.Transaction().TxOut[2].Amount += uint64(totalFees)
		txFees[0] = -totalFees
	}

	// Calculate the required difficulty for the block.  The timestamp
	// is potentially adjusted to ensure it comes after the median time of
	// the last several blocks per the chain consensus rules.
	ts:= MedianAdjustedTime(blockManager.GetChain(),timeSource)
	reqDifficulty, err := blockManager.GetChain().CalcNextRequiredDifficulty(ts)

	if err != nil {
		return nil, miningRuleError(ErrGettingDifficulty, err.Error())
	}

	// Correct transaction index fraud proofs for any transactions that
	// are chains. maybeInsertStakeTx fills this in for stake transactions
	// already, so only do it for regular transactions.
	for i, tx := range blockTxnsRegular {
		// No need to check any of the transactions in the custom first
		// block.
		if nextBlockOrder == 1 {
			break
		}

		utxs, err := blockManager.GetChain().FetchUtxoView(tx)
		if err != nil {
			str := fmt.Sprintf("failed to fetch input utxs for tx %v: %s",
				tx.Hash(), err.Error())
			return nil, miningRuleError(ErrFetchTxStore, str)
		}

		// Copy the transaction and swap the pointer.
		txCopy := types.NewTxDeepTxIns(tx.Transaction())
		blockTxnsRegular[i] = txCopy
		tx = txCopy

		for _, txIn := range tx.Transaction().TxIn {
			originHash := &txIn.PreviousOut.Hash
			utx := utxs.LookupEntry(originHash)
			if utx == nil {
				// Set a flag with the index so we can properly set
				// the fraud proof below.
				txIn.TxIndex = types.NullTxIndex
			} else {
				originIdx := txIn.PreviousOut.OutIndex
				txIn.AmountIn = utx.AmountByIndex(originIdx)
				txIn.BlockOrder = uint32(utx.BlockOrder())
				txIn.TxIndex = utx.TxIndex()
			}
		}
	}

	// Fill in locally referenced inputs.
	for i, tx := range blockTxnsRegular {
		// Skip coinbase.
		if i == 0 {
			continue
		}

		// Copy the transaction and swap the pointer.
		txCopy := types.NewTxDeepTxIns(tx.Transaction())
		blockTxnsRegular[i] = txCopy
		tx = txCopy

		for _, txIn := range tx.Transaction().TxIn {
			// This tx was at some point 0-conf and now requires the
			// correct block height and index. Set it here.
			if txIn.TxIndex == types.NullTxIndex {
				idx := txIndexFromTxList(txIn.PreviousOut.Hash,
					blockTxnsRegular)

				// The input is in the block, set it accordingly.
				if idx != -1 {
					originIdx := txIn.PreviousOut.OutIndex
					amt := blockTxnsRegular[idx].Transaction().TxOut[originIdx].Amount
					txIn.AmountIn = amt
					txIn.BlockOrder = uint32(nextBlockOrder)   //TODO,remove type conversion
					txIn.TxIndex = uint32(idx)
				} else {
					str := fmt.Sprintf("failed find hash in tx list "+
						"for fraud proof; tx in hash %v",
						txIn.PreviousOut.Hash)
					return nil, miningRuleError(ErrFraudProofIndex, str)
				}
			}
		}
	}

	// Choose the block version to generate based on the network.
	blockVersion := uint32(generatedBlockVersion)
	if params.Net != protocol.MainNet {
		blockVersion = generatedBlockVersionTest
	}

	// Create a new block ready to be solved.
	merkles := merkle.BuildMerkleTreeStore(blockTxnsRegular)
	if parents==nil {
		parents=blockManager.GetChain().BlockDAG().GetTips().SortList(false)
	}

	paMerkles :=merkle.BuildParentsMerkleTreeStore(parents)
	var block types.Block
	block.Header = types.BlockHeader{
		Version:      blockVersion,
		ParentRoot:   *paMerkles[len(paMerkles)-1],
		TxRoot:       *merkles[len(merkles)-1],
		StateRoot:    hash.Hash{}, //TODO, state root
		Timestamp:    ts,
		Difficulty:   reqDifficulty,
		// Size declared below
	}
	for _,pb:=range parents{
		if err := block.AddParent(pb); err != nil {
			return nil, err
		}
	}
	for _, tx := range blockTxnsRegular {
		if err := block.AddTransaction(tx.Transaction()); err != nil {
			return nil, miningRuleError(ErrTransactionAppend, err.Error())
		}
	}

	//TODO revisit the size in block header
	/*
	msgBlock.Header.Size = uint32(msgBlock.SerializeSize())
	*/

	// Finally, perform a full check on the created block against the chain
	// consensus rules to ensure it properly connects to the current best
	// chain with no issues.

	sblock := types.NewBlockDeepCopyCoinbase(&block)
	sblock.SetOrder(nextBlockOrder)
	sblock.SetHeight(uint(nextBlockHeight))
	err = blockManager.GetChain().CheckConnectBlockTemplate(sblock)
	if err != nil {
		str := fmt.Sprintf("failed to do final check for check connect "+
			"block when making new block template: %v",
			err.Error())
		return nil, miningRuleError(ErrCheckConnectBlock, str)
	}

	log.Debug("Created new block template",
		"transactions", len(block.Transactions),
		"fees",totalFees,
		"signOp",blockSigOps,
		"bytes", blockSize,
		"target",
		fmt.Sprintf("%064x",blockchain.CompactToBig(block.Header.Difficulty)))

	blockTemplate := &types.BlockTemplate{
		Block:           &block,
		Fees:            txFees,
		SigOpCounts:     txSigOpCounts,
		Height:          nextBlockHeight,
		ValidPayAddress: payToAddress != nil,
	}
	return handleCreatedBlockTemplate(blockTemplate, blockManager)
}

// UpdateBlockTime updates the timestamp in the header of the passed block to
// the current time while taking into account the median time of the last
// several blocks to ensure the new time is after that time per the chain
// consensus rules.  Finally, it will update the target difficulty if needed
// based on the new time for the test networks since their target difficulty can
// change based upon time.
func UpdateBlockTime(msgBlock *types.Block, chain *blockchain.BlockChain, timeSource blockchain.MedianTimeSource,
	activeNetParams *params.Params) error {

	// The new timestamp is potentially adjusted to ensure it comes after
	// the median time of the last several blocks per the chain consensus
	// rules.
	newTimestamp:=MedianAdjustedTime(chain,timeSource)
	msgBlock.Header.Timestamp = newTimestamp

	// If running on a network that requires recalculating the difficulty,
	// do so now.
	if activeNetParams.ReduceMinDifficulty {
		difficulty, err := chain.CalcNextRequiredDifficulty(
			newTimestamp)
		if err != nil {
			return miningRuleError(ErrGettingDifficulty, err.Error())
		}
		msgBlock.Header.Difficulty = difficulty
	}

	return nil
}

// mergeUtxoView adds all of the entries in view to viewA.  The result is that
// viewA will contain all of its original entries plus all of the entries
// in viewB.  It will replace any entries in viewB which also exist in viewA
// if the entry in viewA is fully spent.
func mergeUtxoView(viewA *blockchain.UtxoViewpoint, viewB *blockchain.UtxoViewpoint) {
	viewAEntries := viewA.Entries()
	for h, entryB := range viewB.Entries() {
		if entryA, exists := viewAEntries[h]; !exists ||
			entryA == nil || entryA.IsFullySpent() {
			viewAEntries[h] = entryB
		}
	}
}

// TODO, move the log logic
// logSkippedDeps logs any dependencies which are also skipped as a result of
// skipping a transaction while generating a block template at the trace level.
func logSkippedDeps(tx *types.Tx, deps map[hash.Hash]*txPrioItem) {
	if deps == nil {
		return
	}

	for _, item := range deps {
		log.Trace("Skipping tx %s since it depends on %s\n",
			item.tx.Hash(), tx.Hash())
	}
}

// spendTransaction updates the passed view by marking the inputs to the passed
// transaction as spent.  It also adds all outputs in the passed transaction
// which are not provably unspendable as available unspent transaction outputs.
func spendTransaction(utxoView *blockchain.UtxoViewpoint, tx *types.Tx, order int64) error {
	for _, txIn := range tx.Transaction().TxIn {
		originHash := &txIn.PreviousOut.Hash
		originIndex := txIn.PreviousOut.OutIndex
		entry := utxoView.LookupEntry(originHash)
		if entry != nil {
			entry.SpendOutput(originIndex)
		}

	}

	utxoView.AddTxOuts(tx, order, types.NullTxIndex)
	return nil
}

// standardCoinbaseOpReturn creates a standard OP_RETURN output to insert into
// coinbase to use as extranonces. The OP_RETURN pushes 32 bytes.
func standardCoinbaseOpReturn(height uint32, extraNonce uint64) ([]byte, error) {
	enData := make([]byte, 12)
	binary.LittleEndian.PutUint32(enData[0:4], height)
	binary.LittleEndian.PutUint64(enData[4:12], extraNonce)
	extraNonceScript, err := txscript.GenerateProvablyPruneableOut(enData)
	if err != nil {
		return nil, err
	}

	return extraNonceScript, nil
}
// createCoinbaseTx returns a coinbase transaction paying an appropriate subsidy
// based on the passed block height to the provided address.  When the address
// is nil, the coinbase transaction will instead be redeemable by anyone.
//
// See the comment for NewBlockTemplate for more information about why the nil
// address handling is useful.
func createCoinbaseTx(subsidyCache *blockchain.SubsidyCache, coinbaseScript []byte, opReturnPkScript []byte, nextBlockHeight,nextBlockOrder int64, addr types.Address, voters uint16, params *params.Params) (*types.Tx, error) {
	tx := types.NewTransaction()
	tx.AddTxIn(&types.TxInput{
		// Coinbase transactions have no inputs, so previous outpoint is
		// zero hash and max index.
		PreviousOut: *types.NewOutPoint(&hash.Hash{},
			types.MaxPrevOutIndex ),
		Sequence:        types.MaxTxInSequenceNum,
		BlockOrder:      types.NullBlockOrder,
		TxIndex:         types.NullTxIndex,
		SignScript:      coinbaseScript,
	})

	// Block one is a special block that might pay out tokens to a ledger.
	if nextBlockOrder == 1 && len(params.BlockOneLedger) != 0 {
		// Convert the addresses in the ledger into useable format.
		addrs := make([]types.Address, len(params.BlockOneLedger))
		for i, payout := range params.BlockOneLedger {
			addr, err := address.DecodeAddress(payout.Address)
			if err != nil {
				return nil, err
			}
			addrs[i] = addr
		}

		for i, payout := range params.BlockOneLedger {
			// Make payout to this address.
			pks, err := txscript.PayToAddrScript(addrs[i])
			if err != nil {
				return nil, err
			}
			tx.AddTxOut(&types.TxOutput{
				Amount:   payout.Amount,
				PkScript: pks,
			})
		}

		tx.TxIn[0].AmountIn = params.BlockOneSubsidy()

		return types.NewTx(tx), nil
	}

	// Create a coinbase with correct block subsidy and extranonce.
	subsidy := blockchain.CalcBlockWorkSubsidy(subsidyCache,
		nextBlockHeight,
		voters,
		params)
	tax := blockchain.CalcBlockTaxSubsidy(subsidyCache,
		nextBlockHeight,
		voters,
		params)

	// Tax output.
	if params.BlockTaxProportion > 0 {
		tx.AddTxOut(&types.TxOutput{
			Amount:    uint64(tax),
			PkScript: params.OrganizationPkScript,
		})
	} else {
		// Tax disabled.
		scriptBuilder := txscript.NewScriptBuilder()
		trueScript, err := scriptBuilder.AddOp(txscript.OP_TRUE).Script()
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(&types.TxOutput{
			Amount:    uint64(tax),
			PkScript: trueScript,
		})
	}
	// Extranonce.
	tx.AddTxOut(&types.TxOutput{
		Amount:    0,
		PkScript: opReturnPkScript,
	})
	// AmountIn.
	tx.TxIn[0].AmountIn = subsidy + uint64(tax)  //TODO, remove type conversion

	// Create the script to pay to the provided payment address if one was
	// specified.  Otherwise create a script that allows the coinbase to be
	// redeemable by anyone.
	var pksSubsidy []byte
	if addr != nil {
		var err error
		pksSubsidy, err = txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		scriptBuilder := txscript.NewScriptBuilder()
		pksSubsidy, err = scriptBuilder.AddOp(txscript.OP_TRUE).Script()
		if err != nil {
			return nil, err
		}
	}
	// Subsidy paid to miner.
	tx.AddTxOut(&types.TxOutput{
		Amount:    subsidy,
		PkScript: pksSubsidy,
	})

	return types.NewTx(tx), nil
}

// txIndexFromTxList returns a transaction's index in a list, or -1 if it
// can not be found.
func txIndexFromTxList(hash hash.Hash, list []*types.Tx) int {
	for i, tx := range list {
		h := tx.Hash()
		if hash == *h {
			return i
		}
	}

	return -1
}

// handleCreatedBlockTemplate stores a successfully created block template to
// the appropriate cache if needed, then returns the template to the miner to
// work on. The stored template is a copy of the template, to prevent races
// from occurring in case the template is mined on by the CPUminer.
// TODO, revisit the block template cache design
func handleCreatedBlockTemplate(blockTemplate *types.BlockTemplate, bm *blkmgr.BlockManager) (*types.BlockTemplate, error) {
	curTemplate := bm.GetCurrentTemplate()

	nextBlockHeight := blockTemplate.Height

	// Overwrite the old cached block if it's out of date.
	if curTemplate != nil {
		if curTemplate.Height == nextBlockHeight {
			bm.SetCurrentTemplate(blockTemplate)
		}
	}

	return blockTemplate, nil
}

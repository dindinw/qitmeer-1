#!/bin/bash

set -e

# ---------------------------
# settings
# ---------------------------
DEBUG=0
DEBUG_FILE=/tmp/myeth_curl_debug
ERR_FILE=/tmp/myeth_curl_error

# ---------------------------
# solc call, need solc command line
# ---------------------------

#  eth_compileSolidity removed, use solc comand instaed
function eth_compile(){
  # local payload='{"jsonrpc":"2.0","method":"eth_compileSolidity","params":["'$code'"],"id":1}'
  # get_result "$payload"
  local code="$1"
  local suppress_error=0
  local opts=""
  shift
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -q) suppress_error=1; shift ;;
      -all) opts="--bin --abi --hashes"; shift;;
      *) opts+=" $1"; shift;;
    esac
  done
  if [ "$opts" == "" ]; then
    echo "$code" | solc --bin 2>/dev/null | awk /^[0-9]/
    return
  elif [ "$opts" == " -bin" ] ;then  # by-default, out cleanup binary
    echo "$code" | solc --bin 2>/dev/null | awk /^[0-9]/ |jq -R {"binary":.}
    return
  elif [ "$opts" == " -abi" ];then
    echo "$code" | solc --abi 2>/dev/null | tail -n 1 |jq .
    return
  elif [ "$opts" == " -func" ];then
    echo "$code" | solc --hashes 2>/dev/null |awk /^[0-9a-f]/|jq -R '.|split(": ")|{name:.[1],sign: .[0]}'
    return
  fi;
  #echo "debug opts=$opts s_err=$suppress_error"
  if [[ $suppress_error -eq 0 ]]; then
    echo "$code" | solc $opts
  else
    echo "$code" | solc $opts 2>/dev/null
  fi
}


# ---------------------------
# All jsonrpc calls
# https://github.com/ethereum/wiki/wiki/JSON-RPC
# https://github.com/ethereum/go-ethereum/wiki/Management-APIs
# https://github.com/ethereum/wiki/wiki/JavaScript-API
# https://wiki.parity.io/JSONRPC-eth-module.html
# ---------------------------

# eth_call
#   executes a message call which is directly/immediately executed
#   in the VM of the node,  but never mined into the blockchain.
#   aka, without creating a transaction on the block chain.
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) Call(ctx context.Context, args CallArgs, blockNr rpc.BlockNumber) (hexutil.Bytes, error)
function eth_call(){
  local params="{"
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -from)        shift; params+='"from":"'$1'",'; shift ;;
      -to)          shift; params+='"to":"'$1'",'; shift ;;
      -v|-value)    shift; params+='"value":"'$1'",'; shift ;;
      -d|-data)     shift; params+='"data":"'$1'",'; shift ;;
      -b|-blocknum) shift; block_num=$1; shift ;;
      *)            shift;;
    esac
  done
  if [ "$block_num" == "" ]; then
    block_num="latest"
  fi
  params=${params%,}"}"  # remove the last , and add }
  #echo "debug params=$params"
  local payload='{"jsonrpc":"2.0","method":"eth_call","params":['$params',"'$block_num'"],"id":1}'
  get_result "$payload"
}


#
# debug_dumpBlock
#   retrieves the entire state of the database at a given block.
# geth : eth/api.go
#   func (api *PublicDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error)
#
function dump_block(){
  local number=
  if [ "${1:0:2}" == "0x" ];then
      number=$1
  else
      number=$(to_hex_with_0x_prefix $1)
  fi

  local data='{"jsonrpc":"2.0","method":"debug_dumpBlock","params":["'$number'"],"id":1}'
  get_result "$data"
}

#
# rpc_modules
#   returns the list of RPC services with their version number
# geth : rpc/server.go
#   func (s *RPCService) Modules() map[string]string
#
function get_rpc_modules(){
  local data='{"jsonrpc":"2.0","method":"rpc_modules","params":[],"id":1}'
  get_result "$data"
}

#
# eth_mining
#   if this node is currently mining
# geth : eth/api.go
#   func (api *PublicMinerAPI) Mining() bool
function get_mining(){
  local data='{"jsonrpc":"2.0","method":"eth_mining","params":[],"id":1}'
  get_result "$data"
}

#
# // GetWork returns a work package for external miner. The work package consists of 3 strings
# // result[0], 32 bytes hex encoded current block header pow-hash
# // result[1], 32 bytes hex encoded seed hash used for DAG
# // result[2], 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
# geth : eth/api.go
#   func (api *PublicMinerAPI) GetWork() ([3]string, error)
function get_work(){
  local data='{"jsonrpc":"2.0","method":"eth_getWork","params":[],"id":1}'
  get_result "$data"
}

# miner_start
#   Start the miner with the given number of threads.
# geth : eth/api.go
#   func (api *PrivateMinerAPI) Start(threads *int) error
function miner_start(){
  local data='{"jsonrpc":"2.0","method":"miner_start","params":[],"id":1}'
  get_result "$data"
}

# miner_stop
#   Stop the miner
# geth : eth/api.go
#   func (api *PrivateMinerAPI) Stop() bool
function miner_stop(){
  local data='{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":1}'
  get_result "$data"
}

# returns the POW hashrate
# geth : eth/api.go
#   func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64
function get_hashrate(){
  local data='{"jsonrpc":"2.0","method":"eth_hashrate","params":[],"id":71}'
  get_result "$data"
}

# GasPrice returns a suggestion for a gas price.
# geth :internal/ethapi/api.go
#   func (s *PublicEthereumAPI) GasPrice(ctx context.Context) (*big.Int, error)
function get_gasprice(){
  local data='{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}'
  get_result "$data"
}

# return the address that mining rewards will be send to (alias for Etherbase)
# geth : eth/api.go
#   func (api *PublicEthereumAPI) Coinbase()
function get_coinbase(){
  local data='{"jsonrpc":"2.0","method":"eth_coinbase","params":[],"id":1}'
  get_result "$data"
}

# the collection of accounts this node manages
# geth : eth/api.go
#  func (s *PublicAccountAPI) Accounts() []common.Address
function get_accounts(){
  local data='{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}'
  get_result "$data"
}

# returns the transaction for the given block hash and index.
# It's just a language suger, internally it look up block by block num frist then get tx by index
# geth : internal/ethapi/api.go
#  func (s *PublicTransactionPoolAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction
function get_tx_by_blocknum_and_index_hex(){
  local block_num=$1 #"0x467a65"
  local tx_index=$2  #"0x0"
  local data='{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params":["'$block_num'","'$tx_index'"],"id":1}'
  get_result "$data"
}

# returns the transaction for the given hash
# geth : internal/ethapi/api.go
#   func (s *PublicTransactionPoolAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) *RPCTransaction
function get_tx_by_hash(){
  local tx_hash=$1
  local data='{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["'$tx_hash'"],"id":1}'
  get_result "$data"
}

# the block number of the chain head
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) BlockNumber() *big.Int
function get_block_number(){
  local data='{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}'
  get_result "$data"
}

# Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
# yet received the latest block headers from its pears. In case it is synchronizing:
#   - startingBlock: block number this node started to synchronise from
#   - currentBlock:  block number this node is currently importing
#   - highestBlock:  block number of the highest block header this node has received from peers
#   - pulledStates:  number of state entries processed until now
#   - knownStates:   number of known state entries that still need to be pulled
# geth : internal/ethapi/api.go
#   func (s *PublicEthereumAPI) Syncing() (interface{}, error)
function get_syncing(){
  local data='{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}'
  get_result "$data"
}

# returns the requested block by blockNr
# When fullTx is true all transactions in the block are # returned in full detail, otherwise only the transaction hash is returned.
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) GetBlockByNumber(ctx context.Context, blockNr rpc.BlockNumber, fullTx bool) (map[string]interface{}, error)
function get_block(){
  local block_number=$1
  local hex_prefix=${block_number:0:2}
  local num_hex=${block_number:2}
  if [ ! "$hex_prefix" == "0x" ] ;then
    # $block_number not using hex_prefix, try to convert hex"
    num_hex=$(echo "obase=16;$block_number"|bc)
  fi
  local fullTx=$2
  if [ "$fullTx" == "" ]; then
    fullTx="true"
  fi
  local data='{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x'$num_hex'",'$fullTx'],"id":1}'
  get_result "$data"
}
# return block by hash
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error)
function get_block_by_hash(){
  local block_hash=$1
  local data='{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["'$block_hash'",true],"id":1}'
  get_result "$data"
}

# the amount of wei for the given address in the state of the given block number.
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*big.Int, error)
function get_balance(){
  local addr=$1
  local block_num=$2
  if [ "$block_num" == "" ]; then
    #echo '{"result":"block_tag is latest and addr is '$addr'"}'
    block_num="latest"
  fi
  local data='{"jsonrpc":"2.0","method":"eth_getBalance","params":["'$addr'","'$block_num'"],"id":1}'
  get_result "$data"
}

# the code stored at the given address in the state for the given block number
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (hexutil.Bytes, error)
function get_code(){
  #local contract_addr="0xbdeb5b87843062116b118e574a68a58f511a30e6"
  local contract_addr=$1
  local block_num=$2
  if [ "$block_num" == "" ]; then
    block_num="latest"
  fi
  local data='{"jsonrpc":"2.0","method":"eth_getCode","params":["'$contract_addr'","'$block_num'"],"id":1}'
  get_result "$data"
}

# the storage from the state at the given address, key and block numbe
# geth : internal/ethapi/api.go
#   func (s *PublicBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNr rpc.BlockNumber) (hexutil.Bytes, error)
function get_storage(){
  #local contract_addr="0xbdeb5b87843062116b118e574a68a58f511a30e6"
  #local at="0x0"
  local contract_addr=$1
  local at=$2
  local block_num=$3
  if [ "$block_num" == "" ]; then
    block_num="latest"
  fi
  local data='{"jsonrpc":"2.0","method":"eth_getStorageAt","params":["'$contract_addr'","'$at'","'$block_num'"],"id":1}'
  get_result "$data"
}

# the number of transactions the given address has sent for the given block numbe
# acutally it's return the nouce of the account(address) of the stateDB of the given blockNr
# geth : internal/ethapi/api.go
#   func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*hexutil.Uint64, error)
function get_tx_count_by_addr(){
  local addr=$1
  local data='{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["'$addr'","latest"],"id":1}'
  get_result "$data"
}

function get_result(){
  set +x
  #local site=10.0.0.6
  #local site=10.0.0.8
  if [ -z "$host" ]; then
     host=127.0.0.1
  fi
  if [ -z "$port" ]; then
     port=8545
  fi
  local data=$1
  local curl_result=$(curl -s -X POST -H 'Content-Type: application/json' --data $data http://$host:$port)
  local result=$(echo $curl_result|jq -r -c -M '.result')
  if [ $DEBUG -gt 0 ]; then
    local curl_cmd="curl -s -X POST -H 'Content-Type: application/json' --data '"$data"' http://$host:$port"
    echo "$curl_cmd" > $DEBUG_FILE
    echo "$curl_result" >> $DEBUG_FILE
  fi
  if [ "$result" == "null" ];then
    echo "$curl_result" > $ERR_FILE
    result=""
  fi
  echo $result
}

# -------------------------
# util functions
# -------------------------
function pad_hex_prefix(){
  local input=$1
  if [ "${input:0:2}" == "0x" ];then
    echo "$input"
  else
    echo "0x$input"
  fi
}
function to_hex() {
  printf "%x\n" $1
}
function to_hex_with_0x_prefix(){
  printf "0x%x\n" $1
}
function check_error() {
  if [ -e $ERR_FILE ]; then
    echo "Error:"
    cat $ERR_FILE|jq .error
    rm $ERR_FILE
  fi
}

# The python script for eth currency.
# the orignal script from https://github.com/ethereum/eth-utils
# the -d '\0' trick makes read exit with status 0
read -r -d '\0' py_eth_currency <<-"EOF"
import decimal
from decimal import localcontext

units = {
    'wei':          decimal.Decimal('1'),  # noqa: E241
    'kwei':         decimal.Decimal('1000'),  # noqa: E241
    'babbage':      decimal.Decimal('1000'),  # noqa: E241
    'femtoether':   decimal.Decimal('1000'),  # noqa: E241
    'mwei':         decimal.Decimal('1000000'),  # noqa: E241
    'lovelace':     decimal.Decimal('1000000'),  # noqa: E241
    'picoether':    decimal.Decimal('1000000'),  # noqa: E241
    'gwei':         decimal.Decimal('1000000000'),  # noqa: E241
    'shannon':      decimal.Decimal('1000000000'),  # noqa: E241
    'nanoether':    decimal.Decimal('1000000000'),  # noqa: E241
    'nano':         decimal.Decimal('1000000000'),  # noqa: E241
    'szabo':        decimal.Decimal('1000000000000'),  # noqa: E241
    'microether':   decimal.Decimal('1000000000000'),  # noqa: E241
    'micro':        decimal.Decimal('1000000000000'),  # noqa: E241
    'finney':       decimal.Decimal('1000000000000000'),  # noqa: E241
    'milliether':   decimal.Decimal('1000000000000000'),  # noqa: E241
    'milli':        decimal.Decimal('1000000000000000'),  # noqa: E241
    'ether':        decimal.Decimal('1000000000000000000'),  # noqa: E241
    'kether':       decimal.Decimal('1000000000000000000000'),  # noqa: E241
    'grand':        decimal.Decimal('1000000000000000000000'),  # noqa: E241
    'mether':       decimal.Decimal('1000000000000000000000000'),  # noqa: E241
    'gether':       decimal.Decimal('1000000000000000000000000000'),  # noqa: E241
    'tether':       decimal.Decimal('1000000000000000000000000000000'),  # noqa: E241
}

denoms = type('denoms', (object,), {
    key: int(value) for key, value in units.items()
})

MIN_WEI = 0
MAX_WEI = 2 ** 256 - 1

def is_integer(value):
    return isinstance(value, (int,)) and not isinstance(value, bool)

def is_string(value):
    return isinstance(value, (bytes, str, bytearray))

def from_wei(number, unit):
    """
    Takes a number of wei and converts it to any other ether unit.
    """
    if unit.lower() not in units:
        raise ValueError(
            "Unknown unit.  Must be one of {0}".format('/'.join(units.keys()))
        )

    if number == 0:
        return 0

    if number < MIN_WEI or number > MAX_WEI:
        raise ValueError("value must be between 1 and 2**256 - 1")

    unit_value = units[unit.lower()]

    with localcontext() as ctx:
        ctx.prec = 999
        d_number = decimal.Decimal(value=number, context=ctx)
        result_value = d_number / unit_value

    return result_value

def to_wei(number, unit):
    """
    Takes a number of a unit and converts it to wei.
    """
    if unit.lower() not in units:
        raise ValueError(
            "Unknown unit.  Must be one of {0}".format('/'.join(units.keys()))
        )

    if is_integer(number) or is_string(number):
        d_number = decimal.Decimal(value=number)
    elif isinstance(number, float):
        d_number = decimal.Decimal(value=str(number))
    elif isinstance(number, decimal.Decimal):
        d_number = number
    else:
        raise TypeError("Unsupported type.  Must be one of integer, float, or string")

    s_number = str(number)
    unit_value = units[unit.lower()]

    if d_number == 0:
        return 0

    if d_number < 1 and '.' in s_number:
        with localcontext() as ctx:
            multiplier = len(s_number) - s_number.index('.') - 1
            ctx.prec = multiplier
            d_number = decimal.Decimal(value=number, context=ctx) * 10**multiplier
        unit_value /= 10**multiplier

    with localcontext() as ctx:
        ctx.prec = 999
        result_value = decimal.Decimal(value=d_number, context=ctx) * unit_value

    if result_value < MIN_WEI or result_value > MAX_WEI:
        raise ValueError("Resulting wei value must be between 1 and 2**256 - 1")

    return int(result_value)
\0
EOF

function wei_to_ether() {
    # supress the scientific notation
    # 1.8E-08 -> 0.000000018
    python << EOF
$py_eth_currency
result=from_wei(decimal.Decimal($1),'ether')
out="{:.20f}".format(result)
frac=out.split('.')[1]
# by default, remove all fractional part.
no_zero=len(frac)-1
for i,v in enumerate(reversed(frac)):
    if int(v) > 0 :
        #print "first no zero -> i=%s,v=%s" % (i,v)
        no_zero=i
        break
print out[:len(out)-no_zero]
EOF
}

function ether_to_wei() {
    python << EOF
$py_eth_currency
print "%s" % to_wei(str('$1'),'ether')
EOF
}


function usage(){

  echo "chain    :"
  echo "  get_block_number"
  echo "  get_syncing_info"
  echo "  get_current_block"
  echo "  get_current_block2 <num|hash> [-tx |-txcount|-blocktime|...]"
  echo "  get_highest_block"
  echo "block    :"
  echo "  get_block <num|hash>"
  echo "  get_block <num|hash> -tx [num]"
  echo "  get_block <num|hash> -txcount"
  echo "  get_block <num|hash> -blocktime"
  echo "  get_block <num|hash> -stroot"
  echo "  get_block <num|hash> -txroot"
  echo "  get_block <num|hash> -rcroot"
  echo "  get_block <num|hash> -roots"
  echo "tx       :"
  echo "  get_tx <hash>"
  echo "  get_tx_by_block_and_index <num_hex> <index_hex>"
  echo "account  :"
  echo "  get_accounts"
  echo "  get_balance <addr> [blocknum]"
  echo "  get_tx_count <addr>"
  echo "  get_code <addr> [blocknum]"
  echo "  get_storage <addr> <at> [blocknum]"
  echo "mining  :"
  echo "  get_coinbase"
  echo "  set_coinbase"
  echo "  start_mining"
  echo "  stop_mining"
  echo "status  :"
  echo "  status|get_status|info|get_info [-mining|-module|-hashrate|-work|-all]"
  echo "compile :"
  echo "  eth_complie [-bin|-abi|-fun|-all] [-q]"
  echo "excute  :"
  echo "  dump_block <num>"
  echo "  eth_call -from <addr> -to <addr> -v <value> -d <data>"
  echo "util    :"
  echo "  to_ether <wei>"
  echo "  to_wei <ether>"

}

# -------------------
# level 2 functions
# -------------------

function start_mining(){
  check_api_error_stop miner
  local mining=$(get_mining)
  if [ "$mining" == "false" ]; then
    echo start mining ...
    $(miner_start)
    get_mining
  else
    echo already stated
    get_mining_status
  fi
}

function stop_mining(){
  check_api_error_stop miner
  local mining=$(get_mining)
  if [ "$mining" == "true" ]; then
    echo stop mining ...
    miner_stop
  else
    echo already stopped
    get_mining_status
  fi
}

function check_api_error_stop(){
  if [ "$(check_api $1)" == "" ]; then
    echo "$1 api not find, need to enable the management APIs"
    echo "For example: geth --rpcapi eth,web3,miner --rpc"
    exit
  fi
}

function check_api() {
  local api=$1
  if ! [ "$api" == "" ]; then
    local api_ver=$(get_rpc_modules|jq .$api -r)
     if ! [ "$api_ver" == "null" ]; then
         echo $api_ver
     fi
  fi
}
function get_mining_status(){
  local mining=$(get_mining)
  local hashrate=$(get_hashrate)
  local gasprice=$(get_gasprice)
  local gasprice_dec=$(printf "%d" $gasprice)
  echo "mining   : $mining"
  echo "hashRate : $hashrate"
  echo "gasPrice : $gasprice -> ($gasprice_dec wei/$(wei_to_ether $gasprice) ether)"
}
function get_modules(){
  get_rpc_modules|jq . -r
}

function get_status() {
  #echo "debug get_status $@"
  if [ "$1" == "" ]; then
    get_modules
  elif [ "$1" == "-module" ]; then
    get_modules
  elif [ "$1" == "-mining" ]; then
    get_mining_status
  elif [ "$1" == "-hashrate" ]; then
    get_hashrate
  elif [ "$1" == "-work" ]; then
    get_work|jq .
    check_error
  elif [ "$1" == "-all" ]; then
    echo "modules  : $(get_modules|jq -c -M .)"
    get_mining_status
  else
    echo "unknown opt $@"
  fi
}

function get_current_block_num(){
  get_syncing $@|jq .currentBlock -r|xargs printf "%d\n"
}

# get_block control function
function cmd_get_block(){
  # echo "debug cmd_get_block $@"
  if [ "$1" == "" ] ;then
      echo "get lastet block"
      blocknum=$(get_block_number|xargs printf "%d")
      echo "the lastet block is $blocknum"
      exit
  fi
  if [ "${1:0:2}" == "0x" ];then
    block_hash=$1
    block_result=$(get_block_by_hash $(pad_hex_prefix $block_hash) $@)
  else
    blocknum=$1
    block_result=$(get_block "$blocknum")
  fi
  shift
  #echo debug $block_result
  if [ "$block_result" == "null" ];then
    echo "block $1 not found"
    exit -1
  fi
  if [ -z "$1" ]; then
    echo $block_result|jq '.'
  elif [ "$1" == "-tx" ]; then
    shift
    tx=$1
    if [ "${tx:0:2}" == "0x" ];then
      echo $block_result|jq '.transactions|.[]|select(.transactionIndex == "'$tx'")'
    else
      echo $block_result|jq '.transactions['$tx']'
    fi
  elif [ "$1" == "-txcount" ];then
    shift
    echo $block_result|jq '.transactions|length'
  elif [ "$1" == "-blocktime" ];then
    shift
    echo $block_result|jq '.timestamp'| hex2dec|timestamp
  elif [ "$1" == "-stroot" ];then
    echo $block_result|jq '.stateRoot'
    shift
  elif [ "$1" == "-txroot" ];then
    echo $block_result|jq '.transactionsRoot'
    shift
  elif [ "$1" == "-rcroot" ];then
    echo $block_result|jq '.receiptsRoot'
    shift
  elif [ "$1" == "-roots" ]; then
    echo $block_result|jq '{"stroot":.stateRoot, "txroot":.transactionsRoot, "rcroot":.receiptsRoot}'
  fi
}

# main logic
if [ $? != 0 ]; then
  echo "Usage: -h [host] -p [port] "
  exit;
fi
#echo "$@"
while [ $# -gt 0 ] ;do
  case "$1" in
    -h)
      host=$2
      #echo "host is $host"
      shift;;
    -p)
      port=$2
      #echo "port is $port"
      shift;;
    -D)
      DEBUG=1
      ;;
    *)
      cmd="$@"
      #echo "cmd is $cmd"
      break;;
    esac
  shift
done


## Block
if [ $1 == "get_block" ]; then
  shift
  cmd_get_block $@
elif [ $1 == "get_block_number" ]; then
  shift
  if [ "$1" == "-hex" ]; then # result ishex by default
    get_block_number
  else                        # human can read (hex->decimal)
    get_block_number |xargs printf "%d\n"
  fi
elif [ $1 == "get_syncing_info" ]; then
  shift
  get_syncing $@
elif [ $1 == "get_current_block" ]; then
  shift
  get_current_block_num
elif [ $1 == "get_current_block2" ]; then
  shift
  blocknum=$(get_current_block_num)
  cmd_get_block $blocknum $@
elif [ $1 == "get_highest_block" ]; then
  shift
  get_syncing $@|jq .highestBlock -r|xargs printf "%d\n"

## Tx
elif [ $1 == "get_tx" ]; then
  shift
  get_tx_by_hash $@ |jq .
elif [ $1 == "get_tx_by_block_and_index" ]; then
  shift
  # note: the input is block number & tx index in hex
  get_tx_by_blocknum_and_index_hex $@

## Accounts
elif [ $1 == "get_accounts" ]; then
  shift
  accounts=$(get_accounts)
  check_error
  if [ -z "$1" ]; then
    echo $accounts|jq '.'
  else
    echo $accounts|jq '.['$1']' -r
  fi
elif [ $1 == "get_balance" ]; then
  shift
  addr=$(pad_hex_prefix $1)
  shift
  if [ ! -z "$1" ]; then
    if [ "${1:0:2}" == "0x" ];then
      num=$1
    else
      num=$(to_hex_with_0x_prefix $1)
    fi
    shift
  fi
  # echo "debug get_balance $addr $num"
  balance=$(get_balance $addr $num $@)
  check_error
  # echo "debug get_balance $addr $num --> $balance"
  wei_to_ether $balance
elif [ $1 == "get_tx_count" ]; then
  shift
  addr=$1
  shift
  if [ "$1" == "-h" ]; then
    get_tx_count_by_addr $(pad_hex_prefix $addr) $@|xargs printf "%d\n"
  else
    get_tx_count_by_addr $(pad_hex_prefix $addr) $@
  fi
elif [ $1 == "get_code" ]; then
  shift
  addr=$(pad_hex_prefix $1)
  shift
  if [ ! -z "$1" ]; then
    if [ "${1:0:2}" == "0x" ];then
      num=$1
    else
      num=$(to_hex_with_0x_prefix $1)
    fi
    shift
  fi
  get_code $addr $num
  check_error
elif [ $1 == "get_storage" ]; then
  shift
  addr=$(pad_hex_prefix $1)
  shift
  if [ ! -z "$1" ]; then
    if [ "${1:0:2}" == "0x" ];then
      at=$1
    else
      at=$(to_hex_with_0x_prefix $1)
    fi
    shift
  fi
  if [ ! -z "$1" ]; then
    if [ "${1:0:2}" == "0x" ];then
      num=$1
    else
      num=$(to_hex_with_0x_prefix $1)
    fi
    shift
  fi
  get_storage $addr $at $num
  check_error

## Mining
elif [ $1 == "get_coinbase" ]; then
  shift
  get_coinbase
elif [ $1 == "start_mining" ]; then
  shift
  start_mining
elif [ $1 == "stop_mining" ]; then
  shift
  stop_mining

## INFO & STATUS
elif [ $1 == "status" ] || [ $1 == "info" ] || [ $1 == "get_status" ] || [ $1 == "get_info" ]; then
  shift
  get_status $@

## Execute
elif [ $1 == "eth_compile" ]; then
  shift
  eth_compile "$@"

elif [ $1 == "eth_call" ]; then
  shift
  eth_call $@ |jq -R
  check_error
## DEBUG Moduls
elif [ $1 == "dump_block" ]; then
  shift
  dump_block $@|jq .
  check_error

## UTILS
elif [ $1 == "to_ether" ]; then
  shift
  wei_to_ether $@
  shift
elif [ $1 == "to_wei" ]; then
  shift
  ether_to_wei $@
  shift
elif [ $1 == "list_command" ]; then
  usage
else
  echo "Unkown cmd : $1"
  usage
  exit -1
fi
if [ $DEBUG -gt 0 ]; then
  if [ -e $DEBUG_FILE ]; then
    cat $DEBUG_FILE
    rm $DEBUG_FILE
  fi
fi

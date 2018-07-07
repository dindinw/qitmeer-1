package config

type Config struct {
	HomeDir              string        `short:"A" long:"appdata" description:"Path to application home directory"`
	ShowVersion          bool          `short:"V" long:"version" description:"Display version information and exit"`
	ConfigFile           string        `short:"C" long:"configfile" description:"Path to configuration file"`
	DataDir              string        `short:"b" long:"datadir" description:"Directory to store data"`
	LogDir               string        `long:"logdir" description:"Directory to log output."`
	NoFileLogging        bool          `long:"nofilelogging" description:"Disable file logging."`
	Listeners            []string      `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 8131, testnet: 18131)"`
	MaxPeers             int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	DisableListen        bool          `long:"nolisten" description:"Disable listening for incoming connections"`
	RPCUser              string        `short:"u" long:"rpcuser" description:"Username for RPC connections"`
	RPCPass              string        `short:"P" long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	RPCCert              string        `long:"rpccert" description:"File containing the certificate file"`
	RPCKey               string        `long:"rpckey" description:"File containing the certificate key"`
	RPCMaxClients        int           `long:"rpcmaxclients" description:"Max number of RPC clients for standard connections"`
	DisableRPC           bool          `long:"norpc" description:"Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	DisableTLS           bool          `long:"notls" description:"Disable TLS for the RPC server -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	DisableDNSSeed       bool          `long:"nodnsseed" description:"Disable DNS seeding for peers"`
	DisableCheckpoints   bool          `long:"nocheckpoints" description:"Disable built-in checkpoints.  Don't do this unless you know what you're doing."`
	TxIndex              bool          `long:"txindex" description:"Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	LightNode            bool          `long:"light" description:"start as a nox light node"`
	SigCacheMaxSize      uint          `long:"sigcachemaxsize" description:"The maximum number of entries in the signature verification cache"`
	DumpBlockchain       string        `long:"dumpblockchain" description:"Write blockchain as a flat file of blocks for use with addblock, to the specified filename"`
	TestNet              bool          `long:"testnet" description:"Use the test network"`
	PrivNet              bool          `long:"privnet" description:"Use the private network"`
	DbType               string        `long:"dbtype" description:"Database backend to use for the Block Chain"`
	Profile              string        `long:"profile" description:"Enable HTTP profiling on given [addr:]port -- NOTE port must be between 1024 and 65536"`
	DebugLevel           string        `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical} "`
	DebugPrintOrigins    bool          `long:"printorigin" description:"Print log debug location (file:line) "`
	// Miner
	Generate             bool          `long:"generate" description:"Generate (mine) coins using the CPU"`
	MiningAddrs          []string      `long:"miningaddr" description:"Add the specified payment address to the list of addresses to use for generated blocks -- At least one address is required if the generate option is set"`
	// MemPool Config
	NoRelayPriority      bool          `long:"norelaypriority" description:"Do not require free or low-fee transactions to have high priority for relaying"`
	FreeTxRelayLimit     float64       `long:"limitfreerelay" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
	AcceptNonStd         bool          `long:"acceptnonstd" description:"Accept and relay non-standard transactions to the network regardless of the default settings for the active network."`
	MaxOrphanTxs         int           `long:"maxorphantx" description:"Max number of orphan transactions to keep in memory"`
	MinRelayTxFee        int64
}

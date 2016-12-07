package wallet

var (
	// Paths, TODO: Not hardcode
	walletLDBPath  = "/.factom/m2/wallet/factoid_wallet_ldb"
	walletBoltPath = "/.factom/m2/wallet/factoid_wallet_bolt.db"

	guiLDBPath  = "/.factom/m2/wallet/factoid_gui_ldb"
	guiBoltPath = "/.factom/m2/wallet/factoid_gui_bolt.db"

	txdbLDBPath  = "/.factom/m2/wallet/factoid_blocks_ldb_cache"
	txdbBoltPath = "/.factom/m2/wallet/factoid_blocks_bolt.cache"
)

package wallet

var (
	WalletBoltV1Path = "/.factom/factoid_wallet_bolt.db"

	/* M2 Folder Paths, uses '/wallet/'' */
	walletLDBPath  = "/.factom/wallet/factom_wallet.ldb"
	walletBoltPath = "/.factom/wallet/factom_wallet.db"

	guiLDBPath  = "/.factom/wallet/factoid_gui_ldb.db"
	guiBoltPath = "/.factom/wallet/factom_wallet_gui.db"

	txdbLDBPath  = "/.factom/wallet/factoid_blocks_ldb_cache.db/"
	txdbBoltPath = "/.factom/wallet/factoid_blocks.cache"
)

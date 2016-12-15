package wallet

var (
	WalletBoltV1Path = "/.factom/factoid_wallet_bolt.db"

	/* M2 Fodler Paths, uses '/wallet/'' */
	walletLDBPath  = "/.factom/wallet/factoid_wallet_ldb"
	walletBoltPath = "/.factom/wallet/factom_wallet.db"

	guiLDBPath  = "/.factom/wallet/factoid_gui_ldb"
	guiBoltPath = "/.factom/wallet/factom_wallet_gui.db"

	txdbLDBPath  = "/.factom/wallet/factoid_blocks_ldb_cache"
	txdbBoltPath = "/.factom/wallet/factoid_blocks.cache"
)

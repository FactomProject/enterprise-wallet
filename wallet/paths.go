package wallet

var (
	// Paths, TODO: Not hardcode

	/*  M2 Folder Paths, uses '/m2/wallet' */
	/*
		walletLDBPath  = "/.factom/m2/wallet/factoid_wallet_ldb"
		walletBoltPath = "/.factom/m2/wallet/factoid_wallet_bolt.db"

		guiLDBPath  = "/.factom/m2/wallet/factoid_gui_ldb"
		guiBoltPath = "/.factom/m2/wallet/factoid_gui_bolt.db"

		txdbLDBPath  = "/.factom/m2/wallet/factoid_blocks_ldb_cache"
		txdbBoltPath = "/.factom/m2/wallet/factoid_blocks_bolt.cache"
	*/

	/* M1 Fodler Paths, uses '/wallet/'' */
	walletLDBPath  = "/.factom/wallet/factoid_wallet_ldb"
	walletBoltPath = "/.factom/wallet/factom_wallet.db"

	guiLDBPath  = "/.factom/wallet/factoid_gui_ldb"
	guiBoltPath = "/.factom/wallet/factom_wallet_gui.db"

	txdbLDBPath  = "/.factom/wallet/factoid_blocks_ldb_cache"
	txdbBoltPath = "/.factom/wallet/factoid_blocks.cache"
)

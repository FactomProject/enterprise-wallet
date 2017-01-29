package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	ADD_RANDOM_ADDRESSES = false // Add random addresses. Makes it easier for testing
)

func main() {
	// configure the server
	var (
		guiDB           = flag.String("guiDB", "Bolt", "GUI Database: Bolt, LDB, or Map")
		walDB           = flag.String("walDB", "Bolt", "Wallet Database: Bolt, LDB, or Map")
		txDB            = flag.String("txDB", "Bolt", "Transaction Database: Bolt, LDB, or Map")
		port            = flag.Int("port", 8091, "The port for the GUIWallet")
		compiled        = flag.Bool("compiled", true, "Decides wheter to use the compiled statics or not. Useful for modifying")
		randomAdds      = flag.Bool("randadd", true, "Overrides ADD_RANDOM_ADDRESSES if false and does not add random addresses")
		v1Import        = flag.Bool("i", true, "Search for M1 wallet, if there is no M2 wallet file")
		v1Path          = flag.String("v1path", "/.factom/factoid_wallet_bolt.db", "Change the path for V1 import")
		factomdLocation = flag.String("", "", "Change the location of factomd. Default comes from the config file")

		min = flag.Bool("min", false, "Temporary flag, for testing")
	)
	flag.Parse()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close()
		os.Exit(1)
	}()

	if !(*compiled) {
		COMPILED_STATICS = false
	}

	if *walDB == "Map" {
		if *randomAdds {
			ADD_RANDOM_ADDRESSES = true
		} else {
			ADD_RANDOM_ADDRESSES = false
		}
	}

	if *min {
		FILES_PATH += "min-"
	}

	InitiateWalletAndWeb(*guiDB, *walDB, *txDB, *port, *v1Import, *v1Path, *factomdLocation)
}

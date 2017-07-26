// .
//
// CLI Flags
//
// Enterprise-wallet when launched via CLI has various launch options, see:
//	enterprise-wallet -h
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		factomdLocation = flag.String("factomdlocation", "", "Change the location of factomd. Default comes from the config file")

		min   = flag.Bool("min", false, "Temporary flag, for testing")
		balup = flag.Int64("balup", 10000, "Changes how often the balances of addresses are updated in the cache. Value is in MillSeconds")
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

	if *balup != 10000 {
		BALANCE_UPDATE_INTERVAL = time.Duration(*balup) * time.Millisecond
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

	password := ""
	if *walDB == "ENC" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter password: ")
		text, _, _ := reader.ReadLine()
		password = string(text)
	}

	InitiateWalletAndWeb(*guiDB, *walDB, *txDB, *port, *v1Import, *v1Path, *factomdLocation, password)
}

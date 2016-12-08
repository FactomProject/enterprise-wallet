# M2GUIWallet
GUI Wallet for M2

Must launch from M2GUIWallet directory until static files are compiled into Go.

## Branches to use
 - 'Develop' on everything, 'master' on M2GUIWallet

## To Launch for testing
 - Run 'factomd'
 - Run 'M2GUIWallet', it will populate 5 random factoid and ec addresses, as well as add 1 external.
    - **Must** run from the /M2GUIWallet directory as the web files are not yet compiled in.
    - Just to repeat, **must** run from /M2GUIWallet directory.
 - localhost:8091 to get to wallet


### Testing Notes
 - All databases are configured to be MapDb, so relaunching 'M2GUIWallet' will reset all the data in the wallet.
 - 11 Addresses are preloaded on startup, 5 Random Factoid, 5 random Entry credit, and 1 factoid addresses with factoids in it (in local networks that is)



### Features not working as intended or not Working
  - Import/Export Transactions is **not working**
  - Settings do **not all work**

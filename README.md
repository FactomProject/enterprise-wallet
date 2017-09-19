[![CircleCI](https://circleci.com/gh/FactomProject/enterprise-wallet/tree/develop.svg?style=shield)](https://circleci.com/gh/FactomProject/enterprise-wallet/tree/develop)
[![Coverage Status](https://coveralls.io/repos/github/FactomProject/enterprise-wallet/badge.svg?branch=master)](https://coveralls.io/github/FactomProject/enterprise-wallet?branch=master)

# Enterprise Wallet - GUI Wallet for M2
This uses the same wallet file as factom-walletd and the same port. This means, enterprise-wallet cannot run alongside factom-walletd. enterprise-wallet will import any and all addresses created in the CLI and will monitor any changes the CLI makes and be sure to update itself to reflect those changes. Any addresses created from the CLI however will be marked as not created from the seed, so it is recommended to create all addresses from within the GUI.

Three files are created and used by the wallet:
 1. ~/.factom/wallet/factom_wallet.db
 - ~/.factom/wallet/factom_wallet_gui.db
 - ~/.factom/wallet/factoid_blocks.cache

Database '1' holds all the private keys, this is the main wallet file

Database '2' holds all the nicknames, "seeded" info, and the settings

Database '3' holds every transaction in the Factom blockchain for faster acess for the wallet.

When backing up, backing up #1 is most important. #2 is good to have if you plan on moving to another GUI wallet. #3 does not need to be backed up.


## Branches to use
 - 'Develop' on everything

## To Launch for testing
 - Run 'factomd'
 - Run 'enterprise-wallet'
 - Default, open localhost:8091 in any browser


### Flags
- ```-guiDB=TYPE``` - Gui Database Type, types can be 'Map', 'Bolt', or LDB
  - Default: Bolt (forced to an alternate bolt file when -walDB=ENC)
- ```-walDB=TYPE``` - Wallet Database Type, types can be 'Map', 'Bolt', LDB, or ENC
  - Default: Bolt
- ```-txDB=TYPE``` - Transaction Database Type, Types can be 'Map', 'Bolt', or LDB
  - Default: Bolt
- ```-port=PORT``` - Changes the port the wallet runs on.
  - Default: 8091
- ```-compiled=BOOLEAN``` - Uses statics compiled into GO if true.
  - Default: true
- ```-v1Import=BOOLEAN``` - If true, will look for a V1 database to import. It will only look if there is no M2 database
  - Default: true
- ```-v1Path=PATH_TO_M1``` - The path to look for an M1 wallet.
  - Default: /.factom/factoid_wallet_bolt.db

## Other Flags - Don't bother with these
- ```-randomAdds=BOOLEAN``` - If running on a Map db, this will override adding random addresses on bootup. Put false if you do not want random addresses.
  - Default: true
- ```-min=BOOLEAN``` - If not using compiled in statics, min will decide to use minified versions of the JS and CSS. Reccomend not touching this
  - Default: false

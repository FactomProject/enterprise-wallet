# M2GUIWallet
GUI Wallet for M2

## Branches to use
 - 'Develop' on everything

## To Launch for testing
 - Run 'factomd'
 - Run 'M2GUIWallet'
  - Reccomend using ```-compiled=true``` flag. Will use all static files compiled into Go


### Flags
- ```-guiDB=TYPE``` - Gui Database Type, types can be 'Map', 'Bolt', or LDB
  - Default: Bolt
- ```-walDB=TYPE``` - Wallet Database Type, types can be 'Map', 'Bolt', or LDB
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

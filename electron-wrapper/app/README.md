# Electron Enterprise Wallet Wrapper
This is an electron wrapper for the Enterprise Wallet. It simply makes the wallet a desktop app with a single window.

## To Use

```bash
# Clone this repository
git clone https://github.com/FactomProject/enterprise-wallet.git
# Go into the repository
cd $GOPATH/FactomProject/enterprise-wallet
# Build the Go binary
go build
# Move the Go binary to the bin of the Electron Wrapper
mv enterprise-wallet electron-wrapper/bin/
# Go into electron wrapper directory
cd electron-wrapper
# Install dependencies
npm install
# Run the app
npm start
```
## For Windows users!

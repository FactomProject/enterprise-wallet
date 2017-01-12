rm -r electron-wrapper/bin
mkdir electron-wrapper/bin

echo "Compiling for windows"
env GOOS=windows GOARCH=amd64 go build
mv enterprise-wallet.exe electron-wrapper/bin/

echo "Compiling for Linux"
env GOOS=linux GOARCH=arm go build
mv enterprise-wallet electron-wrapper/bin/enterprise-wallet-lin

echo "Compiling for mac"
env GOOS=darwin GOARCH=amd64 go build
mv enterprise-wallet electron-wrapper/bin/enterprise-wallet-mac
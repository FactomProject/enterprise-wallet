echo "Compiling for windows"
GOOS=windows GOARCH=amd64 go build
# GOOS=windows GOARCH=386 go build -v
# echo "Moving to folder"
# rm windows-binaries/enterprise-wallet.exe
# mv enterprise-wallet.exe windows-binaries/enterprise-wallet.exe

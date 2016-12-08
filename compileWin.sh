echo "Compiling for windows"
GOOS=windows GOARCH=386 go build
# GOOS=windows GOARCH=386 go build -v
echo "Moving to folder"
rm windows-binaries/M2GUIWallet.exe
mv M2GUIWallet.exe windows-binaries/M2GUIWallet.exe

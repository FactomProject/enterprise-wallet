# Release Steps
If you are trying to release cross platform, these steps will not be sufficient and additional dependencies are required.

## Required
 - GoLang Installed
 - Npm Installed
 - NodeJs Installed

# Steps
Most of the steps are already configured, custom configuration will require additional steps
```
# Compile GoLang Binaries for all OS types
cd $GOPATH/src/github.com/FactomProject/enterprise-wallet
# Updates all static files, incase you made a change. It also compiled binaries
# and places the binaries in electron-wrapper/build
sh make.sh
# Time to build the electron-wrapper app
cd electron-wrapper/app
npm install
cd ..
# Build the develop package
npm install
# Build the distribution of YOUR_OS, do not do all of them
# So choose one
# MacOS
npm run dist:mac
# Linux Debian
npm run dist:linux:deb
# Linuc Other
npm run dist:linux:zip
# Windows
npm run dist:win

# Or run all
npm run dist:all

# You distributions will be in /dist
```


# Good links - Custom build options
https://github.com/electron-userland/electron-builder/wiki/Multi-Platform-Build#linux

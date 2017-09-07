# Building Wallet

Building the wallet requires a lot of dependencies to setup for building.

# Setup


## Download Closure compiler

https://developers.google.com/closure/compiler/

Put the jar into `enterprise-wallet/web/closure/compiler.jar`

### Why Closure compiler?

It compiles JS into better JS, and warns you if you fall for any common javascript pitfalls. It's like an automated code review for JS code, so it's pretty cool.


## Install java

To run closure you need java

`sudo apt-get install default-jre`

## Get Staticfiles

This compiles our HTML, CSS, JS, and images into golang .go files. That way we don't need to worry about pathing while serving files as everything is compiled into the binary.

`go get github.com/FactomProject/staticfiles`


## Get Node && NPM

We use electron to package our app as a desktop application, so we will need nodejs and npm. (npm comes with node)

https://docs.npmjs.com/getting-started/installing-node


## Run make.sh

This script will compile our go into windows, mac, and linux binaries. They are placed into `enterprise-wallet/electron-wrapper/bin`

```
cd $GOPATH/src/github.com/FactomProject/enterprise-wallet
bash make.sh closure
```

## Install Node Packages

```
cd electron-wrapper/app
npm install

cd ..
npm install
```

## Dependencies for Electron build


For linux you must get
```
icns2png (icnsutils)
gm (graphicsmagick)
```

For windows cross compiling you need
```
wine 1.8+
```

For mac, you cannot easily cross compile


# Building App

If you ran all the above steps, we already have the binaries for the wallet built and ready to be packaged. We use a tool called Electorn Builder to actually package the app. To run this:

```
cd $GOPATH/src/github.com/FactomProject/enterprise-wallet/electron-wrapper
npm run dist:win
npm run dist:linux:deb
npm run dist:linux:zip
npm run dist:mac

# Or for all
npm run dist:all
```

Try to fix any errors that crop up. If you are cross compiling from linux, you cannot create a .dmg for the mac build. It's unfortunate for sure, so building from a mac is reccomended, as that can build for all OS's. If you only have a linux though, the mac build will say it failed, but it will actually have built the .app, and failed on the .dmg.

To see your newly build packages:

```
ls $GOPATH/src/github.com/FactomProject/enterprise-wallet/electron-wrapper/dist
```


# Good links - Custom build options

https://github.com/electron-userland/electron-builder/wiki/Multi-Platform-Build#linux
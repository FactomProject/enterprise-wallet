const electron = require('electron')
var ps = require('ps-node');
const {dialog} = require('electron')
// Module to control application life.
const app = electron.app
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow

const path = require('path')
const url = require('url')

var exec = require('child_process').exec;

// Detect if windows
var isWin = /^win/.test(process.platform);

// For deployment
const PATH_TO_BIN = "../app.asar.unpacked/bin/"
// For local testing
// const PATH_TO_BIN = "bin/"

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow
let secondaryWindow
let walletd

// Start own enterprise-wallet daemon
const startOwn = true

WALLETD_UP = false

const PORT_TO_SERVE = "8091"
function execWalletd() {
  if(!startOwn || WALLETD_UP){
    return
  }
  console.log("Executing enterprise-wallet...")

  if(isWin){
    walletd = exec(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet.exe -port=' + PORT_TO_SERVE), function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {

      } else {
        console.log("Running as Windows OS")
        mainWindow.loadURL('http://localhost:' + PORT_TO_SERVE + '/');
      }
    });
  } else {
    walletd = exec(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet -port=' + PORT_TO_SERVE), function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {

      } else {
        console.log("Running as Linux/Mac OS")
        mainWindow.loadURL('http://localhost:' + PORT_TO_SERVE + '/');
      }
    });
  }
}

// Before we launch, we need to check if we already have the app running, then check if
// enterprise-wallet has been left hanging around
function startApp(){
  console.log("Checking for app already open...")
  ps.lookup({
    command: 'EnterpriseWallet'
    }, function(err, resultList ) {
    if(resultList.length > 1){
      // App already running, we don't want to spawn two
      dialog.showErrorBox("App Already Running", "EnterpriseWallet is already running, please close it before trying to open it again.")
      app.quit()
    } else {
      // Look for hanging golang process and kill them
      console.log("Checking for hanging enterprise-wallet process...")
      ps.lookup({
        command: 'enterprise-wallet'
        }, function(err, resultList ) {
        if (err) {
          throw new Error( err );
        }

        resultList.forEach(function( process ){
          if( process ){
            if(process.command.endsWith("enterprise-wallet")) {
              console.log( 'Killing PID: %s, COMMAND: %s, ARGUMENTS: %s', process.pid, process.command, process.arguments );
              ps.kill(process.pid, function(err){
                if (err) {
                  console.log( "enterprise-wallet is running, and cannot be stopped: " + err );
                }
              });
            }
          }
        });

        // Now we can start our processes and app
        console.log("Launching...")
        execWalletd()
        WALLETD_UP = true
        // Clear cache always, makes updates easier
        createWindow()
        deleteChromeCache()
      });
    }
  });
}

function createWindow () {


  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1400, 
    height: 800,
    minWidth: 600,
    minHeight: 600,
    center: true
  })

  mainWindow.loadURL('http://localhost:' + PORT_TO_SERVE + '/');

  // Open the DevTools.
  //mainWindow.webContents.openDevTools()

  // Emitted when the window is closed.
  mainWindow.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null
  })
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', startApp)

// Quit when all windows are closed.
app.on('window-all-closed', function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q

  if (process.platform !== 'darwin') {
    if(startOwn){
      walletd.stdin.pause();
      walletd.kill();
    }
    app.quit()
  }
})

app.on('quit', function(){
  if(WALLETD_UP) {
    walletd.stdin.pause();
    walletd.kill();
  }
})

app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow()
    execWalletd()
  }
})

var deleteChromeCache = function() {
  mainWindow.webContents.session.clearCache(function(){});
};

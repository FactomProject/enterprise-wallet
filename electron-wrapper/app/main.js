const electron = require('electron')
var ps = require('ps-node');
var request=require('request');
const {dialog} = require('electron').dialog
// Module to control application life.
const app = electron.app
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow

const path = require('path')
const url = require('url')

var exec = require('child_process').exec;
var spawn = require('child_process').spawn;

// Detect if windows
var isWin = /^win/.test(process.platform);

// For deployment
const PATH_TO_BIN = "../app.asar.unpacked/bin/"
// For local testing
// const PATH_TO_BIN = "bin/"

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow
let loadingWindow
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
   /*walletd = exec(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet.exe -port=' + PORT_TO_SERVE), function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {
        console.log(error)
      }
    });*/
    walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet.exe'),['-port=' + PORT_TO_SERVE])
  } else {
    /*walletd = exec(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet -port=' + PORT_TO_SERVE), function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {
        console.log(error)
      } 
    });*/
    walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet'),['-port=' + PORT_TO_SERVE])
  }
    
  runWhenWalletUp(function(){
    loadMainWindow()
    // Clear cache always, makes updates easier
    deleteChromeCache()
  })
}

// Runs when the wallet is able to start serving web pages
function runWhenWalletUp(callback){
  request.get('http://localhost:' + PORT_TO_SERVE + "/GET?request=on" ,function(err,res,body){
    if (!err && res.statusCode == 200) {
      callback()
    } else {
      console.log("Not up yet, trying again in 0.5seconds...")
      setTimeout(function(){
        runWhenWalletUp(callback)
      }, 500)
    }
  });
}

/* Single Instance Check */
var iShouldQuit = app.makeSingleInstance(function(commandLine, workingDirectory) {
  if (mainWindow) {
    if (mainWindow.isMinimized()) mainWindow.restore();
    mainWindow.show();
    mainWindow.focus();
  }
  return true;
});
if(iShouldQuit){app.quit();return;}

function cleanUp(functionAfterCleanup) {
  console.log("Killing enterprise-wallet processes")
  var commandToKill = "enterprise-wallet"
  if(isWin) {
    commandToKill = "enterprise-wallet.exe"
  }
  ps.lookup({
    command: commandToKill
    }, function(err, resultList ) {
    if (err) {
      console.log(err);
    }

    if(resultList.length == 0) {
      functionAfterCleanup()
    } else {
      resultList.forEach(function( process ){
        if( process ){
          if(process.command.endsWith(commandToKill)) {
            console.log( 'Killing PID: %s, COMMAND: %s, ARGUMENTS: %s', process.pid, process.command, process.arguments );
            ps.kill(process.pid, function(err){
              if (err) {
                console.log( "enterprise-wallet is running, and cannot be stopped: " + err );
              }
            });
          }
        }
      });

      functionAfterCleanup()
    }
  });
}

// Before we launch, we need to check if we already have the app running, then check if
// enterprise-wallet has been left hanging around
function startApp(){
  // Look for hanging golang process and kill them
  createLoadingWindow()
  console.log("Checking for hanging enterprise-wallet process...")
  cleanUp(function(){
    // Now we can start our processes and app
    execWalletd()
    WALLETD_UP = true
  })
}

function loadMainWindow() {
  if(loadingWindow !== null) {
    loadingWindow.close()
    loadingWindow === null
  }
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1400, 
    height: 800,
    minWidth: 600,
    minHeight: 600,
    center: true,
    title: 'EnterpriseWallet'
  })

  // Load loading window
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

function createLoadingWindow() {
  // Create the browser window.
  loadingWindow = new BrowserWindow({
    width: 700, 
    height: 200,
    center: true,
    autoHideMenuBar: true,
    titleBarStyle: 'hidden',
    resizable: false,
    title: 'EnterpriseWallet',
    // transparent: true
  })

  // Load loading window
  console.log("Showing a loading...")
  loadingWindow.loadURL(url.format({
    pathname: path.join(__dirname, 'loading.html'),
    protocol: 'file:',
    slashes: true
  }))

  // Open the DevTools.
  //mainWindow.webContents.openDevTools()

  // Emitted when the window is closed.
  loadingWindow.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    loadingWindow = null
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
    console.log("Properly exiting...")
    WALLETD_UP = false
    cleanUp(function(){app.quit()})
  }
})

// App close handler
app.on('before-quit', function() {
  if(WALLETD_UP) {
    walletd.kill();
  }
});

app.on('quit', function(){
})

app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    loadMainWindow()
    execWalletd()
  }
})

var deleteChromeCache = function() {
  mainWindow.webContents.session.clearCache(function(){});
};
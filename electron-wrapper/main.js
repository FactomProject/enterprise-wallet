const electron = require('electron')
const ipc = require('electron').ipcMain
// Module to control application life.
const app = electron.app
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow

const path = require('path')
const url = require('url')

var exec = require('child_process').exec;

// Detect if windows
var isWin = /^win/.test(process.platform);

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow
let secondaryWindow
let walletd

// Start own enterprise-wallet daemon
const startOwn = true

ipc.on('asynchronous-message', function (event, arg) {
  event.sender.send('asynchronous-reply', 'pong')
})

WALLETD_UP = false

function execWalletd() {
  if(!startOwn){
    return
  }

  if(isWin){
    walletd = exec('./bin/enterprise-wallet.exe -txDB=Map', function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {

      } else {
        console.log("Running as Windows OS")
        mainWindow.loadURL('http://localhost:8091/');
      }
    });
  } else {
    walletd = exec('./bin/enterprise-wallet -txDB=Map', function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {

      } else {
        console.log("Running as Linux/Mac OS")
        mainWindow.loadURL('http://localhost:8091/');
      }
    });
  }
}

function createWindow () {
  execWalletd()

  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1400, 
    height: 800,
    minWidth: 600,
    minHeight: 600,
    center: true
  })

  mainWindow.loadURL('http://localhost:8091/');

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
app.on('ready', createWindow)

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

app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow()
    execWalletd()
  }
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
var deleteChromeCache = function() {
    var chromeCacheDir = path.join(app.getPath('userData'), 'Cache'); 
    if(fs.existsSync(chromeCacheDir)) {
        var files = fs.readdirSync(chromeCacheDir);
        for(var i=0; i<files.length; i++) {
            var filename = path.join(chromeCacheDir, files[i]);
            if(fs.existsSync(filename)) {
                try {
                    fs.unlinkSync(filename);
                }
                catch(e) {
                    console.log(e);
                }
            }
        }
    }
};

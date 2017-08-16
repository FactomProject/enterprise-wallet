const electron = require('electron')
const {ipcMain} = require('electron')
const dialog = require('electron').dialog
var ps = require('ps-node');
var request=require('request');
// Module to control application life.
const app = electron.app
// Module to create menu for copy/paste
var Menu = electron.Menu
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow

const path = require('path')
const url = require('url')

require('electron-context-menu')({
    prepend: (params, browserWindow) => [{
        label: 'Rainbow',
        // only show it when right-clicking images 
        visible: params.mediaType === 'image'
    }]
});

// var exec = require('child_process').exec;
var spawn = require('child_process').spawn;

// Detect if windows
var isWin = /^win/.test(process.platform);

const isDev = require('electron-is-dev');

// For production
var PATH_TO_BIN = "../app.asar.unpacked/bin/"
if (isDev) {
  // Override for development
  console.log('Running in development');
  var PATH_TO_BIN = "bin/"
} else {
  console.log('Running in production');
}

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow
let loadingWindow
let choiceWindow
let walletd

// Start own enterprise-wallet daemon
const startOwn = true

WALLETD_UP = false

const PORT_TO_SERVE = "8091"
function execWalletd(password) {
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
    console.log("Running as Windows")
    if(password === "") {
      console.log("Running without Encryption")
      walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet.exe'),[])
    } else {
      console.log("Running with Encryption")
      walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet.exe'),["-walDB=ENC"])
    }
  } else {
    /*walletd = exec(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet -port=' + PORT_TO_SERVE), function callback(error, stdout, stderr){
      console.log(stdout)
      if (error !== null) {
        console.log(error)
      } 
    });*/
    console.log("Running as Mac/Linux")
    if(password === "") {
      console.log("Running without Encryption")
      walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet'),[])
    } else {
      console.log("Running with Encryption")
      walletd = spawn(path.join(__dirname, PATH_TO_BIN + 'enterprise-wallet'),["-walDB=ENC"])
    }
  }

  walletd.stdout.on('data', function(data) {
    var s = data.toString()
    
    console.log(s)
    // Look for password prompt
    if(s.includes("password")){
      console.log("Found prompt, inputting password to wallet")
      walletd.stdin.setEncoding('utf-8');
      walletd.stdin.write(password + "\n");
      //walletd.stdin.close()
    } 

    if(s.includes("Error in starting wallet")) {
      var errormessage = ""
      if(s.includes("message authentication failed")) {
        errormessage = "The password given to unlock the encrypted database was incorrect. "+
         "Please try launching the wallet again with the correct password. If you feel this is an " +
         "error, please reach out on our Factom Community slack."
      } else {
        var n = s.indexOf("Error in starting wallet");
        errormessage = "There was an error launching the EnterpriseWallet, but that reason was not" +
        "able to be deducted. Below is the error message that was generated.\n\n"+ s.substring(n,s.length)
      }
      sendMsgToLoading(errormessage, function() {
        ChooseWalletType(true)
        loadingWindow.close()
        loadingWindow = null
      }, 2000)
      //dialog.showErrorBox('Error Launching EnterpriseWallet', errormessage)
      return
      //app.quit()
    }
    if(s.includes("Starting GUI")){
      sendMsgToLoading("success", function(){
        runWhenWalletUp(function(){
        loadMainWindow()
          // Clear cache always, makes updates easier
           deleteChromeCache()
        })  
      }, 2000)
    }
  });
}

function sendMsgToLoading(mesg, f, wait) {
  console.log(mesg)
  loadingWindow.webContents.send('info' , {msg:mesg})
  if(f !== undefined) {
    setTimeout(function(){
      f()
    }, wait)
  }
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
if(iShouldQuit && (isWin || process.platform === 'darwin')){app.quit();return;}

function cleanUp(functionAfterCleanup) {
  console.log("Killing enterprise-wallet processes")
  var commandToKill = "enterprise-wallet"
  var sep = "/"

  ps.lookup({
    command: commandToKill
    }, function(err, resultList ) {
    if (err) {
      console.log(err);
    }
    if(resultList === undefined || resultList === null || resultList.length == 0) {
      functionAfterCleanup()
    } else {
      resultList.forEach(function( process ){
        if( process ){
          // The command can come in the form of '/opt/EnterpriseWallet/enterprise-wallet -port=8091'
          // We need to first strip all pathing: enterprise-wallet -port=8091
          // Then test if we start with 'enterprise-wallet'
          rawProcessCommand = process.command
          cmds = rawProcessCommand.split(sep)
          processCommand = cmds[cmds.length - 1]
          if(processCommand.startsWith(commandToKill)) {
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

function ChooseWalletType(witherror) {
  // Create the browser window.
  choiceWindow = new BrowserWindow({
    width: 1000, 
    height: 500,
    center: true,
    autoHideMenuBar: true,
    titleBarStyle: 'hidden',
    resizable: false,
    title: 'EnterpriseWallet',
    // transparent: true
  })

  var pathurl = 'loading/index.html'
  var ext = ''
  if(witherror) {
    ext = `?error=Wrong Password`
    pathurl = 'loading/wallet-secure.html'
  }

  // Load choice window
  console.log("Showing a choice...")
  choiceWindow.loadURL(url.format({
    pathname: path.join(__dirname, pathurl),
    protocol: 'file:',
    slashes: true
  }) + ext)

  // Open the DevTools.
  //mainWindow.webContents.openDevTools()

  // Emitted when the window is closed.
  choiceWindow.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    choiceWindow = null
  })
}

// Before we launch, we need to check if we already have the app running, then check if
// enterprise-wallet has been left hanging around
function startApp(){
  ChooseWalletType(false)
}

function loadMainWindow() {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1400, 
    height: 800,
    minWidth: 600,
    minHeight: 600,
    center: true,
    title: 'EnterpriseWallet'
  })

  if(loadingWindow !== null) {
    loadingWindow.close()
    loadingWindow === null
  }

  console.log("Main window is now open, loading window is closed.")

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

  // Create the Application's main menu
  var template = [{
      label: "Application",
      submenu: [
          { label: "About Application", selector: "orderFrontStandardAboutPanel:" },
          { type: "separator" },
          { label: "Quit", accelerator: "CmdOrCtrl+Q", click: function() { app.quit(); }}
      ]}, {
      label: "Edit",
      submenu: [
          { label: "Undo", accelerator: "CmdOrCtrl+Z", selector: "undo:" },
          { label: "Redo", accelerator: "Shift+CmdOrCtrl+Z", selector: "redo:" },
          { type: "separator" },
          { label: "Cut", accelerator: "CmdOrCtrl+X", selector: "cut:" },
          { label: "Copy", accelerator: "CmdOrCtrl+C", selector: "copy:" },
          { label: "Paste", accelerator: "CmdOrCtrl+V", selector: "paste:" },
          { label: "Select All", accelerator: "CmdOrCtrl+A", selector: "selectAll:" }
      ]}, {
      label: "View",
      submenu: [
        { label: "Reload", accelerator: "CmdOrCtrl+R", role: "reload" },
        { type: "separator" },
        { label: "Reset Zoom", role: "resetzoom" },
        { label: "Zoom In", accelerator: 'CmdOrCtrl+=', role: "zoomin" },
        { label: "Zoom Out", accelerator: "CmdOrCtrl+-", role: "zoomout" },
      ]},
      {
        role: 'window',
        submenu: [
          {
            role: 'minimize'
          },
          {
            role: 'close'
          }
        ]
      }
  ];

  if(isWin) {
    template[0].submenu = [
          { label: "Quit", accelerator: "Command+Q", click: function() { app.quit(); }}
      ]
  }

  Menu.setApplicationMenu(Menu.buildFromTemplate(template));
}

function createLoadingWindow() {
  // Create the browser window.
  loadingWindow = new BrowserWindow({
    width: 1000, 
    height: 500,
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
    pathname: path.join(__dirname, 'loading/loading.html'),
    protocol: 'file:',
    slashes: true
  }))

  if(choiceWindow !== null) {
    choiceWindow.close()
    choiceWindow === null
  }

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

// Recieve the start options and start the wallet
ipcMain.on('submitForm', function(event, data) {
  // Look for hanging golang process and kill them
  createLoadingWindow()
  console.log("Checking for hanging enterprise-wallet process...")
  // Cleanup takes incredibly long on windows, so on bootup it makes things hang and be slow.
  // It is a safeguard to do at launch, and not required. We cleanup on close, so if a user closes
  // incorrectly, we will get a hanging process. They will have to launch, then close properly to
  // clean up the haning processes
  if(isWin) {
    setTimeout(function(){execWalletd(data)}, 300)
    WALLETD_UP = true
  } else {
    cleanUp(function(){
      // Now we can start our processes and app
      setTimeout(function(){
        execWalletd(data)
        WALLETD_UP = true
      }, 300)
    })
  }
});

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
    if(WALLETD_UP) {
      walletd.stdin.pause();
      walletd.kill();
      WALLETD_UP = false
    }
    cleanUp(function(){app.quit()})
  }
})

// App close handler
app.on('before-quit', function() {
  if(WALLETD_UP) {
    walletd.stdin.pause();
    walletd.kill();
    cleanUp(function(){app.quit()})
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
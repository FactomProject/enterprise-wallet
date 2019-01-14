@echo off
REM An attempt at a build file for Windows
REM Yes, yes, I know

SETLOCAL
set wallDir=%GOPATH%/src/github.com/FactomProject/enterprise-wallet

echo Emulating /web/compile.sh
cd %wallDir%/web


REM concatJs.sh
REM ###################################################################

echo JS/CSS
echo   Concatenating JS files into all.js

cd %wallDir%/web/statics/js
del all.js
copy addressbook.js + index.js + send-convert.js + settings.js + recieve-factoids.js + new-address.js + edit-address.js + backup.js all.js /b >NUL


REM web/min-file.sh closure
REM ###################################################################

cd %wallDir%/web/
echo   Removing existing temporary directory
rmdir /Q /S min-statics >NUL 2>NUL
rmdir /Q /S min-templates >NUL 2>NUL

echo   Copying folders to temp destination
robocopy statics min-statics\ /E >NUL
robocopy templates min-templates\ /E >NUL



REM web/non-essentials.sh
REM ###################################################################
echo   (web/non-essential.sh) Removing Non-Essentials
rmdir min-statics\scss /s /q
del min-statics\bower_components\foundation-sites\dist\foundation.css
del min-statics\bower_components\foundation-sites\dist\foundation.js
rmdir min-statics\bower_components\foundation-sites\js /s /q
rmdir min-statics\bower_components\foundation-sites\scss /s /q
del min-statics\bower_components\foundation-sites\README.md
del min-statics\bower_components\foundation-sites\foundation-sites.scss
del min-statics\bower_components\foundation-sites\LICENSE
del min-statics\bower_components\foundation-sites\docslink.sh
del min-statics\bower_components\foundation-sites\bower.json
del min-statics\bower_components\foundation-sites\.bower.json

REM    jquery
rmdir min-statics\bower_components\jquery\src /s /q
del min-statics\bower_components\jquery\dist\jquery.js
del min-statics\bower_components\jquery\dist\jquery.min.map
del min-statics\bower_components\jquery\.bower.json
del min-statics\bower_components\jquery\bower.json
del min-statics\bower_components\jquery\MIT-LICENSE.txt

REM    motion-ui
rmdir min-statics\bower_components\motion-ui\docs /s /q
rmdir min-statics\bower_components\motion-ui\src /s /q
del min-statics\bower_components\motion-ui\dist\motion-ui.css
del min-statics\bower_components\motion-ui\dist\motion-ui.js
del min-statics\bower_components\motion-ui\.bower.json
del min-statics\bower_components\motion-ui\bower.json
del min-statics\bower_components\motion-ui\LICENSE
del min-statics\bower_components\motion-ui\motion-ui.scss
del min-statics\bower_components\motion-ui\package.json
del min-statics\bower_components\motion-ui\README.md

REM    what-input
del min-statics\bower_components\what-input\.bower.json
del min-statics\bower_components\what-input\bower.json
del min-statics\bower_components\what-input\demo.html
del min-statics\bower_components\what-input\LICENSE
del min-statics\bower_components\what-input\package.json
del min-statics\bower_components\what-input\what-input.js



echo   Minfying templates...
echo     Minifying all.js using closure...
java -jar closure/compiler.jar  --js_output_file=min-statics/js/all.js statics/js/all.js
echo     Minifying ajax.js using closure...
java -jar closure/compiler.jar  --js_output_file=min-statics/js/ajax.js statics/js/ajax.js
echo     Minifying app.js using closure...
java -jar closure/compiler.jar  --js_output_file=min-statics/js/app.js statics/js/app.js

copy statics\css\app.css min-statics\css\app.css /Y >NUL
copy statics\css\other.css min-statics\css\other.css /Y >NUL

echo   Compiling statics into GO....
staticfiles -o files\statics\statics.go min-statics
staticfiles -o files\templates\templates.go min-templates

REM compileAll.sh

cd %wallDir%
rmdir electron-wrapper\bin /q /s
mkdir electron-wrapper\bin

echo Compiling for windows
go build
move enterprise-wallet.exe electron-wrapper\bin\ >NUL
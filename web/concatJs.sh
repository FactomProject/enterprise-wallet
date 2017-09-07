#!/bin/bash

rootPath=$GOPATH/src/github.com/FactomProject/enterprise-wallet/web/

# Concatenate all js files
allFile="${rootPath}statics/js/all.js"

rm $allFile

touch $allFile

cat ${rootPath}statics/js/addressbook.js >> $allFile
cat ${rootPath}statics/js/index.js >> $allFile
cat ${rootPath}statics/js/send-convert.js >> $allFile
cat ${rootPath}statics/js/settings.js >> $allFile
cat ${rootPath}statics/js/recieve-factoids.js >> $allFile
cat ${rootPath}statics/js/new-address.js >> $allFile
cat ${rootPath}statics/js/edit-address.js >> $allFile
cat ${rootPath}statics/js/backup.js >> $allFile
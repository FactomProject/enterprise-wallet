#!/bin/bash

# Concatenate all js files
allFile="statics/js/all.js"

rm $allFile

touch $allFile

cat statics/js/addressbook.js >> $allFile
cat statics/js/index.js >> $allFile
cat statics/js/send-convert.js >> $allFile
cat statics/js/settings.js >> $allFile
cat statics/js/recieve-factoids.js >> $allFile
cat statics/js/new-address.js >> $allFile
cat statics/js/edit-address.js >> $allFile
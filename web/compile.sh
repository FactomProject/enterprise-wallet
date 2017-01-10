#!/bin/bash

#
sh concatJs.sh
# Minifies static files
if [ "$1" == "closure" ] # Closure compile by google
	then
	sh min-file.sh closure
else 
	sh min-file.sh
fi

echo "Compiling statics into GO...."
# Compiles static files into binary
staticfiles -o files/statics/statics.go min-statics
staticfiles -o files/templates/templates.go min-templates
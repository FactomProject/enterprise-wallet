#!/bin/bash
# Runs the compile.sh in web/
# Then runs go install

cd web/
# Minifies JS/CSS
# Compiles into binary
#     The input can be "closure", which will use google's closure compile for JS.
#     It is slower to minify, but better results.
sh compile.sh $1

cd ..

# Build binaries
sh compileAll.sh
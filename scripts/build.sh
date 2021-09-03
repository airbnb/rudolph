#!/usr/bin/env bash
set -e

# Figure out the project's basepath
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Establish directories
BUILD_DIR=$DIR/build
LINUX_BUILD_DIR=$BUILD_DIR/linux
MACOS_BUILD_DIR=$BUILD_DIR/macos
APPS_DIR=$DIR/cmd
CLI_NAME=rudolph
PKG_DIR=$BUILD_DIR/package
DEPLOYMENT_ZIP_PATH=$PKG_DIR/deployment.zip

cd "$DIR"

# Cleanup from previous run(s)
rm -rf $BUILD_DIR

# Do the build things
echo "*** compiling application binaries... ***"

echo "  compiling api..."
GOOS=linux GOARCH=amd64 go build -o $LINUX_BUILD_DIR/api $APPS_DIR/api

echo "  compiling authorizer..."
GOOS=linux GOARCH=amd64 go build -o $LINUX_BUILD_DIR/authorizer $APPS_DIR/authorizer

echo "  compiling cli..."
GOOS=darwin go build -o $MACOS_BUILD_DIR/cli $APPS_DIR/cli

ln -sf $MACOS_BUILD_DIR/cli $DIR/$CLI_NAME

echo "*** packaging... ***"

mkdir $PKG_DIR
# remember zip 2nd arg zips up all the specified directories. We omit the dir info by cd,
# but you could use the -j option as well.
cd $LINUX_BUILD_DIR; zip -r $DEPLOYMENT_ZIP_PATH *

echo "*** complete ***"

echo "  created:"
echo "    $DEPLOYMENT_ZIP_PATH"
echo "    $DIR/$CLI_NAME"

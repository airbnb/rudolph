#!/usr/bin/env bash
set -e

# Figure out the project's basepath
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Establish directories
BUILD_DIR=$DIR/build
LINUX_BUILD_DIR=$BUILD_DIR/linux
LINUX_BUILD_DIR_API=$LINUX_BUILD_DIR/api
LINUX_BUILD_DIR_AUTHORIZER=$LINUX_BUILD_DIR/authorizer
MACOS_BUILD_DIR=$BUILD_DIR/macos
APPS_DIR=$DIR/cmd
CLI_NAME=rudolph
CLI_BUILD_DIR=$BUILD_DIR/cli
PKG_DIR=$BUILD_DIR/package
API_DEPLOYMENT_ZIP_PATH=$PKG_DIR/api_deployment.zip
API_AUTHORIZER_DEPLOYMENT_ZIP_PATH=$PKG_DIR/api_authorizer_deployment.zip

cd "$DIR"

# Cleanup from previous run(s)
rm -rf $BUILD_DIR

# Do the build things
echo "*** compiling application binaries... ***"

echo "  compiling api in linux:arm64..."
GOOS=linux GOARCH=arm64 go build -o $LINUX_BUILD_DIR_API/bootstrap $APPS_DIR/api

echo "  compiling authorizer in linux:arm64..."
GOOS=linux GOARCH=arm64 go build -o $LINUX_BUILD_DIR_AUTHORIZER/bootstrap $APPS_DIR/authorizer

if [ "$(uname)" == "Darwin" ]; then
    echo "  compiling cross-compatible macOS cli..."
    GOOS=darwin GOARCH=amd64 go build -o $MACOS_BUILD_DIR/cli_amd64 $APPS_DIR/cli
    GOOS=darwin GOARCH=arm64 go build -o $MACOS_BUILD_DIR/cli_arm64 $APPS_DIR/cli
    lipo -create -output $MACOS_BUILD_DIR/cli $MACOS_BUILD_DIR/cli_amd64 $MACOS_BUILD_DIR/cli_arm64
    ln -sf $MACOS_BUILD_DIR/cli $DIR/$CLI_NAME
else
    echo "  compiling cli..."
    go build -o $CLI_BUILD_DIR/cli $APPS_DIR/cli
    ln -sf $CLI_BUILD_DIR/cli $DIR/$CLI_NAME
fi

echo "*** packaging... ***"

mkdir $PKG_DIR
# remember zip 2nd arg zips up all the specified directories. We omit the dir info by cd,
# but you could use the -j option as well.
cd $LINUX_BUILD_DIR_API; zip -r $API_DEPLOYMENT_ZIP_PATH *
cd $LINUX_BUILD_DIR_AUTHORIZER; zip -r $API_AUTHORIZER_DEPLOYMENT_ZIP_PATH *

echo "*** complete ***"

echo "  created:"
echo "    API: $API_DEPLOYMENT_ZIP_PATH"
echo "    API Authorizer: $API_AUTHORIZER_DEPLOYMENT_ZIP_PATH"
if [ "$(uname)" == "Darwin" ]; then
    echo "    generated cross-compiled macOS cli"
    echo "    CLI: $MACOS_BUILD_DIR/cli"
else
    echo "    CLI: $CLI_BUILD_DIR/cli"
fi

#!/usr/bin/env bash
set -e

# Figure out the project's basepath
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Establish directories
TF_ENV_DIR=$DIR/deployments/environments
EXAMPLE_ENV_DIR=$TF_ENV_DIR/example

if [ $# -eq 0 ]
then
  echo "Not enough arguments"
  exit 1
fi

NEW_ENV=$1
NEW_ENV_DIR=$TF_ENV_DIR/$NEW_ENV
echo "new env dir:"
echo $NEW_ENV_DIR


if [[ -d $NEW_ENV_DIR ]]
then
  echo "Warning: environment $NEW_ENV already exists. Exiting."
  exit 1
fi

# Now create the directory and stuff
mkdir $NEW_ENV_DIR
cd $NEW_ENV_DIR

# Create a .gitignore file which just ignores this current directory
echo "*" >> .gitignore

ln -s ../../terraform_modules/default_main.tf ./main.tf
ln -s ../../terraform_modules/default_variables.tf ./variables.tf

cp $EXAMPLE_ENV_DIR/_backend.tf ./
cp $EXAMPLE_ENV_DIR/_outputs.tf ./
cp $EXAMPLE_ENV_DIR/config.auto.tfvars.json ./
cp $EXAMPLE_ENV_DIR/versions.tf ./

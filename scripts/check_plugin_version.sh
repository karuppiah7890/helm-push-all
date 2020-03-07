#!/bin/sh

set -e

plugin_yaml=$1
expected_version=$2

if [ ! -e $plugin_yaml ];
then
    echo "File $plugin_yaml does not exist!"
    exit 1
fi

plugin_version=$(yq read $plugin_yaml version)

if [ "$plugin_version" != "$expected_version" ];
then
    echo "Plugin version $plugin_version is not the expected version $expected_version"
    exit 1
fi
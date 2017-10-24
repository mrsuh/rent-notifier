#!/bin/sh

cd "$( cd `dirname $0` && pwd )/.."

for f in  'config/config.yml'
do
    if [ ! -f $f ]; then
        cp $f.dist $f
        echo "File created from $f"
    fi
done
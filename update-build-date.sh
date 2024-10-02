#!/bin/sh

DATETIME=$(date +'%Y%m%d%H%M%S')

# Replace versionStamp in local dev file
sed -i.previous -e "s/VERSION_STAMP=[0-9]*/VERSION_STAMP=$DATETIME/" .env
rm .env.previous

# Replace VERSION_STAMP in app.yaml
#sed -i.previous -e "s/VERSION_STAMP: \"[0-9]*\"/VERSION_STAMP: \"$DATETIME\"/" app.yaml
#rm app.yaml.previous

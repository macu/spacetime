#!/bin/sh

npm run dev || { echo 'Client code compilation failed.' ; exit 1; }

sh update-build-date.sh

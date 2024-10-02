#!/bin/sh

sh update-build-date.sh

npm run dev || { echo 'Client code compilation failed.' ; exit 1; }

sh ./bin/restart-web.sh

# Spacetime

## Install Docker and Docker Compose

This command runs on Mac.

```bash
brew install --cask docker
```

## Start up containers

```bash
cd /treetime
sh ./bin/restart-all.sh
```

## Access postgres interactive shell

Run `sql/init.pgsql` by copying contents into shell.

```bash
sh ./bin/psql-shell.sh
```

## Rebuild web app

```bash
sh ./bin/restart-web.sh
```

## Test in browser

http://localhost:8080/

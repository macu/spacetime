#!/bin/bash

# bring down and back up
docker compose down
docker compose up -d --build

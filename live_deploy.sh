#!/bin/bash

DEPLOY_ENV="live"
export DEPLOY_ENV

git checkout main
git pull origin main

go build -o main ./cmd
chmod a+x ./main
./main
#!/bin/bash

DEPLOY_ENV="test"
export DEPLOY_ENV

git checkout develop
git pull origin develop

go build -o main ./cmd
chmod a+x ./main
./main
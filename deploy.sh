#!/bin/bash

git checkout develop
git pull origin develop

go build -o main ./cmd
chmod a+x ./main
./main
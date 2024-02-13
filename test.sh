#!/bin/bash

current_directory=$(pwd)
export workdirectory="$current_directory"

go test -v ./...
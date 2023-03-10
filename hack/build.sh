#!/usr/bin/env bash
set -x

go generate
cat version.txt
go install

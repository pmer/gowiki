#!/bin/sh
go get \
    github.com/DataDog/zstd \
    github.com/mattn/go-sqlite3

go build

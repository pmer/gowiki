#!/bin/sh
rm -f gowiki.db
./tree-sap &

sleep 1

curl --data-binary @tree-sap.go localhost:8080
curl localhost:8080 >tree-sap.go.copy
cmp tree-sap.go tree-sap.go.copy

#! /bin/bash
export PATH=$PWD/../../fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$PWD/../../fabric-samples/config/
peer lifecycle chaincode package docstore.tar.gz --path ./ --lang golang --label docstore_1.0
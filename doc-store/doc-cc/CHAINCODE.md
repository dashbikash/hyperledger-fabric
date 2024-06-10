# CHAINCODE

**Chaincode Life Cycle**

- Package the chaincode
- Install the chaincode on your peers
- Approve a chaincode definition for your organization
- Commit the chaincode definition to the channel

## Package Build Chaincode

~~~
export PATH=$PWD/../../fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$PWD/../../fabric-samples/config/
peer lifecycle chaincode package docstore.tar.gz --path ./ --lang golang --label docstore_1.0

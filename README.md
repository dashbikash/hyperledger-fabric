# Hyperledger Fabric Development

## Start Network and Create Channel for Development
1. curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
2. ./install-fabric.sh --fabric-version 2.5.8 d s b
3. cd fabric-samples/test-network
4. ./network.sh up createChannel -c mychannel -ca
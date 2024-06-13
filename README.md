# Hyperledger Fabric Development

## Start Network and Create Channel for Development
1. curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
2. ./install-fabric.sh --fabric-version 2.5.8 d s b
3. cd fabric-samples/test-network
4. ./network.sh up createChannel -c mychannel -ca
5. export 
   ~~~
   PATH=${PWD}/../bin:$PATH
   
6. Important : Set FABRIC_CFG_PATH
    ~~~
    export FABRIC_CFG_PATH=${PWD}/../config/
7. Run peer as Org1 
    ~~~ 
    export CORE_PEER_TLS_ENABLED=true \
    export CORE_PEER_LOCALMSPID=Org1MSP \
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp \
    export CORE_PEER_ADDRESS=localhost:7051

8. Check Channels 
    ~~~ 
    peer channel list



Troubleshoot :
Error in Alpine: 
apk add --no-cache libaio libnsl libc6-compat gcompat jq


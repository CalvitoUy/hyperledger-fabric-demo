# Hyperledger Fabric Network

The Hyperledger Fabric blockchain network includes the seller and one courier company with possibility to add more courier companies later.

Seller and courier are in different Hyperledger Fabric organizations and have a private communication channel: `shipment`. Each of them has 2 peers and 1 anchor peer.

A solo orderer type is used for simplicity. In production kafka orderer type with multiple orderers must be used.

Both seller and courier have one CA.

Network configuration is based on the Hyperledger Fabric's tutorial ["Build Your First Network"](http://hyperledger-fabric.readthedocs.io/en/release-1.0/build_network.html), Ivan Vankov's videos ["Hyperledger Fabric - build first network"](https://www.youtube.com/playlist?list=PLjsqymUqgpSTGC4L6ULHCB_Mqmy43OcIh) and  [yeasy/docker-compose-files](https://github.com/yeasy/docker-compose-files/tree/master/hyperledger_fabric/v1.0.6) ([docs](https://github.com/yeasy/docker-compose-files/tree/master/hyperledger_fabric/docs)).

Hyperledger Fabric 1.0.6 is used.

## Create the network

1. Configure Orderer and Peer organizations and generate cryptographic materials. Generate orderer genesis block, `shipment` channel configuration file and anchor peers updates.
    - Adjust the file `crypto-config.yaml` with the organization names and domains of the orderer, seller and courier. Setup number of peers and users for the peer organizations.
    - Adjust the file `configtx.yaml` with the organization names and domains of the orderer, seller and courier. Change the profile names of the Orderer genesis and the channel.
    - Generate certificates and orderer, channel and anchor peers artifacts: `./byfn.sh -m generate`
    - In `docker-compose-cli.yaml` change the seller and courier CA containers keyfile, as it's changed when new cryptographic materials are generated.
2. Bring up the network: `./byfn.sh -m up`.

The chaincode_example02 chaincode is instantiated in the network for end-to-end testing purpose.

If you have to start the network from scratch, clear the network, all docker containers, images and volumes: `./byfn.sh -m down` and start from point 1.
    
## Start and stop the network

After the network is created you DON'T need to create it again every time you want to start and stop it. Use these commands:
1. Start already created network: `docker-compose -f docker-compose-cli.yaml up -d`
2. Stop the network without deleting the containers, volumes and networks: `docker-compose -f docker-compose-cli.yaml stop`

## Deploy a chaincode
1. Start the network and enter the CLI container to deploy and test the chaincode: `docker-compose -f docker-compose-cli.yaml run cli bash`
2. Deploy a chaincode:
    - Install `shipment` chaincode on seller/peer0: `CORE_PEER_LOCALMSPID="SellerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/peers/peer0.seller.blockchain.localhost/tls/ca.crt CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/users/Admin@seller.blockchain.localhost/msp CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 peer chaincode install -n shipmentcc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/shipment`
    - Instantiate `shipment` chaincode on seller/peer0: `CORE_PEER_LOCALMSPID="SellerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/peers/peer0.seller.blockchain.localhost/tls/ca.crt CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/users/Admin@seller.blockchain.localhost/msp CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 peer chaincode instantiate -o orderer.blockchain.localhost:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/blockchain.localhost/orderers/orderer.blockchain.localhost/msp/tlscacerts/tlsca.blockchain.localhost-cert.pem -C shipment -n shipmentcc -v 1.0 -c '{"Args":["init"]}' -P "OR ('SellerMSP.member','CourierMSP.member')"`
    - Upgrade `shipment` chaincode on seller/peer0: `...`
    
## Invoke and query chaincode

1. Invoke `shipment` chaincode on seller/peer0: `CORE_PEER_LOCALMSPID="SellerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/peers/peer0.seller.blockchain.localhost/tls/ca.crt CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/users/Admin@seller.blockchain.localhost/msp CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 peer chaincode invoke -o orderer.blockchain.localhost:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/blockchain.localhost/orderers/orderer.blockchain.localhost/msp/tlscacerts/tlsca.blockchain.localhost-cert.pem -C shipment -n shipmentcc -c '{"Args":["initLedger"]}'` 
2. Query `shipment` chaincode on seller\peer0: `CORE_PEER_LOCALMSPID="SellerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/peers/peer0.seller.blockchain.localhost/tls/ca.crt CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.blockchain.localhost/users/Admin@seller.blockchain.localhost/msp CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 CORE_PEER_ADDRESS=peer0.seller.blockchain.localhost:7051 peer chaincode query -C shipment -n shipmentcc -c '{"Args":["queryCar","CAR4"]}'`  
3. See `shipment` chaincode logs on seller\peer0: `dev-peer0.seller.blockchain.localhost-shipmentcc-1.0`

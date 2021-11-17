#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

export CA1KEY=$(cd crypto-config/peerOrganizations/org1.reshop.com/ca/ && ls *_sk)
export CA2KEY=$(cd crypto-config/peerOrganizations/org2.reshop.com/ca/ && ls *_sk)
export CA3KEY=$(cd crypto-config/peerOrganizations/org3.reshop.com/ca/ && ls *_sk)

docker-compose -f docker-compose.yml down

# docker-compose -> 컨테이터수행 및 net_basic 네트워크 생성
docker-compose -f docker-compose.yml up -d ca.org1.reshop.com ca.org2.reshop.com ca.org3.reshop.com orderer.reshop.com peer0.org1.reshop.com  peer0.org2.reshop.com peer0.org3.reshop.com cli

docker ps -a
docker network ls

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=5
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel -> rechannel.block cli working dir 복사
docker exec cli peer channel create -o orderer.reshop.com:7050 -c rechannel -f /etc/hyperledger/configtx/channel.tx
# clie workding dir (/etc/hyperledger/configtx/) rechannel.block

sleep 3

# Join peer0.org1.reshop.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.reshop.com/msp" peer0.org1.reshop.com peer channel join -b /etc/hyperledger/configtx/rechannel.block

sleep 3

# Join peer0.org2.reshop.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.reshop.com/msp" peer0.org2.reshop.com peer channel join -b /etc/hyperledger/configtx/rechannel.block

sleep 3

# Join peer0.org3.reshop.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org3.reshop.com/msp" peer0.org3.reshop.com peer channel join -b /etc/hyperledger/configtx/rechannel.block

sleep 3


# anchor ORG1 rechannel update
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.reshop.com/msp" peer0.org1.reshop.com peer channel update -f /etc/hyperledger/configtx/Org1MSPanchors.tx -c rechannel -o orderer.reshop.com:7050

sleep 3

# anchor ORG2 rechannel update
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.reshop.com/msp" peer0.org2.reshop.com peer channel update -f /etc/hyperledger/configtx/Org2MSPanchors.tx -c rechannel -o orderer.reshop.com:7050

sleep 3

# anchor ORG2 rechannel update
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org3.reshop.com/msp" peer0.org3.reshop.com peer channel update -f /etc/hyperledger/configtx/Org3MSPanchors.tx -c rechannel -o orderer.reshop.com:7050

sleep 3

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.reshop.com/msp" peer0.org1.reshop.com peer channel list
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.reshop.com/msp" peer0.org2.reshop.com peer channel list
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org3.reshop.com/msp" peer0.org3.reshop.com peer channel list



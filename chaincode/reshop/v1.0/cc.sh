#!/bin/bash

set -x 

CCNAME=reshop
VER=1.0.6
CHANNEL=rechannel

# 1. 설치
docker exec cli peer chaincode install -n $CCNAME -v $VER -p github.com/$CCNAME/v1.0
docker exec cli peer chaincode list --installed

docker exec -e "CORE_PEER_ADDRESS=peer0.org2.reshop.com:8051" -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.reshop.com/users/Admin@org2.reshop.com/msp" cli peer chaincode install -n $CCNAME -v $VER -p github.com/$CCNAME/v1.0
docker exec -e "CORE_PEER_ADDRESS=peer0.org2.reshop.com:8051" -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.reshop.com/users/Admin@org2.reshop.com/msp" cli peer chaincode list --installed

# 2. 배포 a = 100
docker exec cli peer chaincode upgrade -n $CCNAME -v $VER -C $CHANNEL -c '{"Args":[]}' -P 'OR ("Org1MSP.member","Org2MSP.member","Org3MSP.member")'
sleep 3
docker exec cli peer chaincode list --instantiated -C $CHANNEL

# register, respond, request, complete, pay
# // 견적요청 파라미터 ContractID, CustomerID, CarInfo
docker exec cli peer chaincode invoke -n $CCNAME -C $CHANNEL -c '{"Args":["register","C0001","bstudent","BMW-320d-2015-....."]}'
sleep 3
# // 견적서등록 파라미터 ContractID, ShopID, ExpectedRepairItems, Price
docker exec cli peer chaincode invoke -n $CCNAME -C $CHANNEL -c '{"Args":["respond","C0001","SHOP101","OIL replace.. FILTER replace...", "100000"]}'
sleep 3
# // 수리요청 파라미터 ContractID, CustomerID (결제정보등록)
docker exec cli peer chaincode invoke -n $CCNAME -C $CHANNEL -c '{"Args":["request","C0001","bstudent"]}'
sleep 3
# // 수리이력등록 파라미터 ContractID, ShopID, RepairRecord
docker exec cli peer chaincode invoke -n $CCNAME -C $CHANNEL -c '{"Args":["complete","C0001","SHOP101","OIL replace.. FILTER replace..."]}'
sleep 3
# // 수리컨펌 결제 파라미터 ContractID, CustomerID
docker exec cli peer chaincode invoke -n $CCNAME -C $CHANNEL -c '{"Args":["pay","C0001","bstudent"]}'
sleep 3

# 3. query history
docker exec cli peer chaincode query -n $CCNAME -C $CHANNEL -c '{"Args":["history","C0001"]}'

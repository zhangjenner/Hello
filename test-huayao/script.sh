#!/bin/bash

echo "#############################################################################"
echo
DELAY=1
TIMEOUT="15"
COUNTER=1
MAX_RETRY=5
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.com/tlsca/tlsca.orderer.com-cert.pem

#==================================================================================================================
#设置全局参数
setGlobals () {
    echo "Connecting to $1:peer$2"
    ORG=$1
    NO=$2
    if [ "$ORG" == "HuaYao" ]; then
        DOMAIN="huayao.com"
    elif [ "$ORG" == "BankA" ]; then
        DOMAIN="bank-a.com"
    elif [ "$ORG" == "BankB" ]; then
        DOMAIN="bank-b.com"
    elif [ "$ORG" == "DrylineA" ]; then
        DOMAIN="dryline-a.com"
    elif [ "$ORG" == "DrylineB" ]; then
        DOMAIN="dryline-b.com"
    elif [ "$ORG" == "CustomerA" ]; then
        DOMAIN="customer-a.com"
    elif [ "$ORG" == "CustomerB" ]; then
        DOMAIN="customer-b.com"
    elif [ "$ORG" == "CustomerC" ]; then
        DOMAIN="customer-c.com"
    fi        
    
    CORE_PEER_LOCALMSPID="${ORG}MSP"
    CORE_PEER_ADDRESS=peer${NO}.${DOMAIN}:7051
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${DOMAIN}/tlsca/tlsca.${DOMAIN}-cert.pem
	CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${DOMAIN}/users/Admin@${DOMAIN}/msp

	#env |grep CORE
}

#验证执行结果
verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute End-2-End Scenario ==========="
		echo
   		exit 1
	fi
}

#==================================================================================================================
#创建通道
createChannel() {
    for chan in chan-pub chan-; do
        sleep $DELAY
		setGlobals $org 0    
        peer channel join -b $CHANNEL_NAME.block  >&log.txt
		res=$?
        cat log.txt
        verifyResult $res "$org has failed to Join the Channel \"$CHANNEL_NAME\""
		echo "===================== $org joined on the channel \"$CHANNEL_NAME\" ===================== "
		echo
	done
    
    sleep $DELAY
	setGlobals HuaYao 0
    local CHANNEL_NAME="channel-ab"
	peer channel create -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/ChannelAB.tx --tls --cafile $ORDERER_CA >&log.txt
	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
    
    sleep $DELAY
    setGlobals HuaYao 1
    local CHANNEL_NAME="channel-bc"
	peer channel create -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/ChannelBC.tx --tls --cafile $ORDERER_CA >&log.txt
	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
}

#==================================================================================================================
#加入通道
joinChannel () {
    local CHANNEL_NAME="channel-ab"
	for org in HuaYao BankA DrylineA CustomerA CustomerB; do
        sleep $DELAY
		setGlobals $org 0    
        peer channel join -b $CHANNEL_NAME.block  >&log.txt
		res=$?
        cat log.txt
        verifyResult $res "$org has failed to Join the Channel \"$CHANNEL_NAME\""
		echo "===================== $org joined on the channel \"$CHANNEL_NAME\" ===================== "
		echo
	done
    
    local CHANNEL_NAME="channel-bc"
	for org in HuaYao BankB DrylineB CustomerB CustomerC; do
        sleep $DELAY
        if [ "$org" == "HuaYao" ]; then
		    setGlobals $org 1
        else 
            setGlobals $org 0
        fi
		peer channel join -b $CHANNEL_NAME.block  >&log.txt
		res=$?
        cat log.txt
        verifyResult $res "$org has failed to Join the Channel \"$CHANNEL_NAME\""
		echo "===================== $org joined on the channel \"$CHANNEL_NAME\" ===================== "
		echo
	done
}

#==================================================================================================================
#更新锚点
updateAnchorPeers() {
    local CHANNEL_NAME="channel-ab"
    for org in HuaYao BankA DrylineA CustomerA CustomerB; do 
        sleep $DELAY
        setGlobals $org 0
        peer channel update -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/AB${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt
        res=$?
        cat log.txt
        verifyResult $res "Anchor peer update failed"
        echo "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
        echo
    done
    
    local CHANNEL_NAME="channel-bc"
    for org in HuaYao BankB DrylineB CustomerB CustomerC; do    
        sleep $DELAY
        if [ "$org" == "HuaYao" ]; then
		    setGlobals $org 1
        else 
            setGlobals $org 0
        fi
        peer channel update -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/BC${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt
        res=$?
        cat log.txt
        verifyResult $res "Anchor peer update failed"
        echo "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
        echo
    done
}

#==================================================================================================================
#安装链码
installChaincode () {
    local CHANNEL_NAME="channel-ab"
    for org in HuaYao BankA DrylineA CustomerA CustomerB; do    
        sleep $DELAY
        setGlobals $org 0
        peer chaincode install -n ccab -v 1.0 -p github.com/hyperledger/fabric/peer/chaincode/example02 >&log.txt
        res=$?
        cat log.txt
        verifyResult $res "Chaincode installation on remote peer $org has Failed"
	    echo "===================== Chaincode is installed on remote $org ===================== "
        echo
    done
    
    local CHANNEL_NAME="channel-bc"
    for org in HuaYao BankB DrylineB CustomerB CustomerC; do    
        sleep $DELAY
        if [ "$org" == "HuaYao" ]; then
		    setGlobals $org 1
        else 
            setGlobals $org 0
        fi
        peer chaincode install -n ccbc -v 1.0 -p github.com/hyperledger/fabric/peer/chaincode/example02 >&log.txt
        res=$?
        cat log.txt
        verifyResult $res "Chaincode installation on remote peer $org has Failed"
        echo "===================== Chaincode is installed on remote $org ===================== "
        echo
    done
}

#==================================================================================================================
#实例化链码
instantiateChaincode () {
    sleep $DELAY
	setGlobals HuaYao 0
    local CHANNEL_NAME="channel-ab"
	peer chaincode instantiate -o ord.orderer.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ccab -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "OR	('HuaYaoMSP.member','BankAMSP.member')" >&log.txt
    res=$?
	cat log.txt
	verifyResult $res "Chaincode instantiation on HuaYao on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Instantiation on channel '$CHANNEL_NAME' is successful ===================== "
	echo
    
    sleep $DELAY
	setGlobals HuaYao 1
    local CHANNEL_NAME="channel-bc"
	peer chaincode instantiate -o ord.orderer.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ccbc -v 1.0 -c '{"Args":["init","c","300","d","400"]}' -P "OR	('HuaYaoMSP.member','BankBMSP.member')" >&log.txt
    res=$?
	cat log.txt
	verifyResult $res "Chaincode instantiation on HuaYao on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Instantiation on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

#==================================================================================================================
#查询连码
chaincodeQuery () {
    if [ "$1" == "HuaYao" -a "$2" == "ccbc" ]; then
        setGlobals $1 1
    else
        setGlobals $1 0
    fi    
    local CHANNEL_NAME="channel-ab"
    if [ "$2" == "ccbc" ]; then
        CHANNEL_NAME="channel-bc"
    fi
    
    local starttime=$(date +%s)
    local querystr="{\"Args\":[\"query\",\"$3\"]}"
    peer chaincode query -C $CHANNEL_NAME -n $2 -c "$querystr" >&log.txt
    res=$?
    echo "Attempting to Query ...$(($(date +%s)-starttime)) secs"
    cat log.txt
    verifyResult $res "Chaincode query on remote peer $1.$2 has Failed!"  
    echo "===================== Chaincode is query on remote $1.$2.$3 is successful ===================== "
    echo    
}

#==================================================================================================================
#查询连码
chaincodeInvoke () {
    if [ "$1" == "HuaYao" -a "$2" == "ccbc" ]; then
        setGlobals $1 1
    else
        setGlobals $1 0
    fi
    local CHANNEL_NAME="channel-ab"
    if [ "$2" == "ccbc" ]; then
        CHANNEL_NAME="channel-bc"
    fi
    
    local invokestr="{\"Args\":[\"invoke\",\"$3\",\"$4\",\"$5\"]}"
	peer chaincode invoke -o ord.orderer.com:7050  --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $2 -c $invokestr >&log.txt
	res=$?
	cat log.txt
	verifyResult $res "Invoke execution on $1.$2.$3.$4 failed "
	echo "===================== Invoke transaction on $1.$2.$3.$4 on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

#==================================================================================================================
## Create channel
if [ "$1" == "-c" ]; then
    echo "Creating channel..."
    createChannel
    echo "Having all peers join the channel..."
    joinChannel
    echo "Updating anchor peers ..."
    updateAnchorPeers

## Set the anchor peers for each org in the channel
elif [ "$1" == "-u" ]; then
    echo "Updating anchor peers ..."
    updateAnchorPeers

## Install chaincode
elif [ "$1" == "-i" ]; then
    echo "Installing chaincode..."
    installChaincode

## Instantiate chaincode
elif [ "$1" == "-s" ]; then
    echo "Instantiating chaincode..."
    instantiateChaincode 

## Query on chaincode
elif [ "$1" == "-q" ]; then
    echo "Querying chaincode ..."
    chaincodeQuery $2 $3 $4

## Invoke on chaincode
elif [ "$1" == "-v" ]; then
    echo "Invoke chaincode ..."
    chaincodeInvoke $2 $3 $4 $5 $6
fi

echo
echo "#############################################################################"

exit 0

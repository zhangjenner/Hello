#!/bin/bash

echo "#############################################################################"
echo
DELAY=1
TIMEOUT="15"
COUNTER=1
MAX_RETRY=5
CMDSTR=""
CHANNEL_NAME="channel-ab"
BASE_PATH=$(cd `dirname $0`; pwd)
ORDERER_CA=${BASE_PATH}/crypto/ordererOrganizations/orderer.com/tlsca/tlsca.orderer.com-cert.pem

#==================================================================================================================
#设置全局参数
setGlobals () {
    echo "Connecting to $1:peer$2"
    ORG=$1
    NO=$2
    if [ "$ORG" == "BankA" ]; then
        DOMAIN="banka.com"
    elif [ "$ORG" == "BankB" ]; then
        DOMAIN="bankb.com"    
    fi        
    
    CORE_PEER_LOCALMSPID="${ORG}MSP"
    CORE_PEER_ADDRESS=peer${NO}.${DOMAIN}:7051
    CORE_PEER_TLS_ROOTCERT_FILE=${BASE_PATH}/crypto/peerOrganizations/${DOMAIN}/tlsca/tlsca.${DOMAIN}-cert.pem
	CORE_PEER_MSPCONFIGPATH=${BASE_PATH}/crypto/peerOrganizations/${DOMAIN}/users/Admin@${DOMAIN}/msp

	#env |grep CORE
}

#调用命令组合
compCmd () {
    args=($@)
    CMDSTR="{\"Args\":["
    for arg in ${args[@]}
    do
        CMDSTR=${CMDSTR}"\"$arg\","
    done
    CMDSTR=${CMDSTR%,*}
    CMDSTR=${CMDSTR}"]}"
    echo "CMD: "$CMDSTR
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
    sleep $DELAY
	setGlobals BankA 0
    if [ "$CORE_PEER_TLS_ENABLED" = "false" ]; then
	    peer channel create -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/ChannelAB.tx >&log.txt
    else
        peer channel create -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/ChannelAB.tx --tls --cafile $ORDERER_CA >&log.txt
    fi
	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
}

#==================================================================================================================
#加入通道
joinChannel () {
	for org in BankA BankB; do
        sleep $DELAY
		setGlobals $org 0    
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
    for org in BankA BankB; do 
        sleep $DELAY
        setGlobals $org 0
        if [ "$CORE_PEER_TLS_ENABLED" = "false" ]; then
            peer channel update -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
        else
            peer channel update -o ord.orderer.com:7050 -c $CHANNEL_NAME -f ./artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt
        fi
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
    sleep $DELAY
    setGlobals $1 0
    peer chaincode install -n $2 -v $3 -p github.com/jenner/chaincode/hycc/src/$4 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode:$2-$3 installation on remote peer $1 has Failed"
    echo "===================== Chaincode is installed on remote $org ===================== "
    echo
}

#==================================================================================================================
#升级链码
upgradeChaincode () {    
    sleep $DELAY
    setGlobals $1 0
    args=($@)
    compCmd "re-init" ${args[@]:3}
    peer chaincode upgrade -C $CHANNEL_NAME -n $2 -v $3 -c $CMDSTR >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "Chaincode:$2-$3 upgrade on remote peer $1 has Failed"
    echo "===================== Chaincode is upgrade on remote $org ===================== "
    echo
}

#==================================================================================================================
#实例化链码
instantiateChaincode () {
    sleep $DELAY
	setGlobals $1 0
    args=($@)
    compCmd "init" ${args[@]:3}
    if [ "$CORE_PEER_TLS_ENABLED" = "false" ]; then
	    peer chaincode instantiate -o ord.orderer.com:7050 -C $CHANNEL_NAME -n $2 -v $3 -c $CMDSTR -P "OR ('BankAMSP.member','BankBMSP.member')" >&log.txt
    else
        peer chaincode instantiate -o ord.orderer.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $2 -v $3 -c $CMDSTR -P "OR ('BankAMSP.member','BankBMSP.member')" >&log.txt
    fi
    res=$?
	cat log.txt
	verifyResult $res "Chaincode:$2-$3 instantiation on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Instantiation on channel '$CHANNEL_NAME' is successful ===================== "
	echo 
    sleep 2
}

#==================================================================================================================
#查询连码
chaincodeInvoke () {   
    sleep $DELAY
    setGlobals $1 0
    args=($@)
    compCmd ${args[@]:2}
    if [ "$CORE_PEER_TLS_ENABLED" = "false" ]; then
	    peer chaincode invoke -o ord.orderer.com:7050 -C $CHANNEL_NAME -n $2 -c $CMDSTR >&log.txt
    else
        peer chaincode invoke -o ord.orderer.com:7050  --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $2 -c $CMDSTR >&log.txt
    fi
	res=$?
	cat log.txt
	verifyResult $res "Invoke execution chaincode:$2 on channel '$CHANNEL_NAME' failed "
	echo "===================== Invoke transaction chaincode:$2 on channel '$CHANNEL_NAME' is successful ===================== "
	echo
}

#==================================================================================================================
## Create channel
if [ "$1" == "-c" ]; then
    echo "Creating channel..."
    createChannel
    
## Create channel
elif [ "$1" == "-j" ]; then
    echo "Having all peers join the channel..."
    joinChannel
    echo "Updating anchor peers ..."
    updateAnchorPeers

## Install chaincode
elif [ "$1" == "-i" ]; then
    echo "Installing chaincode..."
    if [ $# -lt 4 ]; then
        echo "ERROR: Expected parameters"
        echo "-i [org] [name] [ver] [path]"
        exit 1
    fi
    args=($@)
    installChaincode ${args[@]:1}
    
## Upgrade chaincode
elif [ "$1" == "-u" ]; then
    echo "Upgrade chaincode..."
    if [ $# -lt 4 ]; then
        echo "ERROR: Expected parameters"
        echo "-u [org] [name] [ver] [arg...]"
        exit 1
    fi
    args=($@)
    upgradeChaincode ${args[@]:1}

## Instantiate chaincode
elif [ "$1" == "-s" ]; then
    echo "Instantiating chaincode..."
    if [ $# -lt 4 ]; then
        echo "ERROR: Expected parameters"
        echo "-s [org] [name] [ver] [arg...]"
        exit 1
    fi
    args=($@)
    instantiateChaincode ${args[@]:1}

## Invoke on chaincode
elif [ "$1" == "-v" ]; then    
    echo "Invoke chaincode ..."
    if [ $# -lt 3 ]; then
        echo "ERROR: Expected parameters"
        echo "-v [org] [name] [arg...]"
        exit 1
    fi
    args=($@)
    chaincodeInvoke ${args[@]:1}
fi

echo
echo "#############################################################################"

exit 0

#!/bin/bash

echo "#############################################################################"
echo

DELAY=1
FILEN=""

#==================================================================================================================
splitFilename () {
    arg=$*
    OLD_IFS="$IFS" 
    IFS="."
    FILEN=($arg)
    IFS="$OLD_IFS"   
}

#==================================================================================================================
#启动Configtxlator
startConfigtxlator () {
    configtxlator start --hostname="0.0.0.0" --port 7059 &
}

#停止Configtxlator
stopConfigtxlator () {
    ps -ef | grep configtxlator | grep -v grep | awk '{print $2}' | xargs kill -9
}

#==================================================================================================================
#编解码区块/交易
codecBT () {    
    splitFilename $2
    if [ ${#FILEN[@]} -ne 2 ]; then
        echo "ERROR: File name format is error!"
        exit 1
    elif [ ${FILEN[1]} != "block" -a ${FILEN[1]} != "json" ]; then
        echo "ERROR: Unknown file type!"
        exit 1
    fi    
    
    echo "$2 -> $3"
    if [ "$1" == "-b" ]; then
        if [ "${FILEN[1]}" == "block" ]; then            
            curl -X POST --data-binary @$2 http://127.0.0.1:7059/protolator/decode/common.Block > $3
        elif [ "${FILEN[1]}" == "json" ]; then
            curl -X POST --data-binary @$2 http://127.0.0.1:7059/protolator/encode/common.Block > $3
        fi
    elif [ "$1" == "-t" ]; then
        if [ "$FILEN[1]" == "block" ]; then
            curl -X POST --data-binary @$2 http://127.0.0.1:7059/protolator/decode/common.Envelope > $3
        elif [ "$FILEN[1]" == "json" ]; then
            curl -X POST --data-binary @$2 http://127.0.0.1:7059/protolator/encode/common.Envelope > $3
        fi
    fi

}

#==================================================================================================================
#启动Configtxlator
if [ "$1" == "-s" ]; then
    echo "Start configtxlator..."
    startConfigtxlator
    echo 

#编解码区块
elif [ "$1" == "-b" -o "$1" == "-t" ]; then
    echo "CODEC block/transaction..."
    if [ $# -lt 3 ]; then
        echo "ERROR: Expected parameters"
        echo "-b [form] [to]"
        exit 1
    fi
    args=$@
    codecBT $*

#启动Configtxlator
elif [ "$1" == "-d" ]; then
    echo "Stop configtxlator..."
    stopConfigtxlator
fi
    
echo
echo "#############################################################################"
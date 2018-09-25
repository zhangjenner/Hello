#!/bin/bash

export FABRIC_CFG_PATH=${PWD}

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  run.sh -m up|down|gen [-t <timeout>] [-d <delay>] [-f <docker-compose-file>]"
  echo "  run.sh -h|--help (print this message)"
  echo "    -m <mode> - one of 'up', 'down', 'restart' or 'generate'"
  echo "      - 'up' - bring up the network with docker-compose up"
  echo "      - 'down' - clear the network with docker-compose down"
  echo "      - 'gen' - generate required certificates and genesis block"
  echo "    -t <timeout> - CLI timeout duration in microseconds (defaults to 10000)"
  echo "    -d <delay> - delay duration in seconds (defaults to 3)"
  echo "    -f <docker-compose-file> - specify which docker-compose file use (defaults to docker-compose-cli.yaml)"
  echo
  echo "Taking all defaults:"
  echo "	run.sh -m gen"
  echo "	run.sh -m up"
  echo "	run.sh -m down"
}

# =============================================================================
# Generate the needed certificates, the genesis block and start the network.
function networkUp () {
  # generate artifacts if they don't exist
  if [ ! -d "crypto-config" ]; then
    generateCerts
    replacePrivateKey
    generateChannelArtifacts
  fi
  TIMEOUT=$CLI_TIMEOUT DELAY=$CLI_DELAY docker-compose -f $COMPOSE_FILE up -d 2>&1
  if [ $? -ne 0 ]; then
    echo "ERROR !!!! Unable to start network"
    docker logs -f cli
    exit 1
  fi
  #docker logs -f cli
}

# -----------------------------------------------------------------------------
# Tear down running network
function networkDown () {
  docker-compose -f $COMPOSE_FILE down
  # Don't remove containers, images, etc if restarting
  if [ "$MODE" != "restart" ]; then
    #Cleanup the chaincode containers
    clearContainers
    #Cleanup images
    removeUnwantedImages
  fi
}

# Obtain CONTAINER_IDS and remove them
function clearContainers () {
  CONTAINER_IDS=$(docker ps -aq)
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# Delete any images that were generated as a part of this setup
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

# =============================================================================
# Generates Org certs using cryptogen tool
function generateCerts (){
  which cryptogen
  if [ "$?" -ne 0 ]; then
    echo "cryptogen tool not found. exiting"
    exit 1
  fi
  echo
  echo "##########################################################"
  echo "##### Generate certificates using cryptogen tool #########"
  echo "##########################################################"
  if [ -d "crypto-config" ]; then
    rm -Rf crypto-config/*
  else
    mkdir crypto-config
  fi
  cryptogen generate --config=./crypto-config.yaml
  if [ "$?" -ne 0 ]; then
    echo "Failed to generate certificates..."
    exit 1
  fi
  echo
}

# -----------------------------------------------------------------------------
# Generate orderer genesis block, channel configuration transaction 
function generateChannelArtifacts() {
  which configtxgen
  if [ "$?" -ne 0 ]; then
    echo "configtxgen tool not found. exiting"
    exit 1
  fi
  
  if [ -d "channel-artifacts" ]; then
    rm -Rf channel-artifacts/*
  else
    mkdir channel-artifacts
  fi

  echo "##########################################################"
  echo "#########  Generating Orderer Genesis block ##############"
  echo "##########################################################"
  # Note: For some unknown reason (at least for now) the block file can't be
  # named orderer.genesis.block or the orderer will fail to launch!
  configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block
  if [ "$?" -ne 0 ]; then
    echo "Failed to generate orderer genesis block..."
    exit 1
  fi
  
  echo
  echo "#################################################################"
  echo "### Generating channel configuration transaction 'ChannelUA.tx'##"
  echo "#################################################################"
  configtxgen -profile ChannelUA -outputCreateChannelTx ./channel-artifacts/ChannelUA.tx -channelID channel-ua
  if [ "$?" -ne 0 ]; then
    echo "Failed to generate channel configuration transaction..."
    exit 1
  fi
  
  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for ChannelUA   ########"
  echo "#################################################################"
  for MSP in AdminMSP UserMSP; do
    configtxgen -profile ChannelUA -outputAnchorPeersUpdate ./channel-artifacts/${MSP}anchors.tx -channelID channel-ua -asOrg ${MSP}
    if [ "$?" -ne 0 ]; then
      echo "Failed to generate anchor peer update ..."
      exit 1
    fi
  done

}

# ###################################################################################
# Obtain the OS and Architecture string that will be used to select the correct
# native binaries for your platform
OS_ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
# timeout duration - the duration the CLI should wait for a response from
# another container before giving up
CLI_TIMEOUT=10000
#default for delay
CLI_DELAY=3
# use this as the default docker-compose yaml definition
COMPOSE_FILE=docker-compose.yaml

# Parse commandline args
while getopts "h?m:c:t:d:f:s:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    m)  MODE=$OPTARG
    ;;
    t)  CLI_TIMEOUT=$OPTARG
    ;;
    d)  CLI_DELAY=$OPTARG
    ;;
    f)  COMPOSE_FILE=$OPTARG
    ;;
  esac
done

#Create the network using docker compose
if [ "${MODE}" == "up" ]; then
  networkUp
  elif [ "${MODE}" == "down" ]; then 
  networkDown
  elif [ "${MODE}" == "gen" ]; then 
  generateCerts
  generateChannelArtifacts
else
  printHelp
  exit 1
fi

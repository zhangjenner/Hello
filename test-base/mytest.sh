#!/bin/bash


BASE_PATH=$(cd `dirname $0`; pwd)

#==================================================================================================================
#安装公用链码及数据
installPubcc () {
    ./script.sh -i BankA pubcc 1.0 ./code/pubcc
    ./script.sh -s BankA pubcc 1.0 ""
    
    ./script.sh -v BankA pubcc OrgMg add '{\"id\":\"1\",\"mspid\":\"msp1\",\"type\":2,\"name\":\"org1\"}'
    ./script.sh -v BankA pubcc OrgMg add '{\"id\":\"2\",\"mspid\":\"msp2\",\"type\":2,\"name\":\"org2\"}'
    ./script.sh -v BankA pubcc OrgMg add '{\"id\":\"3\",\"mspid\":\"msp3\",\"type\":2,\"name\":\"org3\"}'
    
    ./script.sh -v BankA pubcc PermMg add '{\"id\":\"1\",\"pid\":\"0\",\"perm\":\"perm1\"}'
    ./script.sh -v BankA pubcc PermMg add '{\"id\":\"2\",\"pid\":\"1\",\"perm\":\"perm2\"}'
    ./script.sh -v BankA pubcc PermMg add '{\"id\":\"3\",\"pid\":\"2\",\"perm\":\"perm3\"}'
    
    ./script.sh -v BankA pubcc AuthMg add '{\"opt\":\"read\",\"pexp\":\"perm1\"}'
    ./script.sh -v BankA pubcc AuthMg add '{\"opt\":\"write\",\"pexp\":\"perm2\"}'
    ./script.sh -v BankA pubcc AuthMg add '{\"opt\":\"delete\",\"pexp\":\"perm1&&perm2\"}'

    ./script.sh -v BankA pubcc RoleMg add '{\"id\":\"1\",\"oid\":\"1\",\"pid\":[\"1\"],\"name\":\"role1\"}'
    ./script.sh -v BankA pubcc RoleMg add '{\"id\":\"2\",\"oid\":\"2\",\"pid\":[\"1\",\"2\"],\"name\":\"role2\"}'
    ./script.sh -v BankA pubcc RoleMg add '{\"id\":\"3\",\"oid\":\"3\",\"pid\":[\"1\",\"2\",\"3\"],\"name\":\"role3\"}'

    ./script.sh -v BankA pubcc UserMg add '{\"id\":\"1\",\"oid\":\"1\",\"rid\":[\"1\"],\"name\":\"user1\",\"mobile\":\"1234567890\",\"pubkey\":\"@#$%^&*()et\"}'
    ./script.sh -v BankA pubcc UserMg add '{\"id\":\"2\",\"oid\":\"2\",\"rid\":[\"2\"],\"name\":\"user2\",\"mobile\":\"abcdefghh\",\"pubkey\":\"JOFJS(EI\"}'
    ./script.sh -v BankA pubcc UserMg add '{\"id\":\"3\",\"oid\":\"3\",\"rid\":[\"3\"],\"name\":\"user3\",\"mobile\":\"hddey44434\",\"pubkey\":\":fke)o(o(\"}'
}

#安装订单链码及数据
installOrdercc () {
    ./script.sh -i BankA ordercc 1.0 ./code/ordercc
    ./script.sh -s BankA ordercc 1.0 ""
}

#==================================================================================================================
echo "#############################################################################"
## 安装公用链码及数据
if [ "$1" == "-p" ]; then
    echo "Installing pubcc..."
    installPubcc
    
## 安装订单链码及数据
elif [ "$1" == "-o" ]; then
    echo "Installing ordercc..."
    installOrdercc    
fi

echo
echo "#############################################################################"

exit 0

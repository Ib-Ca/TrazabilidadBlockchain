package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct{ contractapi.Contract }

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create chaincode: %s\n", err.Error())
		return
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err.Error())
	}
}

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) updateShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start updateShipment")

	// todo: must be update of different parts of the shipment, in different methods
	logger.Debug("TODO")

	logger.Debug("End updateShipment")

	return shim.Success(nil)
}

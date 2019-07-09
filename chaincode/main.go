package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AuctionServiceChaincode struct {
}

func (t *AuctionServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}
func (t *AuctionServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}
	if args[0] == "query" {
		return t.query(stub, args)
	}
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	return shim.Error("Unknown action, check the first argument")
}

func (t *AuctionServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### AuctionServiceChaincode query ###########")
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	state, err := stub.GetState(args[1])

	if err != nil {
		return shim.Error(fmt.Sprintf("query error:%v", err))
	}

	return shim.Success(state)
}

func (t *AuctionServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### AuctionServiceChaincode invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	if len(args) == 3 {
		err := stub.PutState(args[1], []byte(args[2]))
		if err != nil {
			return shim.Error(fmt.Sprintf("invoke error:%v", err))
		}

		// 发送事件
		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(fmt.Sprintf("发送事件异常:%v", err.Error()))
		}

		return shim.Success(nil)
	}
	return shim.Error("Unknown invoke action, check the second argument.")
}

func main() {
	err := shim.Start(new(AuctionServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}

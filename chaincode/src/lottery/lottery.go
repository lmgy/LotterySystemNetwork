package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Lottery struct {
  ObjectType string        `json:"docType"` //field for couchdb
  LotteryType string       `json:"lotterytype"`
  LotteryId string         `json:"lotteryid"`
  AwardPool string         `json:"awardpool"`
  AwardNum  string         `json:"awardnum"`
	IsLottery bool         `json:"islottery"`
	Valid     bool           `json:"Valid"`
	LotteryDate  string      `json:"lotterydate"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

//init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  fmt.Println("cc_Lattery_V1.0 Is Starting Up")
}

//invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

   if function == "read" {             //generic read ledger
  		return read(stub, args)
  	}  else if function == "getHistory"{        //read history of a lottery (audit)
  		return getHistory(stub, args)
  	} else if function == "creat_lottery" {    //deletes a lottery from state
  		return creat_lottery(stub, args)
  	} else if function == "settlement" {    //deletes a lottery from state
  		return settlement(stub, args)
  	}

    // error out
  	fmt.Println("Received unknown invoke function name - " + function)
  	return shim.Error("Received unknown invoke function name - '" + function + "'")
  }

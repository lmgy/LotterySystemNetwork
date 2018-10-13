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

// ============================================================================================================================
// Asset Definitions - The ledger will store tickets and agency
// ============================================================================================================================

/*Ticket*/
type Ticket struct {
	ObjectType string        `json:"docType"`     //field for couchdb
  LotteryType string        `json:"lotterytype"`
  TicketId   string        `json:"ticketid"`        // id of the ticket
  LotteryId  string        `json:"lotteryid"`       // issue id of the lottery
  TicketNum  string        `json:"ticketnum"`       // the lottery number of the ticket ,such as : 12 26 00 05 04 06 14 36
  Agency     string        `json:"agency"`          // which agency
  OwnerID      string        `json:"ownerid"`
  Date       string        `json:"date"`            //date of creat the ticket
  Count       int           `json:"count"`        //how many tickets
  IsLottery  bool          `json:"islottery"`
	IsExchange  bool         `json:"isexchange"`
  AwardLevel int           `json:"awardlevel"`
  Award      int           `json:"award"`
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
  fmt.Println("cc_Ticket_V1.0 Is Starting Up")
}

//invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)
//  different functions
  if function == "init" {                    //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "read" {             //generic read ledger
		return read(stub, args)
	} else if function == "creat_ticket" {            //generic writes to ledger
		return creat_ticket(stub, args)
	} else if function == "getHistory"{        //read history of a ticket (audit)
		return getHistory(stub, args)
	} else if function == "delete_tickets" {    //deletes a ticket from state
		return delete(stub, args)
	} else if function == "exchange" {    //exchange ticket
		return ticket_exchange(stub, args)

  // error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}

func read_everything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		Agencys  []Agency   `json:"agency"`
		Tickets  []Ticket  `json:"tickets"`
	}
	var everything Everything

	// ---- Get All Tickets ---- //
	resultsIterator, err := stub.GetStateByRange("m0", "m9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on ticket id - ", queryKeyAsStr)
		var ticket Ticket
		json.Unmarshal(queryValAsBytes, &ticket)                  //un stringify it aka JSON.parse()
		everything.Tickets = append(everything.Tickets, ticket)   //add this ticket to the list
	}
	fmt.Println("ticket array - ", everything.Tickets)

	// ---- Get All Agencys ---- //
	agencysIterator, err := stub.GetStateByRange("o0", "o9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer agencysIterator.Close()

	for agencysIterator.HasNext() {
		aKeyValue, err := agencysIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on agency id - ", queryKeyAsStr)
		var agency Agency
		json.Unmarshal(queryValAsBytes, &agency)                   //un stringify it aka JSON.parse()

		if agency.Enabled {                                        //only return enabled agencys
			everything.Agencys = append(everything.Agencys, agency)  //add this ticket to the list
		}
	}
	fmt.Println("agency array - ", everything.Agencys)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

//delete
func delete(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_ticket")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	ticketid := args[0]
	owner := args[1]

	// get the ticket
	ticket, err := get_ticket(stub, id)
	if err != nil{
		fmt.Println("Failed to find ticket by ticketid " + id)
		return shim.Error(err.Error())
	}

	// check authorizing owner (see note in set_owner() about how this is quirky)
	if ticket.OwnerID != owner {
		return shim.Error("The owner '" + owner + "' cannot authorize deletion for '" + ticket.Owner + "'.")
	}

	// remove the ticket
	err = stub.DelState(id)                                                 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete ticket")
	}

	fmt.Println("- end delete_ticket")
	return shim.Success(nil)
}

// ============================================================================================================================
// Get ticket - get a ticket asset from ledger
// ============================================================================================================================
func get_ticket(stub shim.ChaincodeStubInterface, id string) (ticket, error) {
	var ticket Ticket
	ticketAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return ticket, errors.New("Failed to find ticket - " + id)
	}
	json.Unmarshal(ticketAsBytes, &ticket)                   //un stringify it aka JSON.parse()

	if ticket.Id != id {                                     //test if ticket is actually here or just nil
		return ticket, errors.New("Ticket does not exist - " + id)
	}

	return ticket, nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error{
	for i, val:= range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}


// ========================================================
// creat Ticket : creat a ticket when the costumer buy a ticket
// ========================================================
func creat_ticket(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting creat_ticket")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation3
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	lottery_type := strings.ToLower(args[0])
	ticket_id := args[1]
	lottery_id := args[2]
	ticket_num := args[3]
	agency := arg[4]
	ownerid := strings.ToLower(args[5])
	date := strings.ToLower(args[6])
	count := args[7]
	//islottery := 0
	//isexchange := 0
	//awardlevel := 0
	//award := 0

	if err != nil {
		return shim.Error("8rd argument must be a numeric string")
	}

	//check if ticket id already exists
	ticket, err := get_ticket(stub, id)
	if err == nil {
		fmt.Println("This ticket already exists - " + id)
		fmt.Println(ticket)
		return shim.Error("This ticket already exists - " + id)  //all stop a ticket by this id exists
	}

	//build the ticket json string manually
	str := `{
		"docType":"ticket",
		"lotterytype": "` + ticket_type +` "
		"ticketid": "` + ticket_id +` "
		"lotteryid": "` + lottery_id +` "
		"ticketnum": "` + ticket_num +` "
		"agency": "` + agency + `"
		"ownerid": "` + ownerid +` "
		"date": "` + date +` "
		"count": "` + count +` "
		"lslottery": "` + 0 +` "
		"lsexchange": "` + 0 +` "
		"awardlevel": "` + unknow +` "
		"award": "` + unknow +` "
	}`
	err = stub.PutState(id, []byte(str))                         //store ticket with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end creat_ticket")
	return shim.Success(nil)
}

func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   Ticket   `json:"value"`
	}
	var history []AuditHistory;
	var ticket Ticket

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketid := args[0]
	fmt.Printf("- start getHistoryForTicket: %s\n", ticketid)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(ticketid)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		json.Unmarshal(historyData.Value, &ticket)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //ticket has been deleted
			var emptyTicket Ticket
			tx.Value = emptyTicket                 //copy nil ticket
		} else {
			json.Unmarshal(historyData.Value, &ticket) //un stringify it aka JSON.parse()
			tx.Value = ticket                      //copy ticket over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForTicket returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ========================================================
// exchange Ticket : exchange the ticket
// ========================================================
func ticket_exchange(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting exchange")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var ticketid = args[0]
	var ownerid = args[1]
	var account = args[2]

	fmt.Println(ticketid + "/" + ownerid + "/" + account )

	lotteryid, err := get_ticket(stub, ticketid)
	if err != nil {
		return shim.Error("This ticket does not exist - " + ticketid)
	}

	// get ticket's current state
	ticketAsBytes, err := stub.GetState(ticketid)
	if err != nil {
		return shim.Error("Failed to get ticket")
	}
	res := Ticket{}
	json.Unmarshal(TicketAsBytes, &res)           //un stringify it aka JSON.parse()

	if res.ownerid = hash(ownerid){
		err := exchange_award()                        // unfinish need bank
		if err != nil {
			return shim.Error("This ticket can't exchange - " + ticketid)
		}else {
			res.isExchange = 1
		}
	}

	fmt.Println("- end exchange")
	return shim.Success(nil)
}

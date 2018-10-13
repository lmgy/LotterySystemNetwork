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

func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   Lottery   `json:"value"`
	}
	var history []AuditHistory;
	var lottery Lottery

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lotteryId := args[0]
	fmt.Printf("- start getHistoryForLottery: %s\n", lotteryId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(lotteryId)
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
		json.Unmarshal(historyData.Value, &lottery)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //lottery has been deleted
			var emptyLottery Lottery
			tx.Value = emptyLottery                 //copy nil lottery
		} else {
			json.Unmarshal(historyData.Value, &lottery) //un stringify it aka JSON.parse()
			tx.Value = lottery                      //copy lottery over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForLottery returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

func creat_lottery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting creat_lottery")

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation3
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	lottery_type := strings.ToLower(args[0])
	lottery_id := args[1]
	awardpool := args[2]
  lotterydate := args[3]
	if err != nil {
		return shim.Error("1rd argument must be a numeric string")
	}

	//check if lottery id already exists
	lottery, err := get_lottery(stub, id)
	if err == nil {
		fmt.Println("This lottery already exists - " + id)
		fmt.Println(lottery)
		return shim.Error("This lottery already exists - " + id)  //all stop a lottery by this id exists
	}

	//build the lottery json string manually
	str := `{
		"docType":"lottery",
		"lotterytype": "` + lottery_type +` "
		"lotteryid": "` + lottery_id +` "
		"awardpool": "` + awardpool +` "
		"awardnum": "` + unknown +` "
		"islottery": "` + 0 +` "
		"valid": "` + 1 +` "
	}`
	err = stub.PutState(lotteryid, []byte(str))                         //store lottery with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end creat_lottery")
	return shim.Success(nil)
}

// ============================================================================================================================
// Get lottery - get a lottery asset from ledger
// ============================================================================================================================
func get_lottery(stub shim.ChaincodeStubInterface, id string) (lottery, error) {
	var lottery Lottery
	lotteryAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return lottery, errors.New("Failed to find lottery - " + id)
	}
	json.Unmarshal(lotteryAsBytes, &lottery)                   //un stringify it aka JSON.parse()

	if lottery.Id != id {                                     //test if lottery is actually here or just nil
		return lottery, errors.New("Lottery does not exist - " + id)
	}

	return lottery, nil
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
// settlement - input awardnum
// ========================================================
func settlement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting settlement")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var lotteryid = args[0]
	var awardnum = args[1]

	fmt.Println(lotteryid + "->" + awardnum )

	lotteryid, err := get_lottery(stub, new_owner_id)
	if err != nil {
		return shim.Error("This lottery does not exist - " + new_owner_id)
	}

	// get lottery's current state
	lotteryAsBytes, err := stub.GetState(lotteryid)
	if err != nil {
		return shim.Error("Failed to get lottery")
	}
	res := Lottery{}
	json.Unmarshal(lotteryAsBytes, &res)           //un stringify it aka JSON.parse()

	// transfer the awardnum
	res.AwardNum = awardnum                   //change the awardnum
	res.IsLottery = 1
	err = stub.PutState(args[0], jsonAsBytes)     //rewrite the lottery with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end settlement")
	return shim.Success(nil)
}

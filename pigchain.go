package main

// 引入相關套件
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// 定義所需的資料格式
// Define the Smart Contract structure
type SmartContract struct {
}

/*
Define Pig structure, with 4 properties.
Structure tags are used by encoding/json library
*/
type Pig struct {
	PigId      string `json:"pigid"`
	Timestamp  string `json:"timestamp"`
	Company    string `json:"company"`
	actionName string `json:"actionname"`
}

// 設定必須的 init 與 invoke 方法
/*
 * The Init method *
 called when the Smart Contract "pig-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "pig-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryPig" {
		return s.queryPig(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordPig" {
		return s.recordPig(APIstub, args)
	} else if function == "queryPigHistory" {
		return s.queryPigHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 生成模擬資料
/*
 * The initLedger method *
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	return shim.Success(nil)
}

// 查詢單隻pig
/*
 * The queryPig method *
 */
func (s *SmartContract) queryPig(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	pigAsBytes, _ := APIstub.GetState(args[0])
	if pigAsBytes == nil {
		return shim.Error("Could not locate pig")
	}
	return shim.Success(pigAsBytes)
}

func (s *SmartContract) queryPigHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	pigAsBytes, err := APIstub.GetHistoryForKey(args[0])
	if err != nil {
		fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
		return shim.Error("Failed to get asset: %s with error: %s", args[0], err)
	}

	if pigAsBytes == nil {
		fmt.Errorf("Could not locate pig")
		return shim.Error("Could not locate pig")
	}
	defer pigAsBytes.Close()

	var value string
	for pigAsBytes.HasNext() {
		result, err2 := pigAsBytes.Next()
		if err2 != nil {
			fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
			return shim.Error("Failed to get asset: %s with error: %s", args[0], err)
		}
		value += string(result.Value) + "||"
	}

	return shim.Success(value)
}

// 增加新紀錄

/*
 * The recordPig method *
 */
func (s *SmartContract) recordPig(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var pig = Pig{PigId: args[0], Company: args[1], actionName: args[2], Timestamp: args[3]}

	pigAsBytes, _ := json.Marshal(pig)
	err := APIstub.PutState(args[0], pigAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record pig catch: %s", args[0]))
	}

	return shim.Success(nil)
}

// 主程序
/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

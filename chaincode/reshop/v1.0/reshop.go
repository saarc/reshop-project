// 패키지 정의
package main

// 외부모듈 포함
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 체인코드 객체정의
type SmartContract struct {
}

// 구조체정의 Repair contract , invoice
type REPAIRCONTRACT struct {
	ContractID   string  `json:"contractid"`
	CustomerID   string  `json:"customerid"`
	CarInfo      string  `json:"carinfo"`
	Invoice      INVOICE `json:"invoice"`
	RepairReport string  `json:"repairreport"`
	Status       string  `json:"status"`
}
type INVOICE struct {
	ShopID      string `json:"shopid"`
	RepairItems string `json:"repairitems"`
	Price       string `json:"price"`
}

// Init 함수 초기 필요 파라미터
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// WS: ContractID, CustomerID, Invoice, RepairReport, Status
// Invoice ( ShopID, ExpectedRepairItems, Price )
// 견적요청 -> 견적서등록 -> 수리요청 -> 수리이력등록 -> 수리컨펌(결제)

// Invoke 함수 기능함수이름들
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	// register, respond, request, complete, pay
	if fn == "register" {
		return s.register(stub, args)
	} else if fn == "respond" {
		return s.respond(stub, args)
	} else if fn == "request" {
		return s.request(stub, args)
	} else if fn == "complete" {
		return s.complete(stub, args)
	} else if fn == "pay" {
		return s.pay(stub, args)
	} else if fn == "history" {
		return s.history(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name")
}

// status : requested-0, invoice-registered-1, repair-confirmed-2, repair-complete-3, paid-4
// 견적요청 파라미터 ContractID, CustomerID, CarInfo
func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. (ContractID, CustomerID, CarInfo)")
	}
	var contract = REPAIRCONTRACT{}
	contract.ContractID = args[0]
	contract.CustomerID = args[1]
	contract.CarInfo = args[2]
	contract.Status = "0"
	contAsBytes, _ := json.Marshal(contract)
	stub.PutState(args[0], contAsBytes)

	return shim.Success([]byte(args[0]))
}

// 견적서등록 파라미터 ContractID, ShopID, ExpectedRepairItems, Price
func (s *SmartContract) respond(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. (ContractID, ShopID, ExpectedRepairItems, Price)")
	}
	contAsBytes, _ := stub.GetState(args[0])
	if contAsBytes == nil {
		return shim.Error("ContractID is not vaild")
	}
	var contract = REPAIRCONTRACT{}
	json.Unmarshal(contAsBytes, &contract)

	invoice := INVOICE{args[1], args[2], args[3]}
	contract.Invoice = invoice
	contract.Status = "1"

	contAsBytes, _ = json.Marshal(contract)

	stub.PutState(args[0], contAsBytes)

	return shim.Success([]byte(args[0]))
}

// 수리요청 파라미터 ContractID, CustomerID (결제정보등록)
func (s *SmartContract) request(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. (ContractID, CustomerID)")
	}
	contAsBytes, _ := stub.GetState(args[0])
	if contAsBytes == nil {
		return shim.Error("ContractID is not vaild")
	}
	var contract = REPAIRCONTRACT{}
	json.Unmarshal(contAsBytes, &contract)

	if contract.CustomerID != args[1] {
		return shim.Error("Wrong customer in the contract")
	}

	contract.Status = "2"

	contAsBytes, _ = json.Marshal(contract)

	stub.PutState(args[0], contAsBytes)

	return shim.Success([]byte(args[0]))
}

// 수리이력등록 파라미터 ContractID, ShopID, RepairRecord
func (s *SmartContract) complete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. (ContractID, ShopID, RepairRecord)")
	}
	contAsBytes, _ := stub.GetState(args[0])
	if contAsBytes == nil {
		return shim.Error("ContractID is not vaild")
	}
	var contract = REPAIRCONTRACT{}
	json.Unmarshal(contAsBytes, &contract)

	if contract.Invoice.ShopID != args[1] {
		return shim.Error("Wrong ShopID in the invoice")
	}
	contract.RepairReport = args[2]
	contract.Status = "3"

	contAsBytes, _ = json.Marshal(contract)

	stub.PutState(args[0], contAsBytes)

	return shim.Success([]byte(args[0]))
}

// 수리컨펌 결제 파라미터 ContractID, Customer
func (s *SmartContract) pay(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. (ContractID, Customer)")
	}
	contAsBytes, _ := stub.GetState(args[0])
	if contAsBytes == nil {
		return shim.Error("ContractID is not vaild")
	}
	var contract = REPAIRCONTRACT{}
	json.Unmarshal(contAsBytes, &contract)

	if contract.CustomerID != args[1] {
		return shim.Error("Wrong Customer in the contract")
	}
	contract.Status = "4"

	contAsBytes, _ = json.Marshal(contract)

	stub.PutState(args[0], contAsBytes)

	return shim.Success([]byte(args[0]))
}

func (s *SmartContract) history(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	keyName := args[0]
	// 로그 남기기
	fmt.Println("readTxHistory:" + keyName)

	resultsIterator, err := stub.GetHistoryForKey(keyName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	// 로그 남기기
	fmt.Println("readTxHistory returning:\n" + buffer.String() + "\n")

	return shim.Success(buffer.Bytes())
}

// 메인
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SFDSmartContract is the definition of the chaincode structure.
type SFDSmartContract struct {
}

type Account struct {
	ID     string  `json:"id"`
	Number string  `json:"number"`
	CPF    string  `json:"cpf"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

// Init is called when the SFDSmartContract is instantiated by the blockchain network.
func (cc *SFDSmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)

	accounts := []Account{
		Account{ID: "1111", Number: "1111", CPF: "303.424.538-66", Amount: 230.10, Status: "ATIVO"},
		Account{ID: "2222", Number: "2222", CPF: "202.303.505-11", Amount: 453.20, Status: "ATIVO"},
		Account{ID: "3333", Number: "3333", CPF: "111.111.111-11", Amount: 1203.34, Status: "ATIVO"},
		Account{ID: "4444", Number: "4444", CPF: "222.333.444-55", Amount: 120.30, Status: "ATIVO"},
	}

	i := 0
	for i < len(accounts) {
		fmt.Println("i is ", i)
		accountAsBytes, _ := json.Marshal(accounts[i])
		stub.PutState(accounts[i].ID, accountAsBytes)
		fmt.Println("Added", accounts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (cc *SFDSmartContract) createAccount(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	id := args[0]

	var account = Account{ID: id, Number: args[0], CPF: args[1], Amount: 0, Status: "ATIVA"}
	accountBytes, _ := json.Marshal(account)
	stub.PutState(account.ID, accountBytes)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *SFDSmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryAccount" {
		return cc.queryAccount(stub, args)
	} else if function == "Init" {
		return cc.Init(stub)
	} else if function == "createAccount" {
		return cc.createAccount(stub, args)
	} else if function == "creditAccount" {
		return cc.creditAccount(stub, args)
	} else if function == "debitAccount" {
		return cc.debitAccount(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (cc *SFDSmartContract) queryAccount(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting UserID")
	}
	accountBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountBytes)
}

func (cc *SFDSmartContract) creditAccount(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting UserID")
	}
	id := args[0]
	amount, err := strconv.ParseFloat(args[1], 32)

	if err != nil {
		return shim.Error(err.Error())
	}

	accountBytes, err := stub.GetState(id)
	account := Account{}
	json.Unmarshal(accountBytes, &account)
	account.Amount = account.Amount + amount

	accountBytes, _ = json.Marshal(account)
	stub.PutState(account.ID, accountBytes)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountBytes)
}

func (cc *SFDSmartContract) debitAccount(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting UserID")
	}
	id := args[0]
	amount, err := strconv.ParseFloat(args[1], 32)

	if err != nil {
		return shim.Error(err.Error())
	}

	accountBytes, err := stub.GetState(id)
	account := Account{}
	json.Unmarshal(accountBytes, &account)
	account.Amount = account.Amount - amount

	accountBytes, _ = json.Marshal(account)
	stub.PutState(account.ID, accountBytes)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountBytes)
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func main() {
	err := shim.Start(new(SFDSmartContract))
	if err != nil {
		panic(err)
	}
}

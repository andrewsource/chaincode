/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"strconv"
)

// StorageChaincode example simple Chaincode implementation
type StorageChaincode struct {
}

func (t *StorageChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")
	// if len(args) != 1 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	// }
	//keeperName := args[0]
	var keeper = make(map[string][]string)
	keeperByte, err := json.Marshal(keeper)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("keeper", keeperByte)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *StorageChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	hash := args[0]
	user := args[1]
	Kvalbytes, err := stub.GetState("keeper")
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	Kval := make(map[string][]string)
	err = json.Unmarshal(Kvalbytes, &Kval)
	if err != nil {
		return nil, errors.New("Failed to unmarsal keeper")
	}
	if len(Kval[hash]) > 0 {
		for i := 0; i < len(Kval[hash]); i++ {
			if Kval[hash][i] == user {
				return []byte{0}, nil
			}
		}
		Kval[hash] = append(Kval[hash], user)
	} else {
		Kval[hash] = make([]string, 0)
		Kval[hash] = append(Kval[hash], user)
	}
	//save state
	Kvalbytes, err = json.Marshal(Kval)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("keeper", Kvalbytes)
	if err != nil {
		return nil, err
	}
	return []byte{1}, nil

}

// Deletes an entity from state
func (t *StorageChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *StorageChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *StorageChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (t *StorageChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")
	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments")
	}
	hash := args[0]
	user := args[1]
	Kvalbytes, err := stub.GetState("keeper")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for keeper\"}"
		return nil, errors.New(jsonResp)
	}
	Kval := make(map[string][]string)
	err = json.Unmarshal(Kvalbytes, &Kval)
	if err != nil {
		return nil, errors.New("Failed to unmarsal keeper")
	}
	if len(Kval[hash]) > 0 {
		for i := 0; i < len(Kval[hash]); i++ {
			if Kval[hash][i] == user {
				return []byte{1}, nil
			}
		}
	}
	return []byte{0}, nil
}

func main() {
	err := shim.Start(new(StorageChaincode))
	if err != nil {
		fmt.Printf("Error starting Storage chaincode: %s", err)
	}
}

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
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")

	var Pagador, Recebedor string    // Entities
	var PagadorVal, RecebedorVal int // Asset holdings
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	Pagador = args[0]
	PagadorVal, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	Recebedor = args[2]
	RecebedorVal, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	fmt.Printf("Pagador = %d, Recebedor = %d\n", PagadorVal, RecebedorVal)

	// Write the state to the ledger
	err = stub.PutState(Pagador, []byte(strconv.Itoa(PagadorVal)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Recebedor, []byte(strconv.Itoa(RecebedorVal)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")

	var Pagador, Recebedor string    // Entities
	var PagadorVal, RecebedorVal int // Asset holdings
	var X int                        // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	Pagador = args[0]
	Recebedor = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(Pagador)

	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}

	PagadorVal, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(Recebedor)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	RecebedorVal, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	PagadorVal = PagadorVal - X
	RecebedorVal = RecebedorVal + X
	fmt.Printf("Pagador = %d, Recebedor = %d\n", PagadorVal, RecebedorVal)

	// Write the state back to the ledger
	err = stub.PutState(Pagador, []byte(strconv.Itoa(PagadorVal)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Recebedor, []byte(strconv.Itoa(RecebedorVal)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	Pagador := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(Pagador)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var Pagador string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	Pagador = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(Pagador)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Pagador + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + Pagador + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + Pagador + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

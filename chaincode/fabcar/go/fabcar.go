/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Asset structure, with 3 properties.  Structure tags are used by encoding/json library 
type Asset struct {
	Latitude string `json:"lat"`
	Longitude string `json:"lng"`
	pH  string `json:"ph"`
	Temperature  string `json:"temp"`
	TimeStamp string `json:"timestamp"`
	
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getReadingForID" {
		return s.getReadingForID(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addReading" {
		return s.addReading(APIstub, args)
	} else if function == "getReading" {
		return s.getReading(APIstub)
	} else if function == "updateReading" {
		return s.updateReading(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getReadingForID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(assetAsBytes)
}
//here
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	assets := []Asset{
		Asset{Latitude: "4.2", Longitude: "23.1", pH: "0.9", Temperature: "39.1", TimeStamp: "15/02/2020 15:05"},
		Asset{Latitude: "0.2", Longitude: "3.1", pH: "14", Temperature: "34.1", TimeStamp: "15/02/2020 15:05"},
		Asset{Latitude: "1.2", Longitude: "1.1", pH: "10", Temperature: "32.1", TimeStamp: "15/02/2020 15:05"},
		Asset{Latitude: "3.2", Longitude: "19.1", pH: "9", Temperature: "29.1", TimeStamp: "15/02/2020 15:05"},
	}

	i := 0
	for i < len(assets) {
		fmt.Println("i is ", i)
		assetAsBytes, _ := json.Marshal(assets[i])
		APIstub.PutState("ASSET_"+strconv.Itoa(i), assetAsBytes)
		fmt.Println("Added", assets[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) addReading(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
//here
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var asset = Asset{Latitude: args[1], Longitude: args[2], pH: args[3], Temperature: args[4], TimeStamp: args[5]}

	assetAsBytes, _ := json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getReading(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "ASSET_0"
	endKey := "ASSET_999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record values\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllReadings:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) updateReading(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}


	assetAsBytes, _ := APIstub.GetState(args[0])
	asset := Asset{}

	json.Unmarshal(assetAsBytes, &asset)

	asset.Latitude = args[1]
	asset.Longitude = args[2]
    asset.pH = args[3]
    asset.Temperature = args[4]
    asset.TimeStamp = args[5]


	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

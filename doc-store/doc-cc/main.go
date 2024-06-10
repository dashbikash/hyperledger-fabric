package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}
type Document struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Content    string `json:"Content"`
	Owner      string `json:"Owner"`
	Size       int64  `json:"Size"`
	CreatedAt  int64  `json:"CreatedAt"`
	CreatedBy  string `json:"CreatedBy"`
	ModifiedAt int64  `json:"ModifiedAt"`
	ModifiedBy string `json:"ModifiedBy"`
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) DocExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) CreateDocument(ctx contractapi.TransactionContextInterface, id string, content string, owner string) (*Document, error) {
	isExist, err := s.DocExists(ctx, id)
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return nil, err
	}
	if !isExist {
		doc := &Document{
			ID:        id,
			Content:   content,
			Owner:     owner,
			CreatedAt: time.Now().UnixMilli(),
			CreatedBy: owner,
		}

		docJSON, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}
		ctx.GetStub().PutState(id, docJSON)
		return doc, nil
	}
	return nil, nil
}

func (s *SmartContract) GetDocument(ctx contractapi.TransactionContextInterface, id string) (*Document, error) {

	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return nil, err
	}
	var doc Document
	err = json.Unmarshal(docBytes, &doc)

	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *SmartContract) DeleteDocument(ctx contractapi.TransactionContextInterface, id string) error {
	isExist, err := s.DocExists(ctx, id)
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return err
	}
	if isExist {
		err := ctx.GetStub().DelState(id)
		if err != nil {
			fmt.Printf("failed to read from world state: %v", err)
			return err
		}

		return nil
	}
	return nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating document store chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting document store chaincode: %v", err)
	}
}

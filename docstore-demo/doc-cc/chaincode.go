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
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
	Owner       string `json:"Owner"`
	Size        int64  `json:"Size"`
	ModifiedAt  int64  `json:"ModifiedAt"`
	ModifiedBy  string `json:"ModifiedBy"`
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) DocExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) CreateDocument(ctx contractapi.TransactionContextInterface, id string, content string, user string) (*Document, error) {
	isExist, err := s.DocExists(ctx, id)
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return nil, err
	}
	if !isExist {
		doc := &Document{
			ID:      id,
			Content: content,
			Owner:   user,
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

func (s *SmartContract) UpdateDocument(ctx contractapi.TransactionContextInterface, id string, content string, user string) (*Document, error) {
	isExist, err := s.DocExists(ctx, id)
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return nil, err
	}
	if isExist {
		doc, err := s.GetDocument(ctx, id)
		if err != nil {
			return nil, err
		}
		doc.Content = content
		doc.ModifiedAt = time.Now().UnixMilli()
		doc.ModifiedBy = user

		docJSON, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}
		ctx.GetStub().PutState(id, docJSON)
		return doc, nil
	}
	return nil, nil
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

func (s *SmartContract) GetAll(ctx contractapi.TransactionContextInterface) ([]*Document, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		fmt.Printf("failed to read from world state: %v", err)
		return nil, err
	}
	defer resultsIterator.Close()

	var docs []*Document
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var doc Document
		err = json.Unmarshal(queryResponse.Value, &doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
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

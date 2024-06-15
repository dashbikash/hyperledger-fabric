package main

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

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

// Document CRUD methods
func CreateDocument(contract *client.Contract, id string, content string, owner string) []byte {
	evaluateResult, err := contract.SubmitTransaction("CreateDocument", id, content, owner)
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
	return evaluateResult
}
func GetDocument(contract *client.Contract, id string) []byte {
	evaluateResult, err := contract.EvaluateTransaction("GetDocument", id)
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
	return evaluateResult
}

func GetAll(contract *client.Contract) []byte {
	evaluateResult, err := contract.EvaluateTransaction("GetAll")
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
	return evaluateResult
}

func UpdateDocument(contract *client.Contract, id string, content string, owner string) {
	evaluateResult, err := contract.SubmitTransaction("UpdateDocument", id, content, owner)
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
}
func DeleteDocument(contract *client.Contract, id string) []byte {
	evaluateResult, err := contract.SubmitTransaction("DeleteDocument", id)
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
	return evaluateResult
}

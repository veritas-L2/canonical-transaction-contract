package ctc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CanonicalTransactionContract struct {}

type TransactionInfo struct {
	ChaincodeName   string   `json:"chaincodeName"`
	TransactionName string   `json:"transactionName"`
	Args            []string `json:"args"`
}

type Batch struct {
    Timestamp     int64             `json:"timestamp"`
    Transactions  []TransactionInfo `json:"transactions"`
    PrevStateHash []byte            `json:"prevStateHash"`
    NewStateHash  []byte            `json:"newStateHash"`
}

const genesisKey = "genesis"

func (ctc *CanonicalTransactionContract) PublishToState(ctx contractapi.TransactionContextInterface, payload string) (error) {
	val, _ := ctx.GetStub().GetState(genesisKey)
	key := genesisKey

	if (val != nil){
		var data Batch
		if err := json.Unmarshal([]byte(payload), &data); err != nil{
			return fmt.Errorf("failed to unmarshal payload for genesis: %s", string(key))
		} else {
			key = string(data.PrevStateHash)
		}
	}

	return ctx.GetStub().PutState(key, []byte(payload))
}

func (ctc *CanonicalTransactionContract) GetBatchHistory(ctx contractapi.TransactionContextInterface) (string, error){
	var history []string;

	genesisVal, _ := ctx.GetStub().GetState(genesisKey)
	if (genesisVal == nil){
		return "", nil
	}

	var genesisBatch Batch
	json.Unmarshal([]byte(genesisVal), &genesisBatch)
	history = append(history, string(genesisVal))

	for val, _ := ctx.GetStub().GetState(string(genesisBatch.NewStateHash)); val != nil ; {
		var batch Batch
		json.Unmarshal([]byte(val), &batch)
		history = append(history, string(val))

		val, _ = ctx.GetStub().GetState(string(batch.NewStateHash))
	}

	
	return strings.Join(history, ";"), nil
}
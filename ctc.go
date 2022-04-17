package ctc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CanonicalTransactionContract struct {
	contractapi.Contract
}

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

type BatchReadyForCommit struct {
	Timestamp     int64             `json:"timestamp"`
	Transactions  []TransactionInfo `json:"transactions"`
	PrevStateHash []byte            `json:"prevStateHash"`
	NewStateHash  []byte            `json:"newStateHash"`
	PrevBatchHash []byte            `json:"prevBatchHash"`
}

const genesisKey = "genesis"
const latestBatchHashKey = "latestBatchHash"

func (ctc *CanonicalTransactionContract) Init(ctx contractapi.TransactionContextInterface) error {
	genesisKeyHash := crypto.Keccak256([]byte(genesisKey))
	existingGenesisBatch, err := ctx.GetStub().GetState(hex.EncodeToString(genesisKeyHash))
	if err != nil {
		return fmt.Errorf("failed to check if genesis batch already exists in world state. Error: %s", err.Error())
	}

	if existingGenesisBatch == nil {
		return fmt.Errorf("genesis batch already exists in world state.")
	}

	//setup genesis block and latestBatchHash

	//Thought: we probably do not need a genesis batch, a batch that comes in when no latest batch hash exists can
	//serve as the genesis batch
	genesisBatch := BatchReadyForCommit{
		Timestamp:     time.Now().Unix(),
		Transactions:  []TransactionInfo{},
		PrevStateHash: nil,
		NewStateHash:  nil,
		PrevBatchHash: nil,
	}

	genesisBatchJSON, _ := json.Marshal(genesisBatch)

	err = ctx.GetStub().PutState(hex.EncodeToString(genesisKeyHash), genesisBatchJSON)
	if err != nil {
		return fmt.Errorf("failed to add genesis batch into world state. Error: %s", err.Error())
	}

	err = ctx.GetStub().PutState(latestBatchHashKey, genesisKeyHash)
	if err != nil {
		return fmt.Errorf("failed to add latestBatchHash into world state. Error: %s", err.Error())
	}

	// are state changes atomic?

	return nil
}

func (ctc *CanonicalTransactionContract) PublishToState(ctx contractapi.TransactionContextInterface, payload string) error {
	var batch Batch
	if err := json.Unmarshal([]byte(payload), &batch); err != nil {
		return fmt.Errorf("failed to unmarshal batch payload. Error: %s", err.Error())
	}

	latestBatchHash, err := ctx.GetStub().GetState(latestBatchHashKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve latest batch hash from world state. Error: %s", err.Error())
	}

	if latestBatchHash == nil {
		return fmt.Errorf("Latest batch hash not found. Are you sure you have called Init()?")
	}

	//TODO: check if prevStateHash matches newStateHash of latest batch

	batchReadyForCommit := BatchReadyForCommit{
		Timestamp:     batch.Timestamp,
		Transactions:  batch.Transactions,
		PrevStateHash: batch.PrevStateHash,
		NewStateHash:  batch.NewStateHash,
		PrevBatchHash: latestBatchHash,
	}

	if batchJSON, err := json.Marshal(batchReadyForCommit); err == nil {
		batchHash := crypto.Keccak256(batchJSON)

		return ctx.GetStub().PutState(hex.EncodeToString(batchHash), batchJSON)
	} else {
		return fmt.Errorf("failed to marshal batch struct. Error: %s", err.Error())
	}
}

func (ctc *CanonicalTransactionContract) GetBatchHistory(ctx contractapi.TransactionContextInterface) (string, error) {
	var history []string

	//TODO: handle errors

	latestBatchHash, _ := ctx.GetStub().GetState("latestBatchHash")
	if latestBatchHash == nil {
		return "", fmt.Errorf("latest batch hash does not exist in world state. Are you sure you called Init()?")
	}

	rawBatch, _ := ctx.GetStub().GetState(hex.EncodeToString(latestBatchHash))

	for rawBatch != nil {
		var batch BatchReadyForCommit
		json.Unmarshal(rawBatch, &batch)
		history = append(history, string(rawBatch))

		rawBatch, _ = ctx.GetStub().GetState(hex.EncodeToString(batch.PrevBatchHash))
	}

	return strings.Join(history, ";"), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(CanonicalTransactionContract))

	if err != nil {
		fmt.Printf("Error create statecontract chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting statecontract chaincode: %s", err.Error())
	}
}

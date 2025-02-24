package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	Proof        int
	PreviousHash string
	Hash         string
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}
type Blockchain struct {
	Chain        []Block
	Transactions []Transaction
}

func main() {
	bc := createBlockchain()
	bc.addTransaction("Alice", "Bob", 50)
	bc.addTransaction("Bob", "Charlie", 25)

	start := time.Now()
	proof, candidateTimestamp := bc.proofOfWork()

	duration := time.Since(start)
	fmt.Printf("Proof of work found in %d (execution time: %s)\n", proof, duration)
	previousHash := bc.Chain[len(bc.Chain)-1].Hash
	bc.addBlock(proof, candidateTimestamp, previousHash)

	fmt.Println("Blockchain:", bc.Chain)
}

// create and returns new blockchain with empty slices and first genesis block
func createBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:        []Block{},
		Transactions: []Transaction{},
	}

	bc.createGenesisBlock() // genesis block
	return bc
}

// creates first block in a chain
func (bc *Blockchain) createGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		Proof:        100,
		PreviousHash: "0",
	}

	genesisBlock.Hash = calculateHash(genesisBlock)

	bc.Chain = append(bc.Chain, genesisBlock)
}

// adds new block to blockchain, clears transactions
func (bc *Blockchain) addBlock(proof int, timestamp int64, previousHash string) {
	newBlock := Block{
		Index:        len(bc.Chain),
		Timestamp:    timestamp,
		Transactions: bc.Transactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	newBlock.Hash = calculateHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)
	bc.Transactions = []Transaction{}
}

// adds transaction, returns block index
func (bc *Blockchain) addTransaction(sender, recipient string, amount float64) string {
	tx := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	bc.Transactions = append(bc.Transactions, tx)

	// TXID
	return generateTransactionID(tx)
}

// generates unique transaction id
func generateTransactionID(tx Transaction) string {
	data := fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func isProofValid(lastBlock Block, proof int, transactions []Transaction, candidateTimestamp int64) bool {
	tempBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    candidateTimestamp,
		Transactions: transactions,
		Proof:        proof,
		PreviousHash: lastBlock.Hash,
	}

	guessHash := calculateHash(tempBlock)
	return guessHash[:4] == "0000"
}

func (bc *Blockchain) proofOfWork() (int, int64) {
	lastBlock := bc.Chain[len(bc.Chain)-1]

	candidateTimestamp := time.Now().Unix()
	proof := 0
	for !isProofValid(lastBlock, proof, bc.Transactions, candidateTimestamp) {
		proof++
	}

	return proof, candidateTimestamp
}

// calculates hash of a block
func calculateHash(block Block) string {

	hashInput := fmt.Sprintf("%d%d%d%s%d",
		block.Index, block.Timestamp, block.Proof, block.PreviousHash,
		len(block.Transactions))

	for _, tx := range block.Transactions {
		hashInput += tx.Sender + tx.Recipient + fmt.Sprintf("%f", tx.Amount)
	}

	hash := sha256.Sum256([]byte(hashInput))

	return hex.EncodeToString(hash[:])
}

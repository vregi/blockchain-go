package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block structure contains the index of the block, timestamp of the block, slice of confirmed transactions, proof-of-work (nonce), hash of the previous block and own hash .
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	Nonce        int // nonce
	PreviousHash string
	Hash         string
}

// Transaction structure contains the sender, recipient and amount of medium's of exchange unit.
type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	TXID      string // Transaction ID
}

// Blockchain structure contains the slice of blocks which instantiates the blockchain itself and slice of transaction, which is needed for the temporary pool of unconfirmed transactions - "mempool".
type Blockchain struct {
	Chain        []Block
	Transactions []Transaction // mempool
}

func main() {
	bc := createBlockchain()
	bc.addTransaction("Alice", "Bob", 50)
	bc.addTransaction("Bob", "Charlie", 25)

	start := time.Now()
	nonce, candidateTimestamp := bc.proofOfWork()

	duration := time.Since(start)
	fmt.Printf("Proof of work (nonce) found in iteration %d (execution time: %s)\n", nonce, duration)
	previousHash := bc.Chain[len(bc.Chain)-1].Hash
	bc.addBlock(nonce, candidateTimestamp, previousHash)

	fmt.Println("Blockchain:", bc.Chain)
}

// createBlockchain initializes and returns a new Blockchain instance
// with the genesis block already added to the chain
func createBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:        []Block{},
		Transactions: []Transaction{},
	}

	bc.createGenesisBlock() // genesis block
	return bc
}

// createGenesisBlock creates the very first block of the blockchain (genesis block),
// sets its predefined values, calculates its hash, and appends it to the chain
func (bc *Blockchain) createGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		Nonce:        100,
		PreviousHash: "0",
	}

	genesisBlock.Hash = calculateHash(genesisBlock)

	bc.Chain = append(bc.Chain, genesisBlock)
}

// addBlock creates a new block using the provided nonce, timestamp, and previous hash,
// calculates its hash, appends it to the chain, and clears the mempool
func (bc *Blockchain) addBlock(nonce int, timestamp int64, previousHash string) {
	newBlock := Block{
		Index:        len(bc.Chain),
		Timestamp:    timestamp,
		Transactions: bc.Transactions,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}
	newBlock.Hash = calculateHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)
	bc.Transactions = []Transaction{}
}

// addTransaction adds an unconfirmed transaction to the mempool
// and returns a unique transaction ID generated from its contents
func (bc *Blockchain) addTransaction(sender, recipient string, amount float64) string {
	tx := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	tx.TXID = generateTransactionID(tx)

	bc.Transactions = append(bc.Transactions, tx)

	return tx.TXID
}

// generateTransactionID creates a SHA-256 hash from a transaction's sender, recipient,
// and amount to uniquely identify the transaction and prevent duplication or tampering
func generateTransactionID(tx Transaction) string {
	data := fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// isProofValid verifies whether the hash generated from a block candidate with a given nonce
// satisfies the mining difficulty condition (e.g., hash starts with "0000")
func isProofValid(lastBlock Block, nonce int, transactions []Transaction, candidateTimestamp int64) bool {
	tempBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    candidateTimestamp,
		Transactions: transactions,
		Nonce:        nonce, // nonce
		PreviousHash: lastBlock.Hash,
	}

	guessHash := calculateHash(tempBlock)
	return guessHash[:4] == "0000" // mining difficulty target
}

// proofOfWork iterates over increasing nonce values, generating a hash each time,
// until it finds a hash that satisfies the predefined mining difficulty (e.g. starts with "0000").
// Returns the valid nonce and the timestamp when the proof was found
func (bc *Blockchain) proofOfWork() (int, int64) {
	lastBlock := bc.Chain[len(bc.Chain)-1]

	candidateTimestamp := time.Now().Unix()
	nonce := 0
	for !isProofValid(lastBlock, nonce, bc.Transactions, candidateTimestamp) {
		nonce++
	}

	return nonce, candidateTimestamp
}

// calculateHash generates the SHA-256 hash of a block by concatenating its index, timestamp, nonce,
// previous block's hash, number of transactions, and details of each transaction (sender, recipient, amount).
// Returns the hexadecimal string representation of the resulting hash.
func calculateHash(block Block) string {

	hashInput := fmt.Sprintf("%d%d%d%s%d",
		block.Index, block.Timestamp, block.Nonce, block.PreviousHash,
		len(block.Transactions))

	for _, tx := range block.Transactions {
		hashInput += tx.Sender + tx.Recipient + fmt.Sprintf("%f", tx.Amount)
	}

	hash := sha256.Sum256([]byte(hashInput))

	return hex.EncodeToString(hash[:])
}

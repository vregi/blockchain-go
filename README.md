# minimum viable blockchain in go

this is a minimum vialbe blockchain implementation written in go, created for **educational purposes** to learn the fundamental concepts behind blockchain technology

it is loosely based on the [minimum viable blockchain](https://www.igvita.com/2014/05/05/minimum-viable-block-chain) article

## purpose

the goal of this project is to demonstrate how a basic blockchain works: how blocks are created, how transactions are handled, how hashes are used to link blocks, and how a simple proof-of-work system can simulate mining

## features

- basic `block` and `transaction` structures
- chain of blocks with hashes linking them
- simple mempool (temporary pool of transactions)
- block creation and hash generation
- proof-of-work algorithm (mining by finding a hash starting with "0000")
- transaction id (txid) based on hashed contents

## how it works

1. the blockchain is initialized with a **genesis block**
2. users can add transactions to the **mempool**
3. a simple **proof-of-work** function tries increasing nonce values until a block hash starting with `"0000"` is found
4. once found, a new block is created with that nonce and the transactions are added to it
5. the mempool is cleared and the new block is appended to the chain

## main concepts

- **block**: contains index, timestamp, transaction list, nonce, previous block hash, and its own hash
- **transaction**: has sender, recipient, amount, and a generated transaction id
- **hashing**: each block's hash is based on its contents, ensuring immutability
- **proof-of-work**: a basic mining simulation where we look for a hash with a prefix of "0000"
- **mempool**: stores unconfirmed transactions waiting to be added to the next block

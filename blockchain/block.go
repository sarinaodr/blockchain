package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block represents a single block in the blockchain
// It contains the block's hash, data, and the hash of the previous block
type Block struct {
	Timestamp   int64          // Timestamp of the block
	Hash        []byte         // SHA256 hash of the block's data and previous hash
	Transaction []*Transaction // Data stored in the block
	PrevHash    []byte         // Hash of the previous block in the chain
	Nonce       int
	Height      int
}

// DeriveHash calculates and sets the block's hash based on its data and previous hash
// The hash is computed using SHA256 on the concatenation of data and previous hash

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transaction {
		txHashes = append(txHashes, tx.Serialize())
	}

	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}

// CreateBlock creates a new block with the given data and previous hash
// It automatically calculates and sets the block's hash
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis creates the first block in a blockchain (genesis block)
// It contains the string "Genesis" as data and an empty previous hash
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{},0)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

package database

import (
	"crypto/sha256"
	"encoding/json"
)
type Hash [32]byte

type Block struct {
	Header BlockHeader `json:"header"` // reference to the parent block meta data
	TXs    []Tx        `json:"payload"`//new transactions only (payload)
}

type BlockHeader struct {
	Parent Hash   `json:"parent"` //parent block reference
	Time   uint64 `json:"time"`//time parent block was made
}

type BlockFS struct {
	Key Hash `json:"hash"`
	value Block `json:"block"`
}

func NewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{BlockHeader{parent, time}, txs};
}

func (b Block) Hash (Hash, error){
	blockJson, err := json:Marshal(b);

	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}


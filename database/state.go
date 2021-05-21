package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
	snapshot  Snapshot
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")

	gen, err := LoadGenesis(genFilePath)

	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txFilePath, os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{balances, make([]Tx, 0), f, Snapshot{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		if err := json.Unmarshal(scanner.Bytes(), &tx); err != nil {
			return nil, err
		}

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() (Hash, error) {
	//create a new Block with ONLY new TXs
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)

	blockHash, err := block.Hash()

	if err != nil {
		return Hash{}, err
	}

	blockFS := BlockFS{blockHash, block}

	blockFSJson, err := json.Marshal(blockFS)

	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFSJson)


	err = s.dbFile.Write( append(blockFSJson,"\n") )
	if err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash

	s.txMempool = []Tx
	
	return s.latestBlockHash, nil
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) apply(tx Tx) error {
	// policy guide for applying a transaction

	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) doSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return nil
	}

	txData, err := ioutil.ReadAll(s.dbFile)

	if err != nil {
		return err
	}

	s.snapshot = sha256.Sum256(txData)

	return nil
}

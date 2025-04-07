package blockchain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	// "github.com/Kami0rn/golang-blockchain/blockchain"
	"github.com/dgraph-io/badger"
	// "google.golang.org/grpc/encoding"
)

const (
	dbPath = "./tmp/blocks"
	dbFile = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database *badger.DB
}

func DBexists() bool {
	if _,err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}

	return true
}

func InitBlockChain() *BlockChain {
	var lastHash [] byte

	if DBexists(){
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created!")
		err = txn.Set(genesis.Hash,genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"),genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash,db}

	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash [] byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err
		
	})

	Handle(err)

	chain := BlockChain{lastHash,db}

	return &chain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error{
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)

}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash,chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
    var block *Block
    err := iter.Database.View(func(txn *badger.Txn) error {
        item, err := txn.Get(iter.CurrentHash)
        if err != nil {
            return err
        }
        return item.Value(func(val []byte) error {
            block = Deserialize(val)
            return nil
        })
    })
    Handle(err)

    iter.CurrentHash = block.PrevHash

    return block
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spenTX0s := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Outputs {
					if spenTX0s[txID] != nil {
						for _, spentOut := range spenTX0s[txID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					if out.CanBeUnlocked(address) {
						unspentTxs = append(unspentTxs, *tx)
					}
				}
		}

		if len(block.PrevHash) == 0{
			break
		}
	}

	return unspentTxs
}
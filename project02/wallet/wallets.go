package wallet

import (
	"crypto/elliptic"
	"encoding/gob"
	"log"
)

const walletFile = "./tmp/wallet.data"

type Wallets struct {
	Wallet map[string]*Wallet
}

func (ws *Wallets) SaveFile() {
	var content byte.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
}

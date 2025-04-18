package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

const walletFile = "./tmp/wallet.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

type SerializableWallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var serializableWallets map[string]SerializableWallet
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&serializableWallets)
	if err != nil {
		return err
	}

	ws.Wallets = make(map[string]*Wallet)
	for address, sWallet := range serializableWallets {
		privateKey := ecdsa.PrivateKey{
			D: new(big.Int).SetBytes(sWallet.PrivateKey),
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     new(big.Int).SetBytes(sWallet.PublicKey[:len(sWallet.PublicKey)/2]),
				Y:     new(big.Int).SetBytes(sWallet.PublicKey[len(sWallet.PublicKey)/2:]),
			},
		}
		ws.Wallets[address] = &Wallet{
			PrivateKey: privateKey,
			PublicKey:  sWallet.PublicKey,
		}
	}

	return nil
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer
	serializableWallets := make(map[string]SerializableWallet)

	for address, wallet := range ws.Wallets {
		privateKeyBytes := wallet.PrivateKey.D.Bytes()
		serializableWallets[address] = SerializableWallet{
			PrivateKey: privateKeyBytes,
			PublicKey:  wallet.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(serializableWallets)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

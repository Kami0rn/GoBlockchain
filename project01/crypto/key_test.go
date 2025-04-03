package crypto

import (
	// "encoding/hex"
	"fmt"
	// "io"
	// "crypto/rand"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)
	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar buz")

	sig := privKey.Sign(msg)
	assert.True(t,sig.Verify(pubKey, msg))

	// Test with invalid msg
	assert.False(t,sig.Verify(pubKey, []byte("foo")))

	// Test with invalid key
	invalidPrivKey := GeneratePrivateKey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
	

}

func TestNewPrivateKeyFromString(t *testing.T) {
	var  (
		seed = "29ae93b1f8a5d2780407280e38e0d7d72acf7be19f7bf4ed47d56d2d0b4acba2"
		privKey = NewPrivatKeyFromString(seed)
		addressStr = "2a7c4fdbb83dc7d98609441ea7870854c1dd77e3"
	)
	
	assert.Equal(t, privKeyLen , len(privKey.Bytes()))
	address := privKey.Public().Address()
	assert.Equal(t , addressStr, address.String())


	// seed := make([]byte,32)
	// io.ReadFull(rand.Reader, seed)
	// fmt.Println(hex.EncodeToString(seed))
}

func TestPublicKeyToAdress(t *testing.T) {
	privKey := GeneratePrivateKey()
	publicKey := privKey.Public()
	address := publicKey.Address()
	assert.Equal(t, addressLen, len(address.Bytes()))
	fmt.Println(address)
}
package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/vrecan/death/v3"

	"github.com/Kami0rn/golang-blockchain/blockchain"
)

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAddress    string
	minerAddress   string
	KnownNodes     = []string{"localhost:3000"}
	blockInTransit = [][]byte{}
	memoryPool     = make(map[string]blockchain.Transaction)
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}

type GetBlocks struct {
	AddrFrom string
}

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func CmdToBytes (cmd string) []byte {
	var bytes [commandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCmd(bytes []byte) string{
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}


func ExtractCmd(request []byte) []byte {
	return request[:commandLength]
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	SendData(address,request)
}

func SendBlock(addr string,b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"),payload...)

	SendData(addr, request)
}

func SendData(addr string, data []byte){
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("%s is not availabl\n", addr)
		var updateNodes []string

		for _, node := range KnownNodes{
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}

		KnownNodes = updateNodes

		return
		
	}

	defer conn.Close()

	_,err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func ()  {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain){
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()

	if err != nil {
		log.Panic(err)
	}
	command := BytesToCmd(req[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	default:
		fmt.Println("Unknown command")
	}
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

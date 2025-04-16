package main

import (
	"os"

	// "github.com/Kami0rn/golang-blockchain/cli"
	"github.com/Kami0rn/golang-blockchain/wallet"
)

// "google.golang.org/grpc/balancer"



func main() {
	defer os.Exit(0)
	// cli := cli.CommandLine{}
	// cli.Run()

	w := wallet.MakeWallet()
	w.Address()
}
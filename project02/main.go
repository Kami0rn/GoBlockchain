package main

import (
	"os"
	"github.com/Kami0rn/golang-blockchain/cli"
)
// "google.golang.org/grpc/balancer"



func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
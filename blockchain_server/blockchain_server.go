package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/leoashishs99/go-blockchain/block"
	"github.com/leoashishs99/go-blockchain/wallet"
)

var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type BlockchainServer struct {
	port uint64
}

func (bcs *BlockchainServer) Port() uint64 {
	return bcs.port
}
func NewBlockChainServer(port uint64) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) GetBlockChain() *block.BlockChain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockChain(minersWallet.BlockChainAddress(), bcs.Port())
		cache["blockchain"] = bc
	}
	return bc
}
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!!")
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", HelloWorld)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(bcs.port)), nil))
}

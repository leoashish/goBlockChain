package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/leoashishs99/go-blockchain/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {

	fmt.Printf("timestamp: %d\n", b.timestamp)
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Printf("previousHash: %s\n", b.previousHash)
	for i, tran := range b.transactions {
		fmt.Sprintf("%s %d %s", strings.Repeat("-", 6), i, strings.Repeat("-", 6))
		tran.Print()
	}
	fmt.Println("-----------------------------------\n")
}

func (b *Block) Hash() [32]byte {
	marshalledJson, _ := b.MarshalJSON()
	//fmt.Println(string(marshalledJson))
	return sha256.Sum256([]byte(marshalledJson))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TimeStamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json: "transactions"`
	}{
		TimeStamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

type BlockChain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint64
}

func NewBlockChain(blockchainAddress string, port uint64) *BlockChain {
	b := new(Block)
	bc := new(BlockChain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}
func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, block)
	bc.transactionPool = []*Transaction{}
	return block
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 13), i, strings.Repeat("=", 13))
		block.Print()
	}
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float64
}

func NewTrasaction(sender string, recipient string, value float64) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}
func (t *Transaction) Print() {
	fmt.Println(strings.Repeat("=", 26))
	fmt.Printf("Sender %s \n", t.senderBlockchainAddress)
	fmt.Printf("Reciever %s \n", t.recipientBlockchainAddress)
	fmt.Printf("Value %.1f \n", t.value)

	fmt.Println(strings.Repeat("-", 26))
}

func (t *Transaction) MarshallJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float64 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
func (bc *BlockChain) AddTransaction(sender string, recipient string, value float64,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTrasaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("Error: Insufficient balance!!")
		// 	return true
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("Error: Verify Transaction!!!")
	}
	return false
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range bc.transactionPool {
		transactions = append(transactions,
			NewTrasaction(transaction.senderBlockchainAddress,
				transaction.recipientBlockchainAddress,
				transaction.value))
	}

	return transactions
}

func (bc *BlockChain) VerifyTransactionSignature(
	publicKey *ecdsa.PublicKey, s *utils.Signature, transaction *Transaction) bool {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(publicKey, h[:], s.R, s.S)
}
func (bc *BlockChain) ValidProof(difficultyLevel int, nonce int, previousHash [32]byte, trasactions []*Transaction) bool {
	zeros := strings.Repeat("0", 3)
	guessBlock := Block{0, nonce, previousHash, trasactions}
	guessString := fmt.Sprintf("%x", guessBlock.Hash())
	fmt.Println(guessString)
	return guessString[:difficultyLevel] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	lastTransactionBlock := bc.LastBlock()
	nonce := 0

	for !bc.ValidProof(MINING_DIFFICULTY, nonce, lastTransactionBlock.Hash(), transactions) {
		nonce += 1
	}

	return nonce
}

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nounce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nounce, previousHash)
	log.Println("Action= mining, status= Succees")
	return true
}

func (bc *BlockChain) CalculateTotalAmount(userBlockChainAddress string) float64 {
	totalAmount := 0.0
	for _, perBlockTransactions := range bc.chain {
		for _, transaction := range perBlockTransactions.transactions {
			if transaction.recipientBlockchainAddress == userBlockChainAddress {
				totalAmount += transaction.value
			}
			if transaction.senderBlockchainAddress == userBlockChainAddress {
				totalAmount -= transaction.value
			}
		}
	}

	return totalAmount
}

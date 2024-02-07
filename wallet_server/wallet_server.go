package main

import (
	"encoding/json"
	"fmt"
	"golang-blockchain/utils"
	"golang-blockchain/wallet"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

const tempDir = "/home/ubuntu/Desktop/GO/golang-blockchain/wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

//	func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
//		switch req.Method {
//		case http.MethodGet:
//			t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
//			t.Execute(w, "")
//		default:
//			log.Printf("ERROR: Invalid HTTP Method")
//		}
//	}
func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir, "index.html"))
		if err != nil {
			log.Printf("ERROR: Failed to parse template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, "")
		if err != nil {
			log.Printf("ERROR: Failed to execute template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("FAIL")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: MISSING FIELDS")
			io.WriteString(w, string(utils.JsonStatus("FAIL")))
			return
		}

		fmt.Println(*t.SenderPublicKey)
		fmt.Println(*t.SenderBlockchainAddress)
		fmt.Println(*t.SenderPrivateKey)
		fmt.Println(*t.RecipeintBlockchainAddress)
		fmt.Println(*t.Value)

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Request")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}

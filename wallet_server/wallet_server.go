package main

import (
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

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
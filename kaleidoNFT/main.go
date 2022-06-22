package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	// _ "github.com/pdrum/swagger-automation/docs"
)

var db *sql.DB

func main() {

	//Create an instance to connect to the postgres DB
	var err error
	connStr := "host=localhost port=5432 dbname=kaleido user=postgres password=test connect_timeout=10 sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database", err)
		panic(err)
	}
	fmt.Println("Successfully connected to database!")

	//connection details to the host
	endpoint := "wss://eth-rinkeby.alchemyapi.io/v2/udItV2soFnlRaibwBUbdIUibJ1-88Ali"
	privateKey := "686981ba439b8ef0b8da5cf474d6751b2da6a51fb95b9246d2758dd931e83219"
	chainId := big.NewInt(4)

	// connecting to the alchemy api
	client, err := NewClient(endpoint, privateKey, chainId)
	if err != nil {
		fmt.Println("Error connecting to blockchain network", err)
		panic(err)
	}

	client.SetNonce(big.NewInt(16789))
	client.SetFundValue(big.NewInt(0))
	client.SetGasLimit(uint64(8000000))
	client.SetGasPrice(big.NewInt(1875000000))

	//deploying the smart contract
	contract, err := client.DeployContract()
	if err != nil {
		fmt.Println("Error connecting to blockchain network", err)
		panic(err)
	}
	fmt.Println("Contract address:", contract.Address.Hex())

	//mint an NFT for the given tokenURI
	client.SetNonce(big.NewInt(167500))
	mintTX, err := contract.MintToken(client, "1")
	if err != nil {
		fmt.Println("Error connecting to blockchain network", err)
		panic(err)
	}
	fmt.Println("Mint transaction:", mintTX)

	//burn the given tokenID NFT
	client.SetNonce(big.NewInt(167987))
	burnTX, err := contract.BurnToken(client, big.NewInt(1))
	if err != nil {
		fmt.Println("Error connecting to blockchain network", err)
		panic(err)
	}
	fmt.Println("Mint transaction:", burnTX.Hash())

	// Create a mux router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/kaleido/{module}", ClientAPIHandler)

	fmt.Println("Successfully connected to localhost!")
	log.Fatal(http.ListenAndServe(":8080", r))

}

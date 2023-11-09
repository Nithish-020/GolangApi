package main

import (
	"apipkg/apdb"
	"log"
)

func main() {
	log.Println("Server started")
	db, err := apdb.LocalDbConnect(apdb.IPODB)
	if err != nil {
		log.Println("Error", err)
	} else {
		defer db.Close()
	}
}

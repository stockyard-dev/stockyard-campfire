package main

import (
	"flag"
	"fmt"
	"github.com/stockyard-dev/stockyard-campfire/internal/server"
	"github.com/stockyard-dev/stockyard-campfire/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "")
	dataFlag := flag.String("data", "", "")
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9700"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./campfire-data"
	}
	if *portFlag != "" {
		port = *portFlag
	}
	if *dataFlag != "" {
		dataDir = *dataFlag
	}
	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("campfire: %v", err)
	}
	defer db.Close()
	srv := server.New(db, server.DefaultLimits(), dataDir)
	fmt.Printf("\n  Campfire — Self-hosted async team standup tool\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n", port, port)
	log.Printf("campfire: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}

package main

import (
	"log"
	"net/http"

	"github.com/nytro04/Client-App/v5/database"
	"github.com/nytro04/Client-App/v5/handlers"
)

func main() {

	//connecting to database
	log.Println("Connecting to database")
	db, err := database.New("user=francisbadasu password=superman host=localhost dbname=clientApp sslmode=disable")
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err)
	}
	defer db.Close()

	//multiplexer
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("templates/assets"))

	//to handle static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", files))

	h := handlers.Handlers{Db: db}

	//routes
	// mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/", h.ShowClients)
	mux.HandleFunc("/CreateClients", h.CreateClients)
	mux.HandleFunc("/DeleteClient", h.DeleteClient)
	mux.HandleFunc("/UpdateClient", h.UpdateClient)
	mux.HandleFunc("/Terminate", h.TerminateClient)	
	mux.HandleFunc("/ShowDetails", h.ShowDetails)

	log.Println("Server started on port....8888")
	log.Fatal(http.ListenAndServe(":8888", mux))
}

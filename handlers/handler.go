package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nytro04/Client-App/v5/clients"
	"github.com/nytro04/Client-App/v5/database"
	"strconv"
)

type Handlers struct {
	Db database.DB
}

//to handle landing page
func (h Handlers) Index(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Fatalf("Error parsing templates: %s\n", err)
	}
	templ.Execute(w, nil)
}

//to handle show clients page
func (h Handlers) ShowClients(w http.ResponseWriter, r *http.Request) {
	AllClients, err := h.Db.GetAllClients()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(AllClients)


	templ, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Fatalf("Error parsing templates: %s\n", err)
	}

	templ.Execute(w, AllClients)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("List of Clients: ", AllClients)

}

//handles create and save clients
func (h Handlers) CreateClients(w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err, "Cant not parse Create client form")
	}
	
	c := &clients.Client{
		Name: r.PostFormValue("name"),
		PostalAddress:	r.PostFormValue("postal"),
		PhysicalAddress: 	r.PostFormValue("physical"),
		Email:	r.PostFormValue("email"),
		NatureOfBusiness: r.PostFormValue("business"),
		ContactPerson: r.PostFormValue("contact"),
		ContactNumber:	r.PostFormValue("number"),
		NumberOfGuards:	r.PostFormValue("guards"),
		Location:	r.PostFormValue("location"),

	}

	_, err = h.Db.CreateClient(c)
	if err != nil {
		log.Fatal(err)
	}

	templ, err := template.ParseFiles("templates/new.gohtml")
	if err != nil {
		log.Fatalf("Error parsing templates: %s\n", err)
	}
	templ.Execute(w, c)

	fmt.Println(c)

}

//Get client by name 
func (h Handlers) GetClients(w http.ResponseWriter, r *http.Request) {
	c := "Topp"

	searchClient, err := h.Db.GetClientsByName(c)
	if err != nil {
		log.Println(err)
	}

	for _, client := range searchClient {
		if client.Name == c {
			fmt.Println("match found", client.Name, client.NatureOfBusiness)
		}
	}
	
	//commented while testing GetClients with println
	// templ, err := template.ParseFiles("templates/new.gohtml")
	// if err != nil {
	// 	log.Fatalf("Error parsing templates: %s\n", err)
	// }
	// templ.Execute(w, nil)
}

func (h Handlers) DeleteClient(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = h.Db.RemoveClient(id)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	http.Redirect(w, r, "/ShowClients", 302)

}


func (h Handlers) TerminateClient(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		c := &clients.Client{
			ID: id,
		}

		err = h.Db.TerminateClient(c)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	http.Redirect(w, r, "/ShowClients", 302)

}


func (h Handlers) UpdateClient(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	idStr := params.Get("id")
	
	if len(idStr) > 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		c := &clients.Client{
			ID: id,
		}

		err = h.Db.UpdateClient(c)
		if err != nil {
			log.Fatal(err)
			return
		}

		templ, err := template.ParseFiles("templates/new.gohtml")
			if err != nil {
			log.Fatalf("Error parsing templates: %s\n", err)
		}

		templ.Execute(w, nil)

		// http.Redirect(w, r, "/ShowClients", 302)

	}
}

func (h Handlers) ShowDetails(w http.ResponseWriter, r *http.Request) {
	
	idstr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idstr, 10, 64) 
	if err != nil {
		log.Fatal(err)
		return
	}

	client, err := h.Db.GetClientByID(id)
	if err != nil {
		log.Fatal(err)
		return
	}

	templ, err := template.ParseFiles("templates/details.gohtml")
	if err != nil {
	log.Fatalf("Error parsing templates: %s\n", err)
}

	templ.Execute(w, client)

	// fmt.Printf("%T\n", client)
	// fmt.Println(client)

}



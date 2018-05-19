package database

import (
	"fmt"
	"log"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nytro04/Client-App/v5/clients"
	"github.com/nytro04/Client-App/v5/users"
)

type DB interface {
	clients.ClientDB
	users.UserDB
	Close() error
}

type postgresDB struct {
	db *sql.DB
}

//New func connect to db
func New(conn string) (DB, error) {
	//connecting to postgres db
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	//checking if connection is alive
	if err = db.Ping(); err != nil {
		return nil, err
	}

	pg := &postgresDB{db: db}
	err = pg.ensureTables()

	return pg, err
}

//ensures that the various tables are created if they dont exist
func (db *postgresDB) ensureTables() error {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS clients(id SERIAL, name TEXT, postalAddress TEXT, physicalAddress TEXT, email TEXT, natureOfBusiness TEXT, contactPerson TEXT, contactNumber TEXT, numberOfGuards TEXT, location TEXT);")
	if err != nil {
		return err
	}

	_, err = db.db.Exec("CREATE TABLE IF NOT EXISTS users(id SERIAL, name TEXT, hash TEXT);")
	if err != nil {
		return err
	}

	return err
}

func (db *postgresDB) Close() error {
	return db.db.Close()
}

func (db *postgresDB) CreateClient(client *clients.Client) (int64, error) {
	result, err := db.db.Prepare("INSERT INTO clients(name, postaladdress, physicaladdress, email, natureofbusiness, contactperson, contactnumber, numberofguards, location) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;")
	if err != nil {
		return 0, err
	}
	defer result.Close()

	var lastInsertId int64

	err = result.QueryRow(client.Name, client.PostalAddress, client.PhysicalAddress, client.Email, client.NatureOfBusiness, client.ContactPerson, client.ContactNumber, client.NumberOfGuards, client.Location).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, err

}

func (db *postgresDB) GetClientByID(id int64) (*clients.Client, error) {
	fmt.Printf("getting row for %d\n", id)
	row := db.db.QueryRow("SELECT id, name, postaladdress, physicaladdress, email, natureofbusiness, contactperson, contactnumber, numberofguards, location FROM clients WHERE id= $1;", id)
	fmt.Printf("got row %v\n", row)

	c := new(clients.Client)
	err := row.Scan(&c.ID, &c.Name, &c.PostalAddress, &c.PhysicalAddress, &c.Email, &c.NatureOfBusiness, &c.ContactPerson, &c.ContactNumber, &c.NumberOfGuards, &c.Location); if err != nil {
		log.Fatal(err)
	}
	return c, err
}

func (db *postgresDB) GetClientsByName(name string) ([]*clients.Client, error) {
	rows, err := db.db.Query("SELECT id, name, postaladdress, physicaladdress, email, natureofbusiness, contactperson, contactnumber, numberofguards, location FROM clients WHERE name= $1;;", name)
	if err != nil {
		return nil, err
	}
	var clientsSlice []*clients.Client
	for rows.Next() {
		c := new(clients.Client)
		err = rows.Scan(&c.ID, &c.Name, &c.PostalAddress, &c.PhysicalAddress, &c.Email, &c.NatureOfBusiness, &c.ContactPerson, &c.ContactNumber, &c.NumberOfGuards, &c.Location)
		clientsSlice = append(clientsSlice, c)
	}
	return clientsSlice, err
}

func (db *postgresDB) GetAllClients() ([]*clients.Client, error) {
	rows, err := db.db.Query("SELECT id, name, postaladdress, physicaladdress, email, natureofbusiness, contactperson, contactnumber, numberofguards, location FROM clients;")
	if err != nil {
		return nil, err
	}
	var allClientsSlice []*clients.Client
	for rows.Next() {
		c := new(clients.Client)
		err = rows.Scan(&c.ID, &c.Name, &c.PostalAddress, &c.PhysicalAddress, &c.Email, &c.NatureOfBusiness, &c.ContactPerson, &c.ContactNumber, &c.NumberOfGuards, &c.Location)
		allClientsSlice = append(allClientsSlice, c)
	}
	return allClientsSlice, err

}

func (db *postgresDB) UpdateClient(client *clients.Client) error {
	_, err := db.db.Exec("UPDATE clients SET name = $2, postaladdress = $3, physicaladdress = $4, email =$5, natureofbusiness = $6, contactperson = $7, contactnumber =$8, numberofguards =$9, location =$10 WHERE id = $1;", client.ID, client.Name, client.PostalAddress, client.PhysicalAddress, client.Email, client.NatureOfBusiness, client.ContactPerson, client.ContactNumber, client.NumberOfGuards, client.Location)
	return err
}


func (db *postgresDB) TerminateClient(client *clients.Client) error {
	_, err := db.db.Exec("UPDATE clients SET name = $2, postaladdress = $3, physicaladdress = $4, email =$5, natureofbusiness = $6, contactperson = $7, contactnumber =$8, numberofguards =$9, location =$10 WHERE id = $1 AND deleted = 1;", client.ID, client.Name, client.PostalAddress, client.PhysicalAddress, client.Email, client.NatureOfBusiness, client.ContactPerson, client.ContactNumber, client.NumberOfGuards, client.Location)
	return err
}


func (db *postgresDB) RemoveClient(id int64) error {
	_, err := db.db.Exec("DELETE FROM clients WHERE id=$1;", id)
	return err
}

//USER DB METHODS
func (db *postgresDB) CreateUser(user *users.User) (int64, error) {
	result, err := db.db.Prepare("INSERT INTO users (name, hash) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return 0, err
	}
	defer result.Close()

	var lastInsertId int64
	err = result.QueryRow(user.Name, user.Hash).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, err
}

func (db *postgresDB) GetUserByID(id int64) (*users.User, error) {
	row := db.db.QueryRow("SELECT id, name, hash FROM users WHERE id = $1;", id)

	u := new(users.User)
	err := row.Scan(&u.ID, &u.Name, &u.Hash)
	return u, err

}

func (db *postgresDB) GetUserByName(name string) (*users.User, error) {
	row := db.db.QueryRow("SELECT id, name, hash FROM users WHERE name = $2;", name)

	u := new(users.User)
	err := row.Scan(&u.ID, &u.Name, &u.Hash)
	return u, err
}

func (db *postgresDB) GetAllUsers() ([]*users.User, error) {
	rows, err := db.db.Query("SELECT id, name, hash FROM users")
	if err != nil {
		return nil, err
	}
	var allUsersSlice []*users.User
	for rows.Next() {
		u := new(users.User)
		allUsersSlice = append(allUsersSlice, u)
	}
	return allUsersSlice, err
}

func (db *postgresDB) UpdateUser(user *users.User) error {
	_, err := db.db.Exec("UPDATE users SET name = $2, hash = $3 WHERE id = $1;", user.ID, user.Name, user.Hash)
	return err
}

func (db *postgresDB) RemoveUser(id int64) error {
	_, err := db.db.Exec("DELETE FROM users WHERE id=$1;", id)
	return err
}


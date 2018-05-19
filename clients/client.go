package clients

//Client properties
type Client struct {
	ID               int64
	Name             string
	PostalAddress	string
	PhysicalAddress	string
	Email			string
	NatureOfBusiness string
	ContactPerson	string
	ContactNumber	string
	NumberOfGuards	string
	Location		string
	deleted			int64
}

//Client interface to abstract functionality
type ClientDB interface {
	CreateClient(*Client) (int64, error)
	GetClientByID(id int64) (*Client, error)
	GetClientsByName(name string) ([]*Client, error)
	GetAllClients() ([]*Client, error)
	UpdateClient(client *Client) error
	TerminateClient(client *Client) error
	RemoveClient(id int64) error
}


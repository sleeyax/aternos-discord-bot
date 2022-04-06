package database

type Database interface {
	Connect() error
	Disconnect() error
	ListDatabases() ([]string, error) // TODO: remove this
}

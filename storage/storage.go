package storage

type Request struct {
	ID     string
	Status string
}

type Storage interface {
	FetchAndReserveRequests() ([]Request, error)
}

package store

type Store interface {
	Create(string, string) (string, error)
}

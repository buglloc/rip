package handlers

type Parser interface {
	All() ([]Handler, error)
	Next() (Handler, error)
	NextRaw() (string, error)
}

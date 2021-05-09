package handlers

type Parser interface {
	All() ([]Handler, error)
	Next() (Handler, error)
	RestValues() ([]string, error)
	NextRaw() (string, error)
	FQDN() string
}

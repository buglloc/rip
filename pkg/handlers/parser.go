package handlers

type Parser interface {
	RestHandlers() ([]Handler, error)
	NextHandler() (Handler, error)
	NextValue() (string, error)
	RestValues() ([]string, error)
	NextRaw() (string, error)
	FQDN() string
}

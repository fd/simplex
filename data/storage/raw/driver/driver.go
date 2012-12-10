package driver

type I interface {
	Ids() ([]string, error)
	Get(id string) ([]byte, error)
	Commit(set map[string][]byte, del []string) error
}

package driver

type I interface {
	Ids() ([]string, error)
	Get(id string) (interface{}, error)
	Commit(set map[string]interface{}, del []string) error
}

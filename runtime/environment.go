package runtime

type (
	Environment struct {
		tables    map[string]Table
		terminals []Terminal
	}

	Terminal interface {
		Resolve(txn *Transaction)
	}
)

func (env *Environment) Transaction() *Transaction {
	return &Transaction{
		env: env,
	}
}

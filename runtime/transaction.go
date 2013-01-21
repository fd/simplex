package runtime

type (
	Transaction struct {
		env     *Environment
		changes []*Change

		errors []interface{}
	}

	ChangeKind uint

	Change struct {
		Kind  ChangeKind
		Table string
		Key   interface{}
		Value interface{}
	}
)

const (
	SET ChangeKind = iota
	UNSET
)

func (txn *Transaction) Set(table Table, key interface{}, val interface{}) {
	change := &Change{SET, table.TableId(), key, val}
	txn.changes = append(txn.changes, change)
}

func (txn *Transaction) Unset(table Table, key interface{}) {
	change := &Change{UNSET, table.TableId(), key, nil}
	txn.changes = append(txn.changes, change)
}

func (txn *Transaction) Commit() {
	// wait for prev txn to resolve

	dones := make([]<-chan bool, len(txn.env.terminals))
	for i, t := range txn.env.terminals {
		dones[i] = ResolveTerminal(txn, t)
	}

	for _, done := range dones {
		<-done
	}
}

func ResolveTerminal(txn *Transaction, t Terminal) <-chan bool {
	done := make(chan bool)
	go func() {
		defer func() {
			done <- true
		}()
		t.Resolve(txn)
	}()
	return done
}

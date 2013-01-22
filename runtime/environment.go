package runtime

import (
	"sort"
)

var Env *Environment

type (
	Environment struct {
		tables    map[string]Table
		terminals []Terminal
	}

	Terminal interface {
		Resolve(txn *Transaction, events chan<- Event)
	}
)

func init() {
	Env = &Environment{
		tables: map[string]Table{},
	}
}

func (env *Environment) Transaction() *Transaction {
	return &Transaction{
		env: env,
	}
}

func (env *Environment) RegisterTerminal(terminal Terminal) {
	env.terminals = append(env.terminals, terminal)
}

func (env *Environment) RegisterTable(tab Table) {
	env.tables[tab.TableId()] = tab
}

func (env *Environment) Tables() []string {
	names := make([]string, 0, len(env.tables))
	for name := range env.tables {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

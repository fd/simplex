package runtime

import (
	"github.com/fd/simplex/data/storage"
	"sort"
)

var Env *Environment

type (
	Environment struct {
		tables    map[string]Table
		terminals []Terminal
		store     *storage.S
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

func (env *Environment) ConnectToStorage(url string) error {
	s, err := storage.New(url)
	if err != nil {
		return err
	}

	env.store = s
	return nil
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

package runtime

import (
	"github.com/fd/simplex/data/storage"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

var Env *Environment

type (
	Environment struct {
		tables    map[string]Table
		terminals []Terminal
		services  []Service
		store     *storage.S
	}

	Terminal interface {
		DeferredId() string
		Resolve(txn *Transaction, events chan<- Event)
	}

	Service interface {
		Start() error
		Stop() error
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

func (env *Environment) Start() error {
	for _, srv := range env.services {
		err := srv.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *Environment) Stop() error {
	for _, srv := range env.services {
		err := srv.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *Environment) Run() error {
	err := env.Start()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	_ = <-c

	err = env.Stop()
	return err
}

func (env *Environment) RegisterTerminal(terminal Terminal) {
	env.terminals = append(env.terminals, terminal)
}

func (env *Environment) RegisterTable(tab Table) {
	env.tables[tab.TableId()] = tab
}

func (env *Environment) RegisterService(srv Service) {
	env.services = append(env.services, srv)
}

func (env *Environment) Tables() []string {
	names := make([]string, 0, len(env.tables))
	for name := range env.tables {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (env *Environment) GetCurrentTransaction() (storage.SHA, bool) {
	return env.store.GetEntry()
}

func (env *Environment) SetCurrentTransaction(curr, prev storage.SHA) {
	// conditional atomic ...
	env.store.SetEntry(curr)
}

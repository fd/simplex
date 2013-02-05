package runtime

import (
	"github.com/fd/simplex/cas"
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
		store     cas.Store
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
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)
	sig := <-c
	if sig == syscall.SIGUSR1 {
		panic("dump")
	}

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

func (env *Environment) Store() cas.Store {
	return env.store
}

func (env *Environment) LoadTable(addr cas.Addr) *InternalTable {
	var kv KeyValue

	if !env.store.Get(storage.SHA(sha), &kv) {
		panic("corrupt")
	}

	table := &InternalTable{}
	if !env.store.Get(kv.ValueSha, &table) {
		panic("corrupt")
	}

	table.env = env
	table.setup()
	return table
}

func (env *Environment) GetCurrentTransaction() (cas.Addr, bool) {
	return env.store.GetEntry()
}

func (env *Environment) SetCurrentTransaction(curr, prev cas.Addr) {
	// conditional atomic ...
	env.store.SetEntry(curr)
}

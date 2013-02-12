package runtime

import (
	"os"
	"os/signal"
	go_runtime "runtime"
	"simplex.sh/cas"
	"sort"
	"syscall"
)

var Env *Environment

type (
	Environment struct {
		Store cas.Store

		tables    map[string]Table
		terminals []Terminal
		services  []Service
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
	go_runtime.GOMAXPROCS(go_runtime.NumCPU())

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

func (env *Environment) GetCurrentTransaction() (cas.Addr, error) {
	addr, err := cas.GetRef(env.Store, "_main")
	if cas.IsNotFound(err) {
		return nil, nil
	}
	return addr, err
}

func (env *Environment) SetCurrentTransaction(curr, prev cas.Addr) error {
	return cas.SetRef(env.Store, "_main", curr)
}

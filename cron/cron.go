package cron

import (
	"simplex.sh/runtime"
	"time"
)

type (
	Cron struct {
		RunAt time.Time
	}

	// view[]Cron
	cron_view interface {
		runtime.IndexedView
		EltZero() Cron
	}

	cron_service_t struct {
		views    []cron_view
		combined runtime.Deferred
	}
)

var cron_service = &cron_service_t{}

func init() {
	runtime.Env.RegisterTerminal(cron_service)
}

func (srv *cron_service_t) Register(view cron_view) {
	srv.views = append(srv.views, view)

	def_views := make([]runtime.Deferred, len(srv.views))
	for i, view := range srv.views {
		def_views[i] = view
	}

	srv.combined = runtime.Sort(
		runtime.Select(
			runtime.Union(def_views...),
			func(m interface{}) bool {
				if cron, ok := m.(Cron); ok {
					return cron.RunAt.After(time.Now())
				}
				return false
			},
		),
		func(m interface{}) interface{} {
			if cron, ok := m.(Cron); ok {
				return cron.RunAt
			}
			return nil
		},
	)
}

func (srv *cron_service_t) Resolve(txn *runtime.Transaction, out chan<- runtime.Event) {
	in := txn.Resolve(srv.combined)

	for e := range in {
		switch event := e.(type) {

		case runtime.EvError:
			out <- event

		// resolved a complete table
		case runtime.EvResolvedTable:

		}
	}
}

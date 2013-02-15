package runtime

import (
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *group_op) Resolve(state promise.State, events chan<- event.Event) {
}

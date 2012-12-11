package data

type upstream_state interface {
	Id() []string

	Ids() []string
	Get(id string) Value

	Added() []string
	Changed() []string
	Removed() []string

	NewState(segment ...string) *state
}

type state struct {
	id []string

	Info transformation_state

	added   []string
	changed []string
	removed []string
}

type transformation_state interface {
	Ids() []string
	Get(id string) Value
}

func (s *state) Id() []string {
	return s.id
}

func (s *state) Ids() []string {
	return s.Info.Ids()
}

func (s *state) Get(id string) Value {
	return s.Info.Get(id)
}

func (s *state) Added() []string {
	return s.added
}

func (s *state) Changed() []string {
	return s.changed
}

func (s *state) Removed() []string {
	return s.removed
}

func (s *state) NewState(segment ...string) *state {
	return &state{
		id: append(s.Id(), segment...),
	}
}

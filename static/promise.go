package static

type Promise struct {
	t     Transformation
	tx    *Tx
	Value interface{}
	Valid bool
}

func (p *Promise) Wait() (interface{}, error) {
	err := p.t.Wait()
	p.Valid = (err == nil)
	return p.Value, err
}

func (p *Promise) Err() error {
	return p.t.Err()
}

func (p *Promise) Do(f func() (interface{}, error)) {
	p.t.Do(func() error {
		v, err := f()
		p.Value = v
		return err
	})
}

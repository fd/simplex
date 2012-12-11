package data

type transformation interface {
	Id() string

	Chain() []transformation
	Dependencies() []transformation
	PushDownstream(transformation)

	Transform(upstream upstream_state, txn *transaction)
}

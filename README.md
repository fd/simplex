# Simplex

Simplex is a superset of Go.

## Additional builtin types

```go
view[K]M
view[]M
table[K]M
```

## Additional builtin functions (and methods)

```go
type M struct {}
type V view[string]M
type T view[string]M

V.materialize()                => T
len(V)                         => int
V.inject(func(M)A, func([]A)A) => A
V.collect(func(M)N)            => N.(view)
V.select(func(M)bool)          => V
V.reject(F)                   <=> V.select(func(m M)bool{ return !F(m) })
V.detect(F)                   <=> V.select(F)[0]
V.group(func(M)N)              => struct{ Key N; Members T }.(view)
V.slice(idx, len)              => V
V[idx]                        <=> V.materialize()[idx]
V[idx:len]                    <=> V.slice(idx, len)
V[key]                        <=> V.materialize()[key]
m, idx := range V             <=> m, idx := range V.materialize()
m := range V                  <=> m := range V.materialize()
```


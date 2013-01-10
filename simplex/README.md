# Simplex

Simplex is a superset of Go.

## Additional builtin type

```go
StructType.(view)
StructType.(table)
```

## Additional builtin functions (and methods)

```go
type M view {}
type V M.(view)
type T M.(table)

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


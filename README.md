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

V.wait()                       => V
len(V)                         => int
V.inject(func(M)A, func([]A)A) => A
V.collect(func(M)N)            => view[string]N
V.select(func(M)bool)          => V
V.reject(F)                   <=> V.select(func(m M)bool{ return !F(m) })
V.detect(F)                   <=> V.select(F)[0]
V.group(func(M)N)              => view[K]view[]M
V.slice(idx, len)              => V
V[idx]                        <=> V.wait()[idx]
V[idx:len]                    <=> V.slice(idx, len)
V[key]                        <=> V.wait()[key]
m, idx := range V             <=> m, idx := range V.wait()
m := range V                  <=> m := range V.wait()
```


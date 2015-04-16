# assert

[![GoDoc](https://godoc.org/github.com/arschles/assert?status.svg)](https://godoc.org/github.com/arschles/assert)

`assert` is [Go](http://golang.org/) package that provides convenience methods
for writing assertions in [standard Go tests](http://godoc.org/testing).

You can write this test with `assert`:

```go
func TestSomething(t *testing.T) {
  i, err := doSomething()
  assert.NoErr(err)
  assert.Equal(i, 123, "returned integer")
}
```

Instead of writing this test with only the standard `testing` library:

```go
func TestSomething(t *testing.T) {
  i, err := doSomething()
  if err != nil {
    t.Fatalf("error encountered: %s", err)
  }
  if i != 123 {
    t.Fatalf("returned integer was %d, not %d", i, 123)
  }
}
```

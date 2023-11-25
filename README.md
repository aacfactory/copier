# Copier
copy source value into destination value.

# Feature
* Sample Type
* Tag Match
* Copy sql.NullX into Value Type
* Support Slice
* Support Map

# Install
```bash
go get github.com/aacfactory/copier
```

# Example
destination
```go
type Date time.Time

type Foo struct {
	Str         string
	Int         int
	Bool        bool
	Float       float64
	Date        Date                `copy:"Time"`
	Bytes       json.RawMessage
	Faz         *Faz                `copy:"Faz"`
	Bazs        []*Faz
	Ints        []int
	ISS         [][]int
	FazMap      map[string]*Faz     `copy:"FazMap"`
	SQLTime     time.Time
	SQLString   string
	Ignore      interface{}         `copy:"-"`
}

type Faz struct {
	X string
}
```
source
```go
type Bar struct {
	Str         string
	Int         int
	Bool        bool
	Float       float64
	Time        Date                `copy:"Time"`
	Bytes       json.RawMessage
	Baz         *Baz                `copy:"Faz"`
	Bazs        []*Baz
	Ints        []int
	ISS         [][]int
	BazMap      map[string]*Baz     `copy:"FazMap"`
	SQLTime     time.Time
	SQLString   string
	Ignore      interface{}         `copy:"-"`
}

type Baz struct {
	X string
}
```
copy
```go
foo := &Foo{}
bar := Bar{...}
// use copy
err := copier.Copy(foo, bar)
// use value of
foo, err = copier.ValueOf[Foo](bar)
```
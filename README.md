# Copier
copy source value into destination value.

## Feature
* Support tag
* Support sql.Scanner
* Support anonymous
* Support slice
* Support map
* Support getter or setter
* support type convert

## Note
* DON'T SUPPORT UNEXPORTED EMBED PTR 

## Install
```bash
go get github.com/aacfactory/copier
```

## Example
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

## Type convert
* string
  * bool (false, true)
  * int
  * float
  * uint
  * byte
  * []byte
  * time.Time (RFC3339)
  * sql.NullString 
  * text.TextMarshaler 
* bool
  * string (false, true)
  * int (-1,0,1)
  * uint(0,1)
  * byte(f,t)
  * sql.NullBool
* int
  * string
  * bool (-1,1)
  * float
  * uint
  * time.Time (UnixMilli)
  * sql.NullIntX
* float
  * string 
  * int
  * uint
  * sql.NullFloatX
* uint
  * string
  * bool (0,1)
  * int
  * float
  * time.Time (UnixMilli)
* byte
  * string (first byte)
  * []byte (first byte)
  * bool (t,f)
  * int 
  * encoding.TextMarshaler (first byte)
* []byte
  * string
  * encoding.TextMarshaler
* time.Time
  * string (RFC3339)
  * sql.NullTime
  * int (UnixMilli)
  * uint (UnixMilli)
* encoding.TextUnmarshaler
  * string
  * []byte
  * encoding.TextMarshaler

## Getter 
The recv of getter method must be value, and num of results must be one.

## Setter
The recv of getter method must be ptr, and num of params must be one.  
Note: when src field type equals dst field type, setter will be discarded.
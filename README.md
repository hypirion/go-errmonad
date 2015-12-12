# go-errmonad

This is an implementation of the bind operation (`>>=`) for the equivalent of
the Either/Error type in Go.

## Installation and Usage

The import path for the package is *gopkg.in/hyPiRion/go-errmonad.v1*.

To install it, run:

```shell
go get gopkg.in/hyPiRion/go-errmonad.v1
```

Usage is pretty straightforward: Where you have a chain of `if err != nil {
return ...} `, you replace it with `errmonad.Bind` where possible:

```go
func processEncoded (bs []byte) ([]byte, error) {
    var input MyStruct
    err := encoding.Unmarshal(bs, &input)
    if err != nil {
        return nil, err
    }
    output, err := input.Process()
    if err != nil {
        return nil, err
    }
    return encoding.Marshal(output)
}
```

to

```go
func UnmarshalMyStruct(bs []byte) (ms MyStruct, err error) {
    err = encoding.Unmarshal(bs, &ms)
    return
}

var processEncoded = errmonad.Bind(
    UnmarshalMyStruct,
    (MyStruct).Process,
    encoding.Marshal,
).(func (bs []byte) ([]byte, error))
```

## License

Copyright © 2015 Jean Niklas L'orange

Distributed under the BSD 3-clause license, which is available in the file
LICENSE and at the top of the source code file.

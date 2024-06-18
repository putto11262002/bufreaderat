## bufreaderat

Package bufreaderat implements buffered io.ReaderAt. It wraps io.ReaderAt, by creating a wrapper object that also implement io.ReaderAt but provide buffering.

## Installation

```
go get github.com/putto11262002/bufreaderat
```

## Usage

```
r := bytes.NewReader([]byte("hello world"))

bufr := bufreaderat.New(r, 6)

p := make([]byte, 3)
if _, err := bufr.ReadAt(buf, 2); err != nil {
	log.Fatal(err)
}
fmt.Printf("%s\n", p)

// output: llo
```

## License

This project is licensed under the [MIT License](https://github.com/putto11262002/bufreaderat/blob/master/LICENSE).

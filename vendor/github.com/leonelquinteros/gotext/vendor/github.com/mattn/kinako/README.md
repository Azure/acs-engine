# kinako

Kinako is small VM written in Go.

![](https://raw.githubusercontent.com/mattn/kinako/master/kinako.png)

(Picture licensed under CC BY-SA 3.0 by wikipedia)

## Installation
Requires Go.
```
$ go get -u github.com/mattn/kinako
```

## Usage

Embedding the interpreter into your own program:

```Go
var env = vm.NewEnv()

env.Define("foo", 1)
val, err := env.Execute(`foo + 3`)
if err != nil {
	panic(err)
}

fmt.Println(val)
```

# License

MIT

# Author

Yasuhiro Matsumoto (a.k.a mattn)

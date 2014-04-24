#go-jsarray
This marshals Go interfaces into the 'shorter json' format I've seen Google use, in this format:

```Go
type Person struct {
	Name string
	Age uint
}

x := Person{"David", 10}
```

Becomes the string `["David",10]`. Likewise:

```Go
x := []Person{
	{"Michael", 0},
	{"", 10},
	{"Janet", 23}, 
	{"", 0},
}
```

Would produce:

```Javascript
[["Michael",],[,10],["Janet", 23],]
```

Obviously, this makes interfaces impossible to recover by anything but heuristics, but for anything else
it makes a neat serialisation format.

I'm not sure what this is actually called, if you do know, throw me a mesage.

More information is in the docs http://godoc.org/github.com/TShadwell/go-jsarray/jsarray

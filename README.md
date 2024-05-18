# Columnize

[![GoDoc](https://godoc.org/github.com/rprtr258/columnize?status.svg)](https://godoc.org/github.com/rprtr258/columnize)

Easy column-formatted output for golang. Fork of [ryanuber/columnize](https://github.com/ryanuber/columnize) but more performant and has simpler API.

Columnize is a really small Go package that makes building CLI's a little bit easier. In some CLI designs, you want to output a number similar items in a human-readable way with nicely aligned columns. However, figuring out how wide to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

```go
package main

import (
    "fmt"
    "github.com/rprtr258/columnize"
)

func main() {
    output := [][]string{
        {"Bob", "Male", "38"},
        {"Sally", "Female", "26"},
    }
    fmt.Println(columnize.Columnize(output, columnize.WithHeaders("Name", "Gender", "Age")))
}
```

As you can see, you just pass in a list of strings. And the result:

```
Name   Gender  Age
Bob    Male    38
Sally  Female  26
```

Columnize is tolerant of missing or empty fields, or even empty lines, so passing in extra lines for spacing should show up as you would expect.

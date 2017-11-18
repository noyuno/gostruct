# gostruct

[![Build Status](https://travis-ci.org/noyuno/lgo.svg?branch=master)](https://travis-ci.org/noyuno/lgo)

Golang `struct` graph visualizer

## Example

    gostruct -b1 example A | dot -Tpng -o/tmp/a.png

~~~golang
package example

type E struct {
    key int
    val []rune
    D   *D
}

type D struct {
    E *E
}

type C struct {
    E E
}

type B struct {
    D *D
}
type A struct {
    B1 B
    C2 []*C
    D6 *D
    E9 []E
    *C
}
~~~

![fig](https://raw.githubusercontent.com/noyuno/gostruct/master/example/example.png)

## Requirement

- Graphviz
- Go v1.9.2
  - github.com/sirupsen/logrus

## Usage

    -b    show builtin packages
    -builtin
        show builtin packages
    -d string
        target directory (default ".")
    -debug
        outputs debug text
    -l    show label (default true)
    -label
        show label (default true)


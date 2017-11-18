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

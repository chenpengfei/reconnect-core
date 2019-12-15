package tools

type Cal interface {
	Add(int, int)
}

func Add(a int, b int) int {
	return a + b
}

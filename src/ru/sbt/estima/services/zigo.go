package services

import (
	zg "github.com/glycerine/zygomys/repl"
)


//
// Lisp for Go interpreter
// https://github.com/glycerine/zygomys
// https://github.com/glycerine/zygomys/wiki/Go-API

func test() {
	z := zg.NewGlisp()
	err := z.LoadString("(+ 3 2)")
	expr, err := z.Run()
	println (expr)
	println (err)
}

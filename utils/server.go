package utils

import "net/http"

// shamelessly stolen from https://husobee.github.io/golang/http/middleware/2015/12/22/simple-middleware.html

type middleware func(http.HandlerFunc) http.HandlerFunc

/*
usage: http.HandleFunc("/route/", utils.BuildChain(fun0, fun1, fun2, fun3))
*/
func BuildChain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}

	return m[0](BuildChain(f, m[1:cap(m)]...))
}

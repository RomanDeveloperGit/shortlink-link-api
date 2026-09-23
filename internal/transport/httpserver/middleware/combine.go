package middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Combine(mws ...Middleware) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}

		return h
	}
}

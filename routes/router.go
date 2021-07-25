package routes

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Prevent abnormal shutdown while panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Print(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Put params in context for sharing them between handlers
func wrapHandler(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

//RouterConfig function
func RouterConfig() (router *httprouter.Router) {
	router = httprouter.New()

	indexHandlers := alice.New(recoverHandler)

	registration(router, indexHandlers)
	appointment(router, indexHandlers)

	return
}

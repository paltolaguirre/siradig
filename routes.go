package main

import "github.com/gorilla/mux"
import "net/http"

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"Healthy",
		"GET",
		"/api/siradig/healthy",
		Healthy,
	},
	Route{
		"SiradigList",
		"GET",
		"/api/siradig/siradigs",
		SiradigList,
	},
	/*Route{
		"SiradigShow",
		"GET",
		"/api/siradig/siradigs/{id}",
		SiradigShow,
	},
	Route{
		"SiradigAdd",
		"POST",
		"/api/siradig/siradigs",
		SiradigAdd,
	},
	Route{
		"SiradigUpdate",
		"PUT",
		"/api/siradig/siradigs/{id}",
		SiradigUpdate,
	},
	Route{
		"SiradigRemove",
		"DELETE",
		"/api/siradig/siradigs/{id}",
		SiradigRemove,
	},*/
}

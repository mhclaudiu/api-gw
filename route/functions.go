package route

import (
	"net/http"
)

// ServeHTTP calls f(w, r).
func (f HandlerFuncCustom) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f(w, r)
}

func (obj MuxObj) HandleFuncCustom(pattern string, handler func(http.ResponseWriter, *http.Request)) {

	if handler == nil {
		panic("http: nil handler")
	}

	//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	if obj.Config.API.MainPath == "/" {

		obj.Config.API.MainPath = ""
	}

	obj.Mux.Handle(obj.Config.API.MainPath+pattern, HandlerFuncCustom(handler))
}

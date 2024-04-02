package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/main.html", http.StripPrefix("/static", neuter(fileServer)))
	mux.Handle("/static/main.css", http.StripPrefix("/static", neuter(fileServer)))
	mux.Handle("/static/js/ts/dist/", http.StripPrefix("/static", neuter(fileServer)))
	mux.Handle("/static/posts/", http.StripPrefix("/static", neuter(fileServer)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//https://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc
//type route struct {
//	pattern *regexp.Regexp
//	handler http.Handler
//}
//
//type RegexpHandler struct {
//	routes []*route
//}
//
//func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
//	h.routes = append(h.routes, &route{pattern, handler})
//}
//
//func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
//	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
//}
//
//func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	for _, route := range h.routes {
//		if route.pattern.MatchString(r.URL.Path) {
//			route.handler.ServeHTTP(w, r)
//			return
//		}
//	}
//	// no pattern matched; send 404 response
//	http.NotFound(w, r)
//}

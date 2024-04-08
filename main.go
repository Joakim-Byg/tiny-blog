package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const ServiceConfKey = "SERVICE_CONFIG"

//const AppConfigKey = "APP_CONFIG"

func initConf() {
	serviceConfPath := os.Getenv(ServiceConfKey)
	log.Printf("Reading service config from: '%v'", serviceConfPath)
	serviceConfDir, serviceConfFile := filepath.Split(serviceConfPath)

	slicedFile := strings.Split(serviceConfFile, ".")
	if len(slicedFile) > 0 {
		slicedFile = slicedFile[:len(slicedFile)-1]
	}
	serviceConfFileName := strings.Join(slicedFile, ".")
	log.Printf("Setting config-name: %v", serviceConfFileName)
	viper.SetConfigName(serviceConfFileName)
	serviceConfDir = strings.Trim(serviceConfDir, string(filepath.Separator))
	log.Printf("Setting config path: %v", serviceConfDir)
	viper.AddConfigPath(serviceConfDir)
	viper.ReadInConfig()
	viper.WatchConfig()
	//appConfigPath := os.Getenv(AppConfigKey)
	//if appConfigPath != "" {
	//	//viper.AddConfigPath(appConfigPath)
	//}
}

// inspired by https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
// try having a look at Gin-gonik

func main() {
	if os.Getenv(ServiceConfKey) == "" {
		log.Print("The env-var 'SERVICE_CONFIG' must be defined. Which should at minimum contain configs for 'port' and 'static-serve-path'")
		os.Exit(0)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Working dir is: %v", dir)
	initConf()

	mux := http.NewServeMux()
	staticServePath := viper.GetString("static-serve-path")
	log.Printf("serving static content from '%v'", staticServePath)
	fileServer := http.FileServer(http.Dir(staticServePath))
	mux.Handle("/favicon.ico", neuter(fileServer))
	mux.Handle("/main.html", neuter(fileServer))
	mux.Handle("/main.css", neuter(fileServer))
	mux.Handle("/js/ts/dist/", neuter(fileServer))
	mux.Handle("/posts/", neuter(fileServer))
	mux.Handle("/metrics", promhttp.Handler())

	portStr := fmt.Sprintf(":%v", viper.GetString("port"))
	log.Printf("Listening on %v...", portStr)
	err = http.ListenAndServe(portStr, mux)

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

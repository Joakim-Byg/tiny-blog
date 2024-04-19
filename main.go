package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"tiny-blog/internal"
	"tiny-blog/web/app"
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
	serviceConfDir = strings.TrimRight(serviceConfDir, string(filepath.Separator))
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

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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

	var (
		otelColHost = viper.GetString("otel.collector.host")
		otelColPort = viper.GetString("otel.collector.port")
		otelAppName = viper.GetString("otel.app.name")
	)
	var otelEnabled = !(otelColHost == "" || otelColPort == "" || otelAppName == "")
	if !otelEnabled {
		log.Printf("OpenTelemetry config not correctly setup; values are: {\"otel.collector.host\" = \"%v\", "+
			"\"otel.collector.port\" = \"%v\", \"otel.app.name\" = \"%v\"}",
			otelColHost, otelColPort, otelAppName)
	} else {

		// Set up OpenTelemetry.
		otelShutdown, err := internal.SetupOTelSDK(ctx,
			viper.GetString("otel.collector.host"),
			viper.GetString("otel.collector.port"),
			viper.GetString("otel.app.name"),
		)
		if err != nil {
			return
		}
		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	handler := newHttpHandler(viper.GetString("web.static-content-path"), otelEnabled)

	portStr := fmt.Sprintf(":%v", viper.GetString("web.host.port"))
	log.Printf("Listening on %v...", portStr)
	err = http.ListenAndServe(portStr, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func newHttpHandler(staticServePath string, otelEnabled bool) http.Handler {

	router := chi.NewRouter()

	router.Use(otelchi.Middleware("tiny-blog", otelchi.WithChiRoutes(router)))
	router.Handle("/metrics", promhttp.Handler())

	log.Printf("serving static content from '%v'", staticServePath)
	fileServer := http.FileServer(http.Dir(staticServePath))
	filePaths := []string{"/favicon.ico", "/main.html", "/main.css", "/js/ts/dist/*", "/posts/*"}
	for i := 0; i < len(filePaths); i++ {
		if otelEnabled {
			router.Handle(filePaths[i], otelhttp.WithRouteTag(filePaths[i], neuter(fileServer)))
		} else {
			router.Handle(filePaths[i], neuter(fileServer))
		}
	}
	if otelEnabled {
		router.Handle("/generate-trace", otelhttp.WithRouteTag("generate-trace", http.HandlerFunc(app.GenerateTrace)))
		router.Handle("/add-trace-depth/{depth}", otelhttp.WithRouteTag("generate-trace", http.HandlerFunc(app.AddDepth)))
		return otelhttp.NewHandler(router, "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	} else {
		router.HandleFunc("/generate-trace", app.GenerateTrace)
		router.HandleFunc("/add-trace-depth/{depth}", app.AddDepth)
		return router
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

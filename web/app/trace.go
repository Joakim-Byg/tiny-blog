package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	TraceName string = "generate-trace"
	DepthVal  string = "depth.value"
)

type SimpleResponse struct {
	ResponseMessage string
}

func GenerateTrace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	depth := rand.Intn(12-3) + 3
	response := SimpleResponse{ResponseMessage: fmt.Sprintf("generated trace with depth of %v", depth)}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		labeler.Add(attribute.Bool("error", true))
		return
	}

	labeler.Add(attribute.Int(DepthVal, depth))
	trace.SpanFromContext(ctx).TracerProvider().Tracer("depthTracer").Start(ctx, "addDepth")
	get(ctx, fmt.Sprintf("add-trace-depth/%v", depth))

	w.Header().Set("Content-Type", "application/json")

	w.Write(responseJSON)
}

func AddDepth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	depthString := r.PathValue("depth")
	log.Printf("Depth is: %v", depthString)

	var depth = 0
	var err error
	if depthString != "" {
		depth, err = strconv.Atoi(depthString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			labeler.Add(attribute.Bool("error", true))
			return
		}
	}
	labeler.Add(attribute.Int(DepthVal, depth))
	var responseJSON []byte
	responseJSON, err = json.Marshal(SimpleResponse{ResponseMessage: "hello!"})
	if depth > 0 {
		responseJSON, err = json.Marshal(SimpleResponse{ResponseMessage: fmt.Sprintf("depth was: %v", depthString)})
		trace.SpanFromContext(ctx).TracerProvider().Tracer("depthTracer").Start(ctx, "addDepth")
		get(ctx, fmt.Sprintf("add-trace-depth/%v", depth-1))

	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		labeler.Add(attribute.Bool("error", true))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(responseJSON)

}

func get(ctx context.Context, path string) {
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resourcePath := fmt.Sprintf("http://localhost:%v/%v", viper.GetString("web.host.port"), path)

	req, err := http.NewRequestWithContext(ctx, "GET", resourcePath, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject SimpleResponse
	json.Unmarshal(bodyBytes, &responseObject)
	log.Printf("API Response as struct %+v\n", responseObject)
}

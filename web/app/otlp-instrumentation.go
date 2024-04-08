package main

//func initProvider() func() {
//	ctx := context.Background()
//
//	otelAgentAddr, xtraceToken, ok := common.ObtainXTraceInfo()
//
//	if !ok {
//		log.Fatalf("Cannot init OpenTelemetry, exit")
//		os.Exit(-1)
//	}
//
//	headers := map[string]string{"Authentication": xtraceToken} // Replace xtraceToken with the authentication token obtained in the Prerequisites section.
//	traceClient := otlptracegrpc.NewClient(
//		otlptracegrpc.WithInsecure(),
//		otlptracegrpc.WithEndpoint(otelAgentAddr), // Replace otelAgentAddr with the endpoint obtained in the Prerequisites section.
//		otlptracegrpc.WithHeaders(headers),
//		otlptracegrpc.WithDialOption(grpc.WithBlock()))
//	log.Println("start to connect to server")
//	traceExp, err := otlptrace.New(ctx, traceClient)
//	handleErr(err, "Failed to create the collector trace exporter")
//
//	res, err := resource.New(ctx,
//		resource.WithFromEnv(),
//		resource.WithProcess(),
//		resource.WithTelemetrySDK(),
//		resource.WithHost(),
//		resource.WithAttributes(
//			// Specify the service name displayed on the backend of Managed Service for OpenTelemetry.
//			semconv.ServiceNameKey.String(common.ServerServiceName),
//			semconv.HostNameKey.String(common.ServerServiceHostName),
//		),
//	)
//	handleErr(err, "failed to create resource")
//
//	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
//	tracerProvider := sdktrace.NewTracerProvider(
//		sdktrace.WithSampler(sdktrace.AlwaysSample()),
//		sdktrace.WithResource(res),
//		sdktrace.WithSpanProcessor(bsp),
//	)
//
//	// Set the global propagator to tracecontext. The global propagator is not specified by default.
//	otel.SetTextMapPropagator(propagation.TraceContext{})
//	otel.SetTracerProvider(tracerProvider)
//
//	return func() {
//		cxt, cancel := context.WithTimeout(ctx, time.Second)
//		defer cancel()
//		if err := traceExp.Shutdown(cxt); err != nil {
//			otel.Handle(err)
//		}
//	}
//}

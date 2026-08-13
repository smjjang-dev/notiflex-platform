package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var valkeyClient valkey.Client
var kafkaProducer sarama.SyncProducer

var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by path and response status code",
	},
	[]string{"path", "code"},
)

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

// instrument wraps a handler so every request is counted in http_requests_total,
// which the canary AnalysisTemplate queries to gate rollout promotion.
func instrument(path string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
		h(rec, r)
		httpRequestsTotal.WithLabelValues(path, strconv.Itoa(rec.code)).Inc()
	}
}

// set via -ldflags at build time (see Dockerfile)
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func initTracer(endpoint string) func() {
	if endpoint == "" {
		return func() {}
	}
	ctx := context.Background()
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("OTel gRPC 연결 실패: %v", err)
		return func() {}
	}
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Printf("OTel exporter 생성 실패: %v", err)
		return func() {}
	}
	// resource.Merge keeps the first argument's value on key conflicts, and
	// resource.Default() already sets its own (unknown_service:...) service.name
	res, err := resource.Merge(
		resource.NewSchemaless(attribute.String("service.name", "notiflex-api")),
		resource.Default(),
	)
	if err != nil {
		log.Printf("OTel resource 생성 실패: %v", err)
		res = resource.Default()
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return func() { tp.Shutdown(ctx) }
}

func main() {
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("VALKEY_PASSWORD")
	if pwFile := os.Getenv("VALKEY_PASSWORD_FILE"); pwFile != "" {
		if data, err := os.ReadFile(pwFile); err == nil {
			password = string(data)
		}
	}

	var err error
	for i := 0; i < 10; i++ {
		valkeyClient, err = valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{addr},
			Password:    password,
		})
		if err == nil {
			break
		}
		log.Printf("Valkey 연결 재시도 %d/10: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Valkey 연결 실패: %v", err)
	}
	defer valkeyClient.Close()

	broker := os.Getenv("KAFKA_BROKER")
	if broker != "" {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Version = sarama.V4_1_0_0
		kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, cfg)
		if err != nil {
			log.Printf("Kafka 연결 실패 (계속): %v", err)
		} else {
			defer kafkaProducer.Close()
			go consumeKafka(broker)
		}
	}

	shutdown := initTracer(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	defer shutdown()

	http.HandleFunc("/health", instrument("/health", healthHandler))
	http.HandleFunc("/id", instrument("/id", idHandler))
	http.HandleFunc("/version", instrument("/version", versionHandler))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Notiflex API server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func consumeKafka(broker string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	consumer, err := sarama.NewConsumer([]string{broker}, cfg)
	if err != nil {
		log.Printf("Kafka consumer 생성 실패: %v", err)
		return
	}
	defer consumer.Close()
	pc, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Kafka partition consumer 생성 실패: %v", err)
		return
	}
	defer pc.Close()
	for msg := range pc.Messages() {
		log.Printf("[Kafka] 수신: key=%s value=%s", string(msg.Key), string(msg.Value))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("notiflex")
	_, span := tracer.Start(r.Context(), "health")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version":   version,
		"commit":    commit,
		"buildTime": buildTime,
		"goVersion": runtime.Version(),
	})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("notiflex")
	ctx, span := tracer.Start(r.Context(), "id")
	defer span.End()

	result, err := valkeyClient.Do(ctx, valkeyClient.B().Incr().Key("notiflex:id").Build()).AsInt64()
	if err != nil {
		http.Error(w, "Valkey error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		pod = "unknown"
	}

	if kafkaProducer != nil {
		msg := &sarama.ProducerMessage{
			Topic: "notifications",
			Key:   sarama.StringEncoder(fmt.Sprintf("id-%d", result)),
			Value: sarama.StringEncoder(fmt.Sprintf(`{"id":%d,"pod":"%s"}`, result, pod)),
		}
		if _, _, err := kafkaProducer.SendMessage(msg); err != nil {
			log.Printf("[Kafka] 전송 실패: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": result, "pod": pod})
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/queue"
	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/streamer"
)

func main() {
	log.Println("[streamer] starting telemetry streamer")

	// --- Configuration via environment variables ---
	csvPath := mustGetEnv("CSV_FILE_PATH")
	streamInterval := getDurationEnv("STREAM_INTERVAL_MS", 1000)
	loop := getBoolEnv("LOOP_CSV", true)

	// --- Create context for graceful shutdown ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Handle SIGTERM / SIGINT ---
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Printf("[streamer] received signal: %v, shutting down\n", sig)
		cancel()
	}()

	// --- Initialize Queue ---
	// NOTE: This is a library reference; actual queue instance
	// must be shared with collectors in deployment
	q := queue.GetQueue() // singleton or injected reference

	// --- Initialize CSV Streamer ---
	cfg := streamer.CSVReaderConfig{
		FilePath:       csvPath,
		StreamInterval: streamInterval,
		Loop:           loop,
	}

	csvStreamer, err := streamer.NewCSVStreamer(cfg, q)
	if err != nil {
		log.Fatalf("[streamer] failed to create CSV streamer: %v", err)
	}

	// --- Run streamer ---
	if err := csvStreamer.Run(ctx); err != nil {
		log.Printf("[streamer] exited with error: %v", err)
	}

	log.Println("[streamer] stopped cleanly")
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func getDurationEnv(key string, defaultMs int) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return time.Duration(defaultMs) * time.Millisecond
	}

	ms, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid duration for %s: %v", key, err)
	}
	return time.Duration(ms) * time.Millisecond
}

func getBoolEnv(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("invalid boolean for %s: %v", key, err)
	}
	return b
}

package streamer

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/common"
	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/queue"
)

type CSVReaderConfig struct {
	FilePath       string        // Path to telemetry CSV
	StreamInterval time.Duration // Delay between telemetry events
	Loop           bool          // Loop CSV indefinitely
}

// CSVStreamer reads telemetry from CSV and streams it into the queue
type CSVStreamer struct {
	cfg   CSVReaderConfig
	queue queue.Queue
}

// NewCSVStreamer creates a new CSVStreamer
func NewCSVStreamer(cfg CSVReaderConfig, q queue.Queue) (*CSVStreamer, error) {
	if cfg.FilePath == "" {
		return nil, errors.New("csv file path cannot be empty")
	}
	if q == nil {
		return nil, errors.New("queue cannot be nil")
	}

	return &CSVStreamer{
		cfg:   cfg,
		queue: q,
	}, nil
}

// Run starts streaming telemetry until context is cancelled
func (s *CSVStreamer) Run(ctx context.Context) error {
	for {
		err := s.streamOnce(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		if !s.cfg.Loop {
			return nil
		}
	}
}

// streamOnce streams the CSV file exactly once
func (s *CSVStreamer) streamOnce(ctx context.Context) error {
	file, err := os.Open(s.cfg.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Optional: skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		telemetry, err := parseTelemetry(record)
		if err != nil {
			// Skip malformed rows, do not crash streamer
			continue
		}

		telemetry.Timestamp = time.Now().UTC()

		s.queue.Enqueue(telemetry)

		if s.cfg.StreamInterval > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(s.cfg.StreamInterval):
			}
		}
	}
}

// parseTelemetry converts a CSV row into a Telemetry struct
// Expected CSV format (example):
// host_id,gpu_id,utilization,memory_used,temperature
func parseTelemetry(record []string) (common.Telemetry, error) {
	if len(record) < 5 {
		return common.Telemetry{}, errors.New("invalid telemetry record")
	}

	utilization, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return common.Telemetry{}, err
	}

	memoryUsed, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return common.Telemetry{}, err
	}

	temperature, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return common.Telemetry{}, err
	}

	return common.Telemetry{
		ID:          uuid.NewString(),
		HostID:      record[0],
		GPUId:       record[1],
		Utilization: utilization,
		MemoryUsed:  memoryUsed,
		Temperature: temperature,
		// Timestamp set at processing time
	}, nil
}

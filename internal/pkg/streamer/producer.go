package streamer

import (
	"context"
	"errors"
	"time"

	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/common"
	"github.com/mythilirajendra-new/gpu-telemetry-service/internal/pkg/queue"
)

// ProducerConfig controls producer behavior
type ProducerConfig struct {
	// Optional max rate (telemetry per second)
	MaxRate int

	// Optional batch size (0 or 1 = single message)
	BatchSize int
}

// Producer is responsible for sending telemetry to the queue
type Producer struct {
	queue queue.Queue
	cfg   ProducerConfig
}

// NewProducer creates a new Producer
func NewProducer(q queue.Queue, cfg ProducerConfig) (*Producer, error) {
	if q == nil {
		return nil, errors.New("queue cannot be nil")
	}

	if cfg.BatchSize < 0 {
		return nil, errors.New("batch size cannot be negative")
	}

	return &Producer{
		queue: q,
		cfg:   cfg,
	}, nil
}

// Produce sends telemetry into the queue
func (p *Producer) Produce(ctx context.Context, telemetry common.Telemetry) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Timestamp must reflect processing time
	if telemetry.Timestamp.IsZero() {
		telemetry.Timestamp = time.Now().UTC()
	}

	p.queue.Enqueue(telemetry)
	return nil
}

package queue

import "sync"

type Telemetry struct {
	GPUId       string `json:"gpu_id"`
	ModelName   string `json:"model_name"`
	Hostname    string `json:"hostname"`
	ContainerId string `json:"container_id"`
	Device      string `json:"device"`
	MetricName  string `json:"metric_name"`
	Timestamp   int64  `json:"timestamp"`
}

type Queue interface {
	Enqueue(msg Telemetry)
	Dequeue(batch int) []Telemetry
}

type InMemoryQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []Telemetry
}

func (q *InMemoryQueue) Enqueue(t Telemetry) {
	q.mu.Lock()
	q.queue = append(q.queue, t)
	q.mu.Unlock()
	q.cond.Signal()
}

func (q *InMemoryQueue) Dequeue(n int) []Telemetry {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.queue) == 0 {
		q.cond.Wait()
	}

	batch := min(n, len(q.queue))
	msgs := q.queue[:batch]
	q.queue = q.queue[batch:]
	return msgs
}

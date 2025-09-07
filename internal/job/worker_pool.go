package job

import "indico-be/internal/service"

// WorkerPool simply carries configuration.
type WorkerPool struct {
	Count   int
	Service *service.SettlementService
}

func NewWorkerPool(count int, svc *service.SettlementService) *WorkerPool {
	return &WorkerPool{Count: count, Service: svc}
}

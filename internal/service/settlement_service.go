package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"indico-be/internal/models"
	"indico-be/internal/repository"
)

type SettlementService struct {
	txRepo    repository.TransactionRepository
	setRepo   repository.SettlementRepository
	JobRepo   repository.JobRepository
	mu        sync.Mutex // protects updates to the same job record
	batchSize int
}

// NewSettlementService injects required repos.
func NewSettlementService(tx repository.TransactionRepository,
	set repository.SettlementRepository,
	job repository.JobRepository) *SettlementService {

	return &SettlementService{
		txRepo:    tx,
		setRepo:   set,
		JobRepo:   job,
		batchSize: 5000,
	}
}

// RunJob is called by a worker for a specific job ID.
func (s *SettlementService) RunJob(ctx context.Context, jobID, fromStr, toStr string) error {
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return fmt.Errorf("invalid from date: %w", err)
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return fmt.Errorf("invalid to date: %w", err)
	}

	total, err := s.txRepo.CountByPeriod(context.Background(), from, to)
	if err != nil {
		return fmt.Errorf("failed counting transactions: %w", err)
	}
	if err := s.JobRepo.UpdateTotal(context.Background(), jobID, total); err != nil {
		return fmt.Errorf("failed updating total: %w", err)
	}

	var offset int64 = 0
	for {
		batch, err := s.txRepo.GetBatch(context.Background(), from, to, int(offset), s.batchSize)
		if err != nil {
			return fmt.Errorf("failed fetching batch: %w", err)
		}
		if len(batch) == 0 {
			break
		}

		// Insert/update ke tabel settlement
		settlements := make([]*models.Settlement, 0, len(batch))
		for _, tx := range batch {
			settlements = append(settlements, &models.Settlement{
				MerchantID:  tx.MerchantID,            
				Date:        tx.PaidAt,                
				GrossCents:  int64(tx.FeeCents * 100), 
				FeeCents:    int64(tx.FeeCents),       
				NetCents:    int64((tx.AmountCents - tx.FeeCents) * 100),
				TxnCount:    1,
				GeneratedAt: time.Now(),
				RunID:       jobID, 
			})
		}

		for _, d := range settlements {
			if err := s.setRepo.Upsert(ctx, d); err != nil {
				return fmt.Errorf("failed upserting settlement (merchant_id=%d, date=%v): %w", d.MerchantID, d.Date, err)
			}
		}

		processedBatch := int64(len(batch))

		// Update job progress & processed count
		if err := s.JobRepo.IncrementProcessed(
			context.Background(),
			jobID,
			processedBatch,
			float64(processedBatch)/float64(total)*100,
		); err != nil {
			return fmt.Errorf("failed updating progress: %w", err)
		}

		log.Printf("[Job %s] batch offset %d → processed %d (total %d)", jobID, offset, processedBatch, total)

		offset += processedBatch
	}

	return nil
}

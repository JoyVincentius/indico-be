package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"indico-be/internal/models"
	"indico-be/internal/repository"
)

type SettlementService struct {
	txRepo    repository.TransactionRepository
	setRepo   repository.SettlementRepository
	JobRepo   repository.JobRepository
	mu        sync.Mutex
	batchSize int
}

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

	var allSettlements []*models.Settlement

	var offset int64 = 0
	for {
		batch, err := s.txRepo.GetBatch(context.Background(), from, to, int(offset), s.batchSize)
		if err != nil {
			return fmt.Errorf("failed fetching batch: %w", err)
		}
		if len(batch) == 0 {
			break
		}

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

		allSettlements = append(allSettlements, settlements...)

		processedBatch := int64(len(batch))
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

	if err := s.generateCSV(ctx, jobID, allSettlements); err != nil {
		return fmt.Errorf("failed generating CSV: %w", err)
	}

	if err := s.JobRepo.UpdateStatus(ctx, jobID, "FINISHED"); err != nil {
		return fmt.Errorf("failed updating job status to FINISHED: %w", err)
	}

	log.Printf("[Job %s] COMPLETED successfully: %d settlements written to CSV", jobID, len(allSettlements))
	return nil
}

func (s *SettlementService) generateCSV(ctx context.Context, jobID string, settlements []*models.Settlement) error {
	filePath := filepath.Join("public/downloads", jobID+".csv")

	if err := os.MkdirAll("public/downloads", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create public/downloads directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"merchant_id", "date", "gross_cents", "fee_cents", "net_cents", "txn_count", "generated_at", "run_id"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, s := range settlements {
		row := []string{
			fmt.Sprintf("%d", s.MerchantID),
			s.Date.Format("2006-01-02"),
			fmt.Sprintf("%d", s.GrossCents),
			fmt.Sprintf("%d", s.FeeCents),
			fmt.Sprintf("%d", s.NetCents),
			fmt.Sprintf("%d", s.TxnCount),
			s.GeneratedAt.Format(time.RFC3339),
			s.RunID,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	log.Printf("CSV file generated: %s (%d records)", filePath, len(settlements))
	return nil
}

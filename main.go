package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"indico-be/config"
	"indico-be/internal/handler"
	"indico-be/internal/job"
	"indico-be/internal/repository"
	"indico-be/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ---------- Load config ----------
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] .env file not found or could not be loaded: %v", err)
	} else {
		log.Println("[INFO] .env file loaded successfully")
	}

	cfg := config.Load()

	// ---------- Koneksi DB ----------
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	// Auto-migrate model
	if err := db.AutoMigrate(
		&repository.Order{},
		&repository.Transaction{},
		&repository.Settlement{},
		&repository.JobRecord{},
	); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	// ---------- 3️⃣ Repositories ----------
	orderRepo := repository.NewOrderRepo(db)
	txRepo := repository.NewTransactionRepo(db)
	settleRepo := repository.NewSettlementRepo(db)
	jobRepo := repository.NewJobRepository(db)

	// ---------- 4️⃣ Services ----------
	orderSvc := service.NewOrderService(orderRepo)
	settleSvc := service.NewSettlementService(txRepo, settleRepo, jobRepo)

	// ---------- 5️⃣ Job System ----------
	workerPool := job.NewWorkerPool(cfg.WorkerCount, settleSvc)
	jobQueue := job.NewJobQueue(workerPool)
	jobQueue.SetRepository(jobRepo)
	
	log.Printf("🔧 Workers loaded: %d ", cfg.WorkerCount)

	// ---------- 6️⃣ HTTP Router ----------
	router := gin.Default()
	handler.RegisterOrderRoutes(router, orderSvc)
	handler.RegisterJobRoutes(router, jobQueue, jobRepo)

	// ---------- 7️⃣ Server & Shutdown ----------
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)

		jobQueue.Close()
	}()

	log.Printf("🚀 Server listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %v", err)
	}
}

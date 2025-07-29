package main

import (
	"log"

	worker "github.com/Grafiters/archive/cmd/workers/cron"
	"github.com/Grafiters/archive/configs"
	customerMysql "github.com/Grafiters/archive/internal/customer/repository"
	transactionMysql "github.com/Grafiters/archive/internal/transaction/repository"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load config dan dependency
	if err := configs.Initialize(); err != nil {
		log.Fatal(err)
		return
	}

	customerRepo := customerMysql.NewCustomerRepository(configs.DataBase, configs.Logger)
	transactionRepo := transactionMysql.NewTranscationRepository(configs.DataBase, configs.Logger)

	// Inisialisasi job
	reminderJob := worker.NewInstallmentReminderJob(transactionRepo, customerRepo)

	// Jalankan cron
	c := cron.New()
	_, err := c.AddFunc("0 8 * * *", func() {
		reminderJob.Run()
	})
	if err != nil {
		log.Fatalf("failed to schedule job: %v", err)
	}

	log.Println("Worker cronjob berjalan... (CTRL+C untuk keluar)")
	c.Run()
}

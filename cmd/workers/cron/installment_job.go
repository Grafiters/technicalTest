package worker

import (
	"fmt"
	"time"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
)

type InstallmentReminderJob struct {
	TransactionRepo domain.TransactionRepository
	CustomerRepo    domain.CustomerRepository
}

func NewInstallmentReminderJob(tr domain.TransactionRepository, cr domain.CustomerRepository) *InstallmentReminderJob {
	return &InstallmentReminderJob{
		TransactionRepo: tr,
		CustomerRepo:    cr,
	}
}

func (j *InstallmentReminderJob) Run() {
	today := time.Now().Format("2006-01-02")
	installments, err := j.TransactionRepo.GetInstallmentsByDueDate(today)
	if err != nil {
		fmt.Printf("[ReminderJob] failed to fetch installments: %v\n", err)
		return
	}

	for _, i := range installments {
		transaction, err := j.TransactionRepo.GetByID(i.TransactionID)
		if err != nil {
			fmt.Printf("[ReminderJob] transaction fetch failed: %v\n", err)
			continue
		}

		customer, err := j.CustomerRepo.GetByID(transaction.CustomerID)
		if err != nil {
			fmt.Printf("[ReminderJob] customer fetch failed: %v\n", err)
			continue
		}

		subject := fmt.Sprintf("Tagihan Cicilan Bulan ke-%d", i.Month)
		body := fmt.Sprintf("Halo %s,\n\nIni adalah pengingat bahwa cicilan Anda sebesar Rp%d jatuh tempo hari ini (%s).\n\nTerima kasih.\n", customer.FullName, i.Amount, today)

		configs.Logger.Info(subject)
		configs.Logger.Info(body)
	}
}

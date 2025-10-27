package services

import (
        "database/sql"
        "fmt"
        "time"

        "github.com/wiederin/go-invoicer-app/internal/models"
)

type UsageService struct {
        db *sql.DB
}

func NewUsageService(db *sql.DB) *UsageService {
        return &UsageService{db: db}
}

func (s *UsageService) GetCurrentUsage(userID int) (*models.UsageStatus, error) {
        now := time.Now()
        periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

        var usage models.UsageRecord
        var plan models.Plan

        err := s.db.QueryRow(`
                SELECT COALESCE(ur.invoices_generated, 0) as invoices_generated, p.name, p.monthly_quota
                FROM users u
                JOIN plans p ON u.plan_id = p.id
                LEFT JOIN usage_records ur ON ur.user_id = u.id 
                        AND ur.period_start = $2
                WHERE u.id = $1
        `, userID, periodStart).Scan(&usage.InvoicesGenerated, &plan.Name, &plan.MonthlyQuota)

        if err != nil {
                return nil, fmt.Errorf("failed to get usage: %w", err)
        }

        remaining := plan.MonthlyQuota - usage.InvoicesGenerated
        if plan.MonthlyQuota == -1 {
                remaining = -1
        } else if remaining < 0 {
                remaining = 0
        }

        return &models.UsageStatus{
                CurrentUsage: usage.InvoicesGenerated,
                MonthlyQuota: plan.MonthlyQuota,
                Remaining:    remaining,
                PlanName:     plan.Name,
        }, nil
}

func (s *UsageService) CanGenerateInvoice(userID int) (bool, error) {
        status, err := s.GetCurrentUsage(userID)
        if err != nil {
                return false, err
        }

        if status.MonthlyQuota == -1 {
                return true, nil
        }

        return status.CurrentUsage < status.MonthlyQuota, nil
}

func (s *UsageService) IncrementUsage(userID int) error {
        now := time.Now()
        periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
        periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

        tx, err := s.db.Begin()
        if err != nil {
                return fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        var currentUsage int
        var monthlyQuota int

        err = tx.QueryRow(`
                SELECT COALESCE(ur.invoices_generated, 0), p.monthly_quota
                FROM users u
                JOIN plans p ON u.plan_id = p.id
                LEFT JOIN usage_records ur ON ur.user_id = u.id AND ur.period_start = $2
                WHERE u.id = $1
                FOR UPDATE
        `, userID, periodStart).Scan(&currentUsage, &monthlyQuota)

        if err != nil {
                return fmt.Errorf("failed to get current usage: %w", err)
        }

        if monthlyQuota != -1 && currentUsage >= monthlyQuota {
                return fmt.Errorf("quota exceeded")
        }

        _, err = tx.Exec(`
                INSERT INTO usage_records (user_id, period_start, period_end, invoices_generated)
                VALUES ($1, $2, $3, 1)
                ON CONFLICT (user_id, period_start)
                DO UPDATE SET 
                        invoices_generated = usage_records.invoices_generated + 1,
                        updated_at = CURRENT_TIMESTAMP
        `, userID, periodStart, periodEnd)

        if err != nil {
                return fmt.Errorf("failed to increment usage: %w", err)
        }

        if err := tx.Commit(); err != nil {
                return fmt.Errorf("failed to commit transaction: %w", err)
        }

        return nil
}

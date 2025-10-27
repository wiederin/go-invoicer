package services

import (
	"database/sql"
	"fmt"

	"github.com/wiederin/go-invoicer-app/internal/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetOrCreateUser(replitUserID, email string) (*models.User, error) {
	var user models.User

	err := s.db.QueryRow(`
		SELECT id, email, replit_user_id, plan_id, status, stripe_customer_id, created_at, updated_at
		FROM users
		WHERE replit_user_id = $1
	`, replitUserID).Scan(
		&user.ID, &user.Email, &user.ReplitUserID, &user.PlanID,
		&user.Status, &user.StripeCustomerID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		err = s.db.QueryRow(`
			INSERT INTO users (email, replit_user_id, plan_id, status)
			VALUES ($1, $2, 1, 'active')
			RETURNING id, email, replit_user_id, plan_id, status, stripe_customer_id, created_at, updated_at
		`, email, replitUserID).Scan(
			&user.ID, &user.Email, &user.ReplitUserID, &user.PlanID,
			&user.Status, &user.StripeCustomerID, &user.CreatedAt, &user.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		return &user, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByID(userID int) (*models.User, error) {
	var user models.User

	err := s.db.QueryRow(`
		SELECT id, email, replit_user_id, plan_id, status, stripe_customer_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Email, &user.ReplitUserID, &user.PlanID,
		&user.Status, &user.StripeCustomerID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUserPlan(userID, planID int) error {
	_, err := s.db.Exec(`
		UPDATE users
		SET plan_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, planID, userID)

	if err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}

	return nil
}

func (s *UserService) GetAllPlans() ([]models.Plan, error) {
	rows, err := s.db.Query(`
		SELECT id, name, monthly_quota, price_cents, stripe_price_id, created_at
		FROM plans
		ORDER BY price_cents ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var plan models.Plan
		err := rows.Scan(&plan.ID, &plan.Name, &plan.MonthlyQuota, &plan.PriceCents, &plan.StripePriceID, &plan.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

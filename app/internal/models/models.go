package models

import (
        "time"
)

type Plan struct {
        ID            int       `json:"id"`
        Name          string    `json:"name"`
        MonthlyQuota  int       `json:"monthly_quota"`
        PriceCents    int       `json:"price_cents"`
        StripePriceID *string   `json:"stripe_price_id,omitempty"`
        CreatedAt     time.Time `json:"created_at"`
}

type User struct {
        ID               int       `json:"id"`
        Email            string    `json:"email"`
        ReplitUserID     *string   `json:"replit_user_id,omitempty"`
        PlanID           int       `json:"plan_id"`
        Status           string    `json:"status"`
        StripeCustomerID *string   `json:"stripe_customer_id,omitempty"`
        CreatedAt        time.Time `json:"created_at"`
        UpdatedAt        time.Time `json:"updated_at"`
}

type Subscription struct {
        ID                   int        `json:"id"`
        UserID               int        `json:"user_id"`
        PlanID               int        `json:"plan_id"`
        Status               string     `json:"status"`
        StripeSubscriptionID *string    `json:"stripe_subscription_id,omitempty"`
        CurrentPeriodStart   *time.Time `json:"current_period_start,omitempty"`
        CurrentPeriodEnd     *time.Time `json:"current_period_end,omitempty"`
        QuotaOverride        *int       `json:"quota_override,omitempty"`
        CreatedAt            time.Time  `json:"created_at"`
        UpdatedAt            time.Time  `json:"updated_at"`
}

type UsageRecord struct {
        ID                 int       `json:"id"`
        UserID             int       `json:"user_id"`
        PeriodStart        time.Time `json:"period_start"`
        PeriodEnd          time.Time `json:"period_end"`
        InvoicesGenerated  int       `json:"invoices_generated"`
        CreatedAt          time.Time `json:"created_at"`
        UpdatedAt          time.Time `json:"updated_at"`
}

type Invoice struct {
        ID            int                    `json:"id"`
        UserID        int                    `json:"user_id"`
        InvoiceNumber string                 `json:"invoice_number"`
        CompanyName   string                 `json:"company_name"`
        CustomerName  string                 `json:"customer_name"`
        TotalAmount   float64                `json:"total_amount"`
        Metadata      map[string]interface{} `json:"metadata"`
        CreatedAt     time.Time              `json:"created_at"`
}

type UsageStatus struct {
        CurrentUsage int `json:"current_usage"`
        MonthlyQuota int `json:"monthly_quota"`
        Remaining    int `json:"remaining"`
        PlanName     string `json:"plan_name"`
}

package handlers

import (
        "database/sql"
        "encoding/json"
        "fmt"
        "net/http"

        "github.com/wiederin/go-invoicer-app/internal/middleware"
        "github.com/wiederin/go-invoicer-app/internal/services"
)

type Handler struct {
        userService    *services.UserService
        usageService   *services.UsageService
        invoiceService *services.InvoiceService
}

func NewHandler(db *sql.DB) *Handler {
        return &Handler{
                userService:    services.NewUserService(db),
                usageService:   services.NewUsageService(db),
                invoiceService: services.NewInvoiceService(db),
        }
}

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
        replitUserID := r.Context().Value(middleware.UserIDKey).(string)
        userEmail := r.Context().Value(middleware.UserEmailKey).(string)

        user, err := h.userService.GetOrCreateUser(replitUserID, userEmail)
        if err != nil {
                http.Error(w, "Failed to get user", http.StatusInternalServerError)
                return
        }

        usage, err := h.usageService.GetCurrentUsage(user.ID)
        if err != nil {
                http.Error(w, "Failed to get usage", http.StatusInternalServerError)
                return
        }

        response := map[string]interface{}{
                "user":  user,
                "usage": usage,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleGenerateInvoice(w http.ResponseWriter, r *http.Request) {
        replitUserID := r.Context().Value(middleware.UserIDKey).(string)
        userEmail := r.Context().Value(middleware.UserEmailKey).(string)

        user, err := h.userService.GetOrCreateUser(replitUserID, userEmail)
        if err != nil {
                http.Error(w, "Failed to get user", http.StatusInternalServerError)
                return
        }

        var req services.InvoiceRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                http.Error(w, "Invalid request", http.StatusBadRequest)
                return
        }

        pdfBytes, err := h.invoiceService.GenerateInvoice(user.ID, &req)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to generate invoice: %v", err), http.StatusInternalServerError)
                return
        }

        if err := h.usageService.IncrementUsage(user.ID); err != nil {
                if err.Error() == "quota exceeded" {
                        w.Header().Set("Content-Type", "application/json")
                        w.WriteHeader(http.StatusPaymentRequired)
                        json.NewEncoder(w).Encode(map[string]string{
                                "error": "Monthly quota exceeded. Please upgrade your plan.",
                        })
                        return
                }
                http.Error(w, "Failed to update usage", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/pdf")
        w.Header().Set("Content-Disposition", "attachment; filename=invoice.pdf")
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Write(pdfBytes)
}

func (h *Handler) HandleGetInvoices(w http.ResponseWriter, r *http.Request) {
        replitUserID := r.Context().Value(middleware.UserIDKey).(string)
        userEmail := r.Context().Value(middleware.UserEmailKey).(string)

        user, err := h.userService.GetOrCreateUser(replitUserID, userEmail)
        if err != nil {
                http.Error(w, "Failed to get user", http.StatusInternalServerError)
                return
        }

        invoices, err := h.invoiceService.GetUserInvoices(user.ID, 50)
        if err != nil {
                http.Error(w, "Failed to get invoices", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(invoices)
}

func (h *Handler) HandleGetPlans(w http.ResponseWriter, r *http.Request) {
        plans, err := h.userService.GetAllPlans()
        if err != nil {
                http.Error(w, "Failed to get plans", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(plans)
}

func (h *Handler) HandleGetUsage(w http.ResponseWriter, r *http.Request) {
        replitUserID := r.Context().Value(middleware.UserIDKey).(string)
        userEmail := r.Context().Value(middleware.UserEmailKey).(string)

        user, err := h.userService.GetOrCreateUser(replitUserID, userEmail)
        if err != nil {
                http.Error(w, "Failed to get user", http.StatusInternalServerError)
                return
        }

        usage, err := h.usageService.GetCurrentUsage(user.ID)
        if err != nil {
                http.Error(w, "Failed to get usage", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(usage)
}

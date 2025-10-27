package main

import (
        "fmt"
        "log"
        "net/http"

        "github.com/gorilla/mux"
        "github.com/wiederin/go-invoicer-app/internal/database"
        "github.com/wiederin/go-invoicer-app/internal/handlers"
        "github.com/wiederin/go-invoicer-app/internal/middleware"
)

func main() {
        db, err := database.Connect()
        if err != nil {
                log.Fatalf("Failed to connect to database: %v", err)
        }
        defer db.Close()

        if err := db.RunMigrations(); err != nil {
                log.Printf("Warning: Failed to run migrations: %v", err)
        }

        h := handlers.NewHandler(db.DB)

        r := mux.NewRouter()

        r.HandleFunc("/", serveHome).Methods("GET")
        r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

        api := r.PathPrefix("/api").Subrouter()
        api.Use(middleware.AuthMiddleware)

        api.HandleFunc("/user", h.HandleGetUser).Methods("GET")
        api.HandleFunc("/usage", h.HandleGetUsage).Methods("GET")
        api.HandleFunc("/plans", h.HandleGetPlans).Methods("GET")
        api.HandleFunc("/invoices", h.HandleGetInvoices).Methods("GET")
        api.HandleFunc("/invoices/generate", h.HandleGenerateInvoice).Methods("POST")

        port := "5000"
        fmt.Printf("Server starting on http://0.0.0.0:%s\n", port)
        log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/index.html")
}

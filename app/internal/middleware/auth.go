package middleware

import (
        "context"
        "crypto/rand"
        "encoding/base64"
        "log"
        "net/http"
        "os"

        "github.com/gorilla/sessions"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserEmailKey contextKey = "userEmail"

var store *sessions.CookieStore

func init() {
        secret := os.Getenv("SESSION_SECRET")
        if secret == "" {
                secret = generateRandomSecret()
                log.Printf("WARNING: SESSION_SECRET not set, using generated secret. Set SESSION_SECRET environment variable for production.")
        }
        if len(secret) < 32 {
                log.Fatal("SESSION_SECRET must be at least 32 characters long")
        }
        store = sessions.NewCookieStore([]byte(secret))
}

func generateRandomSecret() string {
        b := make([]byte, 32)
        _, err := rand.Read(b)
        if err != nil {
                log.Fatal("Failed to generate session secret:", err)
        }
        return base64.StdEncoding.EncodeToString(b)
}

type ReplitUserInfo struct {
        ID       string   `json:"id"`
        Username string   `json:"username"`
        Name     string   `json:"name"`
        Email    string   `json:"email"`
        Roles    []string `json:"roles"`
}

func AuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                userInfo := r.Header.Get("X-Replit-User-Id")
                userEmail := r.Header.Get("X-Replit-User-Email")
                userName := r.Header.Get("X-Replit-User-Name")

                if userInfo == "" {
                        session, _ := store.Get(r, "auth-session")
                        if uid, ok := session.Values["user_id"]; ok {
                                userInfo = uid.(string)
                                userEmail = session.Values["user_email"].(string)
                        } else {
                                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                                return
                        }
                } else {
                        session, _ := store.Get(r, "auth-session")
                        session.Values["user_id"] = userInfo
                        session.Values["user_email"] = userEmail
                        session.Values["user_name"] = userName
                        session.Save(r, w)
                }

                ctx := context.WithValue(r.Context(), UserIDKey, userInfo)
                ctx = context.WithValue(ctx, UserEmailKey, userEmail)
                next.ServeHTTP(w, r.WithContext(ctx))
        })
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                userInfo := r.Header.Get("X-Replit-User-Id")
                userEmail := r.Header.Get("X-Replit-User-Email")

                if userInfo != "" {
                        ctx := context.WithValue(r.Context(), UserIDKey, userInfo)
                        ctx = context.WithValue(ctx, UserEmailKey, userEmail)
                        next.ServeHTTP(w, r.WithContext(ctx))
                } else {
                        next.ServeHTTP(w, r)
                }
        })
}

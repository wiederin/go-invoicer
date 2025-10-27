# Invoice Generator Web Application

A subscription-based SaaS platform for generating professional PDF invoices with usage limits and user accounts.

> **Note**: This is a private application that uses the public [go-invoicer](https://github.com/wiederin/go-invoicer) library.

## Overview

This web application provides a complete invoice generation platform with:
- User authentication via Replit Auth
- Subscription-based pricing tiers
- Monthly usage quotas
- Invoice history tracking
- RESTful API
- Modern web interface

## Subscription Plans

| Plan | Monthly Quota | Price | Best For |
|------|--------------|-------|----------|
| **Free** | 20 invoices | $0/month | Individuals & freelancers |
| **Basic** | 100 invoices | $9.99/month | Small businesses |
| **Pro** | 500 invoices | $29.99/month | Growing companies |
| **Business** | Unlimited | $99.99/month | Enterprises |

## Features

### Core Functionality
- ✅ PDF invoice generation using go-invoicer library
- ✅ User account management
- ✅ Monthly usage tracking and quota enforcement
- ✅ Invoice history with metadata storage
- ✅ Real-time usage statistics
- 🔄 Stripe payment integration (coming soon)
- 🔄 Email invoice delivery (coming soon)

### Technical Features
- PostgreSQL database for data persistence
- RESTful API with authentication middleware
- Session-based authentication
- Responsive web UI
- Secure quota enforcement at the API level

## Architecture

```
app/
├── cmd/server/           # Main application entry point
├── internal/
│   ├── auth/            # Authentication logic
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, logging)
│   ├── models/          # Data models
│   └── services/        # Business logic (users, usage, invoices)
├── migrations/          # Database schema migrations
├── static/              # Frontend assets (HTML, CSS, JS)
└── go.mod              # Go module definition
```

## Database Schema

### Tables

**plans**
- Subscription tier definitions
- Quota limits and pricing

**users**
- User accounts and authentication
- Plan assignments
- Stripe customer IDs

**subscriptions**
- Active subscription tracking
- Billing period management
- Stripe subscription IDs

**usage_records**
- Monthly usage tracking
- Period-based quotas
- Usage history

**invoices**
- Generated invoice metadata
- Customer and company details
- Creation timestamps

## API Endpoints

### User Management
- `GET /api/user` - Get current user info and usage stats
- `GET /api/usage` - Get current month's usage details

### Subscriptions
- `GET /api/plans` - List all available subscription plans
- `POST /api/subscription/upgrade` - Upgrade to a paid plan (coming soon)

### Invoice Generation
- `POST /api/invoices/generate` - Generate new invoice (enforces quota)
- `GET /api/invoices` - Get user's invoice history
- `GET /api/invoices/:id/download` - Download specific invoice (coming soon)

## Setup and Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database
- Replit account (for authentication)

### Environment Variables

**Required:**
```bash
DATABASE_URL=postgresql://user:password@host:port/dbname
SESSION_SECRET=your-secure-random-secret-minimum-32-chars
```

**Optional:**
```bash
REPLIT_DB_URL=<auto-provided-by-replit>
STRIPE_SECRET_KEY=sk_...
STRIPE_PUBLISHABLE_KEY=pk_...
```

**Generating SESSION_SECRET:**
```bash
# On Linux/Mac
openssl rand -base64 32

# Or use a password generator to create a 32+ character random string
```

### Installation

1. Clone this repository (private repo)
2. Install dependencies:
```bash
cd app
go mod download
```

3. Run database migrations:
```bash
# Migrations run automatically on server start
```

4. Start the server:
```bash
go run ./cmd/server
```

The application will be available at `http://localhost:5000`

## Development

### Project Dependencies

This application depends on the public [go-invoicer](https://github.com/wiederin/go-invoicer) library for PDF generation.

During development, use a local reference:
```go
// In go.mod
replace github.com/wiederin/go-invoicer => ../
```

For production, reference the published version:
```go
// In go.mod
require github.com/wiederin/go-invoicer v1.0.0
```

### Adding New Features

1. Update database schema in `migrations/`
2. Add/update models in `internal/models/`
3. Implement business logic in `internal/services/`
4. Create API handlers in `internal/handlers/`
5. Update frontend in `static/`

## Usage Limit Enforcement

The application enforces usage limits at the API level:

1. User makes request to generate invoice
2. Middleware authenticates user
3. Handler checks current monthly usage
4. If quota exceeded, returns 402 Payment Required
5. If quota available, generates invoice and increments usage
6. Returns PDF to user

## Security

- ✅ Session-based authentication
- ✅ Database-backed usage tracking
- ✅ Secure password storage (when email/password auth added)
- ✅ CSRF protection for web forms
- ✅ SQL injection prevention via parameterized queries
- 🔄 Rate limiting (coming soon)
- 🔄 API key authentication (coming soon)

## Deployment

### Replit Deployment

1. Configure deployment in Replit
2. Set environment variables
3. Deploy to production

### Manual Deployment

```bash
# Build the application
cd app
go build -o invoice-server ./cmd/server

# Run the server
./invoice-server
```

## Roadmap

- [ ] Stripe payment integration
- [ ] Email invoice delivery
- [ ] PDF storage in object storage (S3-compatible)
- [ ] Invoice templates
- [ ] Team accounts with multiple users
- [ ] API keys for programmatic access
- [ ] Webhook notifications
- [ ] Custom branding per account

## Support

For questions or issues specific to this application, contact the development team.

For issues with the underlying invoice generation library, see the [go-invoicer repository](https://github.com/wiederin/go-invoicer).

---

**Private Repository** - Authorized access only

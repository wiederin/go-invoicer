# Invoice Generator Project

## Project Structure

This project is organized into two separate components:

### 1. Public Library: `go-invoicer`
- **Location**: Root directory (components/, constants/, build.go, generator.go, etc.)
- **Purpose**: Reusable Go library for generating PDF invoices
- **Visibility**: Public repository
- **Module**: `github.com/wiederin/go-invoicer`

### 2. Private Web Application: `app/`
- **Location**: `app/` directory
- **Purpose**: Subscription-based SaaS application with user accounts and usage limits
- **Visibility**: Should be in a separate private repository
- **Module**: `github.com/wiederin/go-invoicer-app`

## Architecture

The web application uses the public library as a dependency. The `app/go.mod` references the library using a local replace directive for development:

```go
replace github.com/wiederin/go-invoicer => ../
```

## Separating Into Two Repositories

To deploy this architecture:

1. **Public Library Repository** (this repo)
   - Keep: All library code in the root
   - Remove: The `app/` directory
   - Publish to GitHub as a public repository

2. **Private Web App Repository** (new repo)
   - Move: The entire `app/` directory to a new repository
   - Update `app/go.mod` to reference the public library:
     ```go
     require github.com/wiederin/go-invoicer v1.0.0
     ```
   - Remove the `replace` directive once the library is published

## Database Schema

The web application uses PostgreSQL with the following tables:
- **plans**: Subscription tiers (Free: 20/month, Basic: 100/month, Pro: 500/month, Business: unlimited)
- **users**: User accounts with plan assignments
- **subscriptions**: Stripe subscription tracking
- **usage_records**: Monthly usage tracking per user
- **invoices**: Generated invoice history

## Features

### Library (Public)
- PDF invoice generation
- Configurable company and customer details
- Line items with quantities and prices
- Custom headers, footers, and notes
- Multi-currency support

### Web App (Private)
- Replit authentication integration
- Usage-based subscription plans
- Monthly quota enforcement
- Invoice history tracking
- RESTful API
- Modern web UI

## Subscription Plans

| Plan | Monthly Quota | Price |
|------|--------------|-------|
| Free | 20 invoices | $0 |
| Basic | 100 invoices | $9.99 |
| Pro | 500 invoices | $29.99 |
| Business | Unlimited | $99.99 |

## API Endpoints

- `GET /api/user` - Get user info and usage stats
- `GET /api/usage` - Get current month's usage
- `GET /api/plans` - List available subscription plans
- `GET /api/invoices` - Get invoice history
- `POST /api/invoices/generate` - Generate new invoice (checks quota)

## Development Notes

- The web app runs on port 5000
- Database migrations are handled in `app/migrations/`
- Static files are served from `app/static/`
- Uses Gorilla Mux for routing
- Uses Gorilla Sessions for session management

## Future Enhancements

- Stripe integration for payment processing
- PDF storage in object storage
- Email delivery of invoices
- Invoice templates
- Multi-user team accounts

# Invoice Generator Application

## Project Overview
A subscription-based invoice PDF generator built with Go. Users can create professional PDF invoices through a web interface, with usage limits based on their subscription tier.

## Project Structure

### Library (`/` root directory)
The core invoice generation library (`go-invoicer`) - a reusable Go module for creating PDF invoices.

**Key Files:**
- `generator.go` - Main invoice document initialization
- `build.go` - PDF document builder
- `components/` - Invoice components (address, items, tax, etc.)
- `constants/` - Configuration constants

### Web Application (`/app` directory)
Full-stack web application that uses the library.

**Structure:**
```
app/
├── cmd/server/           # Main server entry point
├── internal/
│   ├── models/          # Data models (User, Plan, Invoice, etc.)
│   ├── database/        # Database connection and migrations
│   ├── services/        # Business logic (User, Usage, Invoice services)
│   ├── handlers/        # HTTP API handlers
│   └── middleware/      # Authentication middleware
├── static/              # Frontend files (HTML, CSS, JS)
└── migrations/          # Database schema migrations
```

## Features

### Subscription Plans
- **Free**: 20 invoices/month ($0)
- **Basic**: 100 invoices/month ($9.99)
- **Pro**: 500 invoices/month ($29.99)
- **Business**: Unlimited invoices ($99.99)

### Current Capabilities
1. **User Authentication**: Replit Auth integration (auto-creates user accounts)
2. **Usage Tracking**: Monthly quota enforcement per plan
3. **Invoice Generation**: Professional PDF invoices with company and customer details
4. **Invoice History**: View all generated invoices
5. **Quota Management**: Real-time usage tracking and limits

## Database Schema

**Tables:**
- `plans` - Subscription plans and pricing
- `users` - User accounts with Replit Auth integration
- `subscriptions` - User subscription details
- `usage_records` - Monthly usage tracking
- `invoices` - Invoice metadata and history

## API Endpoints

All API endpoints require authentication via Replit Auth headers.

- `GET /api/user` - Get current user and usage stats
- `GET /api/usage` - Get current month usage
- `GET /api/plans` - List all subscription plans
- `GET /api/invoices` - Get invoice history
- `POST /api/invoices/generate` - Generate new PDF invoice (enforces quotas)

## How It Works

1. **Authentication**: Users are authenticated via Replit Auth headers (`X-Replit-User-Id`, `X-Replit-User-Email`)
2. **User Creation**: First-time users are automatically created with the Free plan
3. **Invoice Generation**:
   - Check if user is within monthly quota
   - Generate PDF using go-invoicer library
   - Save invoice metadata to database
   - Increment usage counter
   - Return PDF to user
4. **Quota Enforcement**: Returns 402 Payment Required if quota exceeded

## Running the Application

The server runs on port 5000:
```bash
cd app && go run ./cmd/server
```

## Environment Variables
- `DATABASE_URL` - PostgreSQL connection string (automatically provided by Replit)
- Other `PG*` variables for database connection

## Future Enhancements (Not Implemented)

### Stripe Integration
For payment processing and subscription management:
- User can upgrade/downgrade plans
- Webhook handling for subscription events
- Customer portal for billing management

### Additional Features to Consider
- Invoice templates and customization
- Logo upload for company branding
- PDF storage and download links
- Email delivery of invoices
- Multi-currency support
- Tax calculations
- Recurring invoice templates

## Development Notes

- The library and web app are separate Go modules
- The app uses a `replace` directive in `go.mod` to use the local library
- Database migrations are in SQL format
- Frontend is vanilla JavaScript (no framework)
- Authentication uses Replit's built-in auth system

## Recent Changes
- Initial project setup with library extraction
- Database schema created with subscription and usage tracking
- REST API implemented with quota enforcement
- Web UI created for invoice generation
- Free tier with 20 invoices/month enabled

Last Updated: October 27, 2025

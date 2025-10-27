# Guide: Separating Library and Web App into Two Repositories

This guide explains how to split the current monorepo into two separate repositories:
1. **Public Library**: `go-invoicer` (invoice PDF generation)
2. **Private Web App**: `go-invoicer-app` (SaaS application)

## Current Structure

```
workspace/
├── components/          # Library code
├── constants/           # Library code
├── build.go            # Library code
├── generator.go        # Library code
├── generator_test.go   # Library code
├── go.mod              # Library go.mod
├── go.sum              # Library go.sum
├── README.md           # Library README
└── app/                # Web application code
    ├── cmd/
    ├── internal/
    ├── migrations/
    ├── static/
    ├── go.mod          # App go.mod
    ├── go.sum          # App go.sum
    └── README.md       # App README
```

## Step 1: Create Public Library Repository

### On GitHub

1. Create a new **public** repository: `go-invoicer`
2. Clone it locally or on Replit

### Files to Include

Copy these files from the current workspace to the public library repo:

```
go-invoicer/
├── components/
│   ├── address.go
│   ├── config.go
│   ├── contact.go
│   ├── description.go
│   ├── discount.go
│   ├── document.go
│   ├── header_footer.go
│   ├── item.go
│   ├── meta.go
│   ├── notes.go
│   ├── setters.go
│   ├── tax.go
│   ├── term.go
│   ├── title.go
│   └── totals.go
├── constants/
│   └── constants.go
├── build.go
├── generator.go
├── generator_test.go
├── go.mod
├── go.sum
├── .gitignore
├── LICENSE
└── README.md
```

### Commands

```bash
# In the public repo
git init
git add .
git commit -m "Initial commit: Invoice PDF generation library"
git remote add origin https://github.com/your-username/go-invoicer.git
git push -u origin main
```

### Tag a Version

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Step 2: Create Private Web App Repository

### On GitHub

1. Create a new **private** repository: `go-invoicer-app`
2. Clone it locally or on Replit

### Files to Include

Copy the entire `app/` directory contents to the private repo root:

```
go-invoicer-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── services/
├── migrations/
│   └── 001_initial_schema.sql
├── static/
│   ├── index.html
│   ├── style.css
│   └── app.js
├── go.mod
├── go.sum
├── .gitignore
├── README.md
└── .replit (optional - for Replit deployment)
```

### Update go.mod

Replace the local reference with the published library version:

**Before:**
```go
module github.com/your-username/go-invoicer-app

go 1.21

require (
	github.com/wiederin/go-invoicer v0.0.0
	// ... other dependencies
)

replace github.com/wiederin/go-invoicer => ../
```

**After:**
```go
module github.com/your-username/go-invoicer-app

go 1.21

require (
	github.com/your-username/go-invoicer v1.0.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.4.0
	github.com/lib/pq v1.10.9
)
```

### Update Import Paths

Change all imports in the code from:
```go
import "github.com/wiederin/go-invoicer-app/internal/..."
```

To your actual GitHub username:
```go
import "github.com/your-username/go-invoicer-app/internal/..."
```

### Commands

```bash
# In the private repo
cd path/to/go-invoicer-app
go mod tidy
git init
git add .
git commit -m "Initial commit: Invoice SaaS application"
git remote add origin https://github.com/your-username/go-invoicer-app.git
git push -u origin main
```

## Step 3: Environment Variables

In your private repository, ensure you have these environment variables set:

```bash
# Database
DATABASE_URL=postgresql://...

# Replit Auth (if deploying on Replit)
REPLIT_DB_URL=<provided-by-replit>

# Future: Stripe
STRIPE_SECRET_KEY=sk_...
STRIPE_PUBLISHABLE_KEY=pk_...
```

## Step 4: Deployment

### Deploying on Replit

1. Create a new Replit from the private GitHub repository
2. Set environment variables in Secrets
3. Configure deployment:
   - Build Command: None (Go compiles on run)
   - Run Command: `cd app && go run ./cmd/server`  → UPDATE to just `go run ./cmd/server` in new repo
   - Port: 5000

### Deploying Elsewhere

```bash
# Build the application
go build -o invoice-server ./cmd/server

# Run
./invoice-server
```

## Step 5: Maintaining Both Repositories

### When Updating the Library

1. Make changes in the public `go-invoicer` repository
2. Commit and tag a new version:
   ```bash
   git tag v1.1.0
   git push origin v1.1.0
   ```
3. Update the private app's `go.mod`:
   ```bash
   go get github.com/your-username/go-invoicer@v1.1.0
   go mod tidy
   ```

### Development Workflow

For local development with unreleased library changes:

1. Clone both repositories as siblings:
   ```
   projects/
   ├── go-invoicer/
   └── go-invoicer-app/
   ```

2. Temporarily use a local replace in `go-invoicer-app/go.mod`:
   ```go
   replace github.com/your-username/go-invoicer => ../go-invoicer
   ```

3. Before committing, remove the replace directive and use a version tag

## Benefits of Separation

✅ **Library stays public** - Others can use your invoice generation library
✅ **App stays private** - Your business logic and secret sauce remain confidential
✅ **Clear separation** - Business logic vs. reusable library
✅ **Independent versioning** - Update each at its own pace
✅ **Better documentation** - Each repo has focused documentation
✅ **Security** - Sensitive credentials only in private repo

## Checklist

- [ ] Create public `go-invoicer` repository
- [ ] Copy library files and push
- [ ] Tag version v1.0.0
- [ ] Create private `go-invoicer-app` repository
- [ ] Copy app files to new repo root
- [ ] Update go.mod to reference public library
- [ ] Update import paths with correct GitHub username
- [ ] Set up environment variables
- [ ] Test build locally
- [ ] Deploy to production
- [ ] Update documentation links

## Need Help?

- **Library issues**: Open issue on public repo
- **App issues**: Contact your development team
- **Deployment**: See Replit or hosting provider documentation

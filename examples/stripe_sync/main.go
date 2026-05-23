// Example: fetch a Stripe invoice via API and render HTML.
//
// Usage:
//
//	STRIPE_SECRET_KEY=sk_test_... go run ./examples/stripe_sync -id in_xxx
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wiederin/go-invoicer/integrations/stripe"
	"github.com/wiederin/go-invoicer/render"
)

func main() {
	invoiceID := flag.String("id", "", "Stripe invoice ID (in_xxx)")
	out := flag.String("out", "stripe-invoice.html", "output HTML path")
	flag.Parse()

	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		log.Fatal("set STRIPE_SECRET_KEY")
	}
	if *invoiceID == "" {
		log.Fatal("usage: stripe_sync -id in_xxx")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := stripe.NewClient(key)
	inv, err := client.SyncInvoice(ctx, *invoiceID, stripe.ImportOptions{
		Supplier: stripe.Supplier{Name: os.Getenv("STRIPE_SUPPLIER_NAME")},
	})
	if err != nil {
		log.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		log.Fatal(err)
	}
	html, err := engine.RenderDefault(inv)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", *out, "for", inv.Number)
}

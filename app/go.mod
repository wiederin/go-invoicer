module github.com/wiederin/go-invoicer-app

go 1.23

toolchain go1.24.9

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.4.0
	github.com/lib/pq v1.10.9
	github.com/wiederin/go-invoicer v0.0.0
)

require (
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/go-pdf/fpdf v0.9.0 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/leekchan/accounting v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
)

replace github.com/wiederin/go-invoicer => ../

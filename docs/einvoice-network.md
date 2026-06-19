# E-invoice network submission

Submit generated XRechnung or PEPPOL UBL XML to an access point or middleware.

## Hosted API

1. Configure your network endpoint (HTTPS webhook):

   ```http
   PATCH /v1/org/compliance-network
   { "enabled": true, "url": "https://ap.example.com/invoices" }
   ```

2. Submit an invoice:

   ```http
   POST /v1/invoices/{id}/compliance/send
   { "format": "peppol", "signed": true, "buyer_endpoint_id": "DE987654321" }
   ```

3. Review delivery log:

   ```http
   GET /v1/org/compliance-deliveries?limit=50
   ```

Without per-org configuration, the platform uses `EINVOICE_NETWORK_SENDER` (`log` records to API stdout, `webhook` posts to `EINVOICE_NETWORK_WEBHOOK_URL`).

## XML signing

Append `signed: true` on send or export. Configure the signer on the API:

| Env | Values |
|-----|--------|
| `EINVOICE_SIGNER` | `dev` (default), `file`, `remote` |
| `EINVOICE_SIGN_KEY_PATH` | RSA private key PEM when `file` |
| `EINVOICE_SIGN_REMOTE_URL` | Signing service POST URL when `remote` |

Production QES/HSM providers can implement the remote signer contract or replace the `file` signer with a mounted HSM key.

## Related

- [Document types](document-types.md)
- [Hosted platform](hosted-platform.md)

package swiss

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

// QRPNGBase64 encodes payload as a PNG QR code (base64, no data-URL prefix).
func QRPNGBase64(payload string, size int) (string, error) {
	if size <= 0 {
		size = 256
	}
	png, err := qrcode.Encode(payload, qrcode.Medium, size)
	if err != nil {
		return "", fmt.Errorf("swiss qr: encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

// QRDataURL returns a data:image/png URL for HTML embedding.
func QRDataURL(payload string, size int) (string, error) {
	b64, err := QRPNGBase64(payload, size)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + b64, nil
}

package imagehelper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png" // for decoding PNG
	"strings"

	"golang.org/x/image/draw"
)

// maxSize is 5MB. In base64, size = bytes * 1.33
const maxBase64Size = 5 * 1024 * 1024 * 4 / 3

// ResizeBase64Image resizes and compresses a base64 image (reducing to max width 800px)
func ResizeBase64Image(base64Str string) (string, error) {
	if len(base64Str) > maxBase64Size {
		return "", errors.New("image size exceeds 5MB")
	}

	parts := strings.Split(base64Str, ",")
	b64Data := base64Str
	prefix := ""
	if len(parts) == 2 {
		prefix = parts[0] + ","
		b64Data = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Max width
	maxWidth := 800
	if width > maxWidth {
		ratio := float64(maxWidth) / float64(width)
		width = maxWidth
		height = int(float64(height) * ratio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	// Encode as JPEG with 75% quality
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 75}); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	if prefix != "" {
		if !strings.Contains(prefix, "image/jpeg") { // standardize if we convert to jpeg
			prefix = "data:image/jpeg;base64,"
		}
		return prefix + encoded, nil
	}
	return "data:image/jpeg;base64," + encoded, nil
}

package middlewares

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HmacMiddleware struct {
	secretKey     string
	maxAgeSeconds int
	redisClient   *redis.Client
}

func NewHmacMiddleware(secretKey string, maxAgeSeconds int, redisClient *redis.Client) *HmacMiddleware {
	return &HmacMiddleware{
		secretKey:     secretKey,
		maxAgeSeconds: maxAgeSeconds,
		redisClient:   redisClient,
	}
}

func (h *HmacMiddleware) Validate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timestamp := c.Request().Header.Get("X-Timestamp")
			nonce := c.Request().Header.Get("X-Nonce")
			signature := c.Request().Header.Get("X-Signature")

			if timestamp == "" || nonce == "" || signature == "" {
				return response.Error(c, appErrors.ErrMissingHmacHeaders)
			}

			// Validate timestamp
			ts, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				return response.Error(c, appErrors.ErrInvalidTimestamp)
			}

			now := time.Now().Unix()
			diff := math.Abs(float64(now - ts))
			if diff > float64(h.maxAgeSeconds) {
				return response.Error(c, appErrors.ErrRequestExpired)
			}

			// Read body for signature verification, then restore it
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return response.Error(c, appErrors.ErrInvalidInput)
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Verify HMAC signature: HMAC-SHA256(timestamp + nonce + body, secret)
			message := fmt.Sprintf("%s.%s.%s", timestamp, nonce, string(bodyBytes))
			mac := hmac.New(sha256.New, []byte(h.secretKey))
			mac.Write([]byte(message))
			expectedSignature := hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
				return response.Error(c, appErrors.ErrInvalidSignature)
			}

			// Check nonce uniqueness via Redis (prevent replay)
			nonceKey := fmt.Sprintf("hmac:nonce:%s", nonce)
			ctx := context.Background()
			set, err := h.redisClient.SetNX(ctx, nonceKey, 1, time.Duration(h.maxAgeSeconds)*time.Second).Result()
			if err != nil {
				return response.Error(c, appErrors.ErrInternal)
			}
			if !set {
				return response.Error(c, appErrors.ErrNonceReused)
			}

			return next(c)
		}
	}
}

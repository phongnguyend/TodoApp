package security_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/security"
)

func jwt(payload, secret string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	body := base64.RawURLEncoding.EncodeToString([]byte(payload))
	content := header + "." + body
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(content))
	return content + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := security.HashPassword("password123", 1000)
	require.NoError(t, err)
	assert.True(t, security.VerifyPassword("password123", hash))
	assert.False(t, security.VerifyPassword("wrong", hash))
}

func TestAuthenticatedUserIDValidatesJWT(t *testing.T) {
	now := time.Unix(2_000_000_000, 0)
	token := jwt(fmt.Sprintf(`{"sub":"42","exp":%d}`, now.Unix()+60), "secret")
	id, err := security.AuthenticatedUserID("Bearer "+token, "secret", now)
	require.NoError(t, err)
	assert.Equal(t, uint(42), id)
	_, err = security.AuthenticatedUserID("Bearer "+token, "wrong", now)
	assert.ErrorIs(t, err, security.ErrInvalidToken)
}

func TestResetTokenExpires(t *testing.T) {
	now := time.Unix(2_000_000_000, 0)
	token, err := security.CreateResetToken(7, "hash", "secret", now.Add(time.Minute))
	require.NoError(t, err)
	id, fingerprint, err := security.DecodeResetToken(token, "secret", now)
	require.NoError(t, err)
	assert.Equal(t, uint(7), id)
	assert.Equal(t, security.PasswordFingerprint("hash"), fingerprint)
	_, _, err = security.DecodeResetToken(token, "secret", now.Add(2*time.Minute))
	assert.ErrorIs(t, err, security.ErrInvalidToken)
}

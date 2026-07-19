package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/pbkdf2"
)

var ErrInvalidToken = errors.New("invalid or expired token")

func HashPassword(password string, iterations int) (string, error) {
	if iterations < 1 {
		iterations = 120000
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	digest := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iterations, hex.EncodeToString(salt), hex.EncodeToString(digest)), nil
}

func VerifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2_sha256" {
		return false
	}
	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations < 1 {
		return false
	}
	salt, err := hex.DecodeString(parts[2])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[3])
	if err != nil {
		return false
	}
	actual := pbkdf2.Key([]byte(password), salt, iterations, len(expected), sha256.New)
	return hmac.Equal(actual, expected)
}

func encode(data []byte) string { return base64.RawURLEncoding.EncodeToString(data) }
func sign(content, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(content))
	return encode(mac.Sum(nil))
}
func validSignature(content, supplied, secret string) bool {
	a, err := base64.RawURLEncoding.DecodeString(supplied)
	if err != nil {
		return false
	}
	b, _ := base64.RawURLEncoding.DecodeString(sign(content, secret))
	return hmac.Equal(a, b)
}

type tokenPayload struct {
	Sub      any    `json:"sub"`
	Exp      int64  `json:"exp,omitempty"`
	Password string `json:"password,omitempty"`
}

func CreateJWT(userID uint, secret string, issuedAt, expires time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expires),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func payloadUserID(v any) (uint, error) {
	var raw string
	switch n := v.(type) {
	case string:
		raw = n
	case float64:
		raw = strconv.FormatInt(int64(n), 10)
	default:
		return 0, ErrInvalidToken
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, ErrInvalidToken
	}
	return uint(id), nil
}

// AuthenticatedUserID validates a standard HS256 JWT from an Authorization header.
func AuthenticatedUserID(header, secret string, now time.Time) (uint, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return 0, ErrInvalidToken
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		parts[1],
		claims,
		func(token *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}
	return payloadUserID(claims.Subject)
}

func PasswordFingerprint(hash string) string {
	sum := sha256.Sum256([]byte(hash))
	return hex.EncodeToString(sum[:])
}

func CreateResetToken(userID uint, passwordHash, secret string, expires time.Time) (string, error) {
	body, err := json.Marshal(tokenPayload{Sub: strconv.FormatUint(uint64(userID), 10), Exp: expires.Unix(), Password: PasswordFingerprint(passwordHash)})
	if err != nil {
		return "", err
	}
	encoded := encode(body)
	return encoded + "." + sign(encoded, secret), nil
}

func DecodeResetToken(token, secret string, now time.Time) (uint, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 || !validSignature(parts[0], parts[1], secret) {
		return 0, "", ErrInvalidToken
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, "", ErrInvalidToken
	}
	var payload tokenPayload
	if json.Unmarshal(body, &payload) != nil || payload.Exp < now.Unix() || payload.Password == "" {
		return 0, "", ErrInvalidToken
	}
	id, err := payloadUserID(payload.Sub)
	return id, payload.Password, err
}

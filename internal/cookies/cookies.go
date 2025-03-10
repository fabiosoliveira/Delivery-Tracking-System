package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Define custom errors for cookie value constraints.
var (
	ErrValueTooLong = errors.New("cookie value exceeds size limit")
	ErrInvalidValue = errors.New("cookie value is not valid")
)

// Write safely encodes and sets the cookie if within size limits.
func Write(w http.ResponseWriter, cookie *http.Cookie) error {
	// Base64 encode the cookie value.
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Validate the cookie size, ensuring it doesn't surpass 4096 bytes.
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	// Set the cookie in the HTTP response.
	http.SetCookie(w, cookie)

	return nil
}

// Read decodes a base64-encoded cookie value from the incoming request.
func Read(r *http.Request, name string) (string, error) {
	// Retrieve the cookie by its name.
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	// Base64 decode the cookie value, returning an error for invalid encoding.
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	// Return the decoded value.
	return string(value), nil
}

// WriteSigned enhances the cookie value with a HMAC signature.
func WriteSigned(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) error {
	// Generate HMAC using SHA256 and the secret key.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	// Combine the signature with the original cookie value.
	cookie.Value = string(signature) + cookie.Value

	// Utilize the Write function to encode and set the cookie.
	return Write(w, cookie)
}

// ReadSigned extracts and verifies the HMAC signature from the cookie.
func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	// Retrieve the signed value, which includes the signature and the value.
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	// Ensure the signed value is sufficiently long to contain the signature.
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	// Separate the signature from the value.
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	// Recalculate the HMAC to confirm the cookie's integrity.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	// Validate the signature against the expected one.
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	// Return the original, verified cookie value.
	return value, nil
}

func WriteEncrypted(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) error {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	cookie.Value = base64.URLEncoding.EncodeToString(encryptedValue)

	http.SetCookie(w, cookie)
	return nil
}

func ReadEncrypted(r *http.Request, name string, secretKey []byte) (string, error) {
	encryptedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.URLEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return "", ErrInvalidValue
	}

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok || expectedName != name {
		return "", ErrInvalidValue
	}

	return value, nil
}

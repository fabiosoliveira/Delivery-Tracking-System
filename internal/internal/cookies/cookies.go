package cookies

import (
	"github.com/gorilla/securecookie"
)

// Hash keys should be at least 32 bytes long
var hashKey = []byte("12345678901234567890123456789012")

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte("1234567890123456")
var S = securecookie.New(hashKey, blockKey)

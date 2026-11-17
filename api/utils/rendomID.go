package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var objectIDCounter uint32

func init() {
	var b [4]byte
	_, _ = io.ReadFull(rand.Reader, b[:])
	objectIDCounter = binary.BigEndian.Uint32(b[:])
}

// GenerateObjectID generates a 24-character hexadecimal unique identifier
// identical to the MongoDB ObjectID specification:
// 4-byte unix timestamp + 5-byte cryptographically secure random value + 3-byte atomic counter
// Example output: "670c5e7b8a1c9f2d3e4b5a6c"
func GenerateID() string {
	var b [12]byte

	// 4 bytes: timestamp (seconds since unix epoch)
	binary.BigEndian.PutUint32(b[0:4], uint32(time.Now().Unix()))

	// 5 bytes: crypto random value
	_, _ = io.ReadFull(rand.Reader, b[4:9])

	// 3 bytes: process-level atomic counter to guarantee uniqueness under high concurrency
	c := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(c >> 16)
	b[10] = byte(c >> 8)
	b[11] = byte(c)

	return hex.EncodeToString(b[:])
}

// GenerateRandomID is an alias for GenerateObjectID for generic ID usage across modules
func GenerateRandomID() string {
	return GenerateID()
}

// GeneratePrefixedID generates a MongoDB-style ObjectID with a custom module prefix
// Example: GeneratePrefixedID("cat_") => "cat_670c5e7b8a1c9f2d3e4b5a6c"
func GeneratePrefixedID(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, GenerateID())
}

// IsValidObjectID checks if the given string is a valid 24-character hex MongoDB ObjectID
func IsValidObjectID(id string) bool {
	if len(id) != 24 {
		return false
	}
	_, err := hex.DecodeString(id)
	return err == nil
}

package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Fingerprint creates a deterministic hash of the request context
func Fingerprint(ctx map[string]any) string {
	// Create a sorted list of keys for deterministic ordering
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create a map with sorted keys
	sortedCtx := make(map[string]any)
	for _, k := range keys {
		sortedCtx[k] = ctx[k]
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(sortedCtx)
	if err != nil {
		// Fallback to a simple string representation
		return simpleFingerprint(ctx)
	}

	// Create SHA256 hash
	hash := sha256.Sum256(jsonData)

	// Return the first 16 bytes as hex string
	return hex.EncodeToString(hash[:16])
}

// simpleFingerprint creates a simple fingerprint as fallback
func simpleFingerprint(ctx map[string]any) string {
	// Create a simple string representation
	var result string
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result += k + ":" + toString(ctx[k]) + ";"
	}

	// Create a hash of the string
	hash := sha256.Sum256([]byte(result))

	return hex.EncodeToString(hash[:16])
}

// toString converts any value to string
func toString(v any) string {
	return fmt.Sprint(v)
}

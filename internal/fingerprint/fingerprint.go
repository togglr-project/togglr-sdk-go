package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

func Fingerprint(ctx map[string]any) string {
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedCtx := make(map[string]any)
	for _, k := range keys {
		sortedCtx[k] = ctx[k]
	}

	jsonData, err := json.Marshal(sortedCtx)
	if err != nil {
		return simpleFingerprint(ctx)
	}

	hash := sha256.Sum256(jsonData)

	return hex.EncodeToString(hash[:16])
}

func simpleFingerprint(ctx map[string]any) string {
	var result string
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result += k + ":" + toString(ctx[k]) + ";"
	}

	hash := sha256.Sum256([]byte(result))

	return hex.EncodeToString(hash[:16])
}

func toString(v any) string {
	return fmt.Sprint(v)
}

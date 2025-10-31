package togglr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

func TestNewTrackEvent(t *testing.T) {
	te := NewTrackEvent("test-variant", EventTypeSuccess)
	assert.Equal(t, "test-variant", te.VariantKey)
	assert.Equal(t, EventTypeSuccess, te.EventType)
	assert.Nil(t, te.Reward)
	assert.Empty(t, te.Context)
	assert.Nil(t, te.CreatedAt)
	assert.Nil(t, te.DedupKey)
}

func TestTrackEventMethods(t *testing.T) {
	te := NewTrackEvent("test-variant", EventTypeSuccess).
		WithReward(10.5).
		WithContext("key1", "value1").
		WithContexts(map[string]any{"key2": 123, "key3": true}).
		WithCreatedAt(time.Now()).
		WithDedupKey("dedup-123")

	assert.Equal(t, float32(10.5), *te.Reward)
	assert.Equal(t, "value1", te.Context["key1"])
	assert.Equal(t, 123, te.Context["key2"])
	assert.True(t, te.Context["key3"].(bool))
	assert.NotNil(t, te.CreatedAt)
	assert.Equal(t, "dedup-123", *te.DedupKey)
}

func TestToAPIRequest(t *testing.T) {
	createdAt := time.Now()
	te := NewTrackEvent("test-variant", EventTypeSuccess).
		WithReward(10.5).
		WithContext("key1", "value1").
		WithCreatedAt(createdAt).
		WithDedupKey("dedup-123")

	req := te.toAPIRequest()
	assert.Equal(t, "test-variant", req.VariantKey)
	assert.Equal(t, api.TrackRequestEventType(EventTypeSuccess), req.EventType)
	vReward, okReward := req.Reward.Get()
	assert.True(t, okReward)
	assert.Equal(t, float32(10.5), vReward)
	assert.NotNil(t, req.Context)
	vCreatedAt, okCreatedAt := req.CreatedAt.Get()
	assert.True(t, okCreatedAt)
	assert.Equal(t, createdAt, vCreatedAt)
	vDedupKey, okDedupKey := req.DedupKey.Get()
	assert.True(t, okDedupKey)
	assert.Equal(t, "dedup-123", vDedupKey)
}

func TestToAPIRequestEmptyFields(t *testing.T) {
	te := NewTrackEvent("test-variant", EventTypeSuccess)
	req := te.toAPIRequest()
	assert.Equal(t, "test-variant", req.VariantKey)
	assert.Equal(t, api.TrackRequestEventType(EventTypeSuccess), req.EventType)
	_, okRewardEmpty := req.Reward.Get()
	assert.False(t, okRewardEmpty)
	assert.Nil(t, req.Context.Value)
	_, okCreatedAtEmpty := req.CreatedAt.Get()
	assert.False(t, okCreatedAtEmpty)
	_, okDedupKeyEmpty := req.DedupKey.Get()
	assert.False(t, okDedupKeyEmpty)
}

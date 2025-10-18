package togglr

import (
	"encoding/json"
	"time"

	"github.com/go-faster/jx"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

type TrackEvent struct {
	VariantKey string
	EventType  EventType
	Reward     *float32
	Context    RequestContext
	CreatedAt  *time.Time
	DedupKey   *string
}

type EventType string

const (
	EventTypeSuccess EventType = "success"
	EventTypeFailure EventType = "failure"
	EventTypeError   EventType = "error"
)

func NewTrackEvent(variantKey string, eventType EventType) *TrackEvent {
	return &TrackEvent{
		VariantKey: variantKey,
		EventType:  eventType,
		Context:    make(RequestContext),
	}
}

func (te *TrackEvent) WithReward(reward float32) *TrackEvent {
	te.Reward = &reward
	return te
}

func (te *TrackEvent) WithContext(key string, value any) *TrackEvent {
	te.Context[key] = value
	return te
}

func (te *TrackEvent) WithContexts(contexts map[string]any) *TrackEvent {
	for k, v := range contexts {
		te.Context[k] = v
	}
	return te
}

func (te *TrackEvent) WithCreatedAt(createdAt time.Time) *TrackEvent {
	te.CreatedAt = &createdAt
	return te
}

func (te *TrackEvent) WithDedupKey(dedupKey string) *TrackEvent {
	te.DedupKey = &dedupKey
	return te
}

func (te *TrackEvent) toAPIRequest() *api.TrackRequest {
	req := &api.TrackRequest{
		VariantKey: te.VariantKey,
		EventType:  api.TrackRequestEventType(te.EventType),
	}

	if te.Reward != nil {
		req.Reward = api.NewOptFloat32(*te.Reward)
	}

	if len(te.Context) > 0 {
		contextData := make(api.TrackRequestContext)
		for k, v := range te.Context {
			if raw, err := json.Marshal(v); err == nil {
				contextData[k] = jx.Raw(raw)
			}
		}
		req.Context = api.NewOptTrackRequestContext(contextData)
	}

	if te.CreatedAt != nil {
		req.CreatedAt = api.NewOptDateTime(*te.CreatedAt)
	}

	if te.DedupKey != nil {
		req.DedupKey = api.NewOptString(*te.DedupKey)
	}

	return req
}

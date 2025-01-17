package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type Subscription struct {
	Status        string               `json:"status"`
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// GetSubscription - Handles managing subscriptions for changed appointment slots.
// GET /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-7
func (h *HTTPClient) GetSubscription(ctx context.Context, feedType string) (*Subscription, error) {
	out := &Subscription{}

	_, err := h.Get(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), nil, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type SubscriptionEvent struct {
	EventName string `json:"eventname"`
}

type listSubscriptionEventsResponse struct {
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// ListSubscriptionEvents - Returns the list of events you can subscribe to for changed appointment slots.
// GET /v1/{practiceid}/appointments/changed/subscription/events.
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-8
func (h *HTTPClient) ListSubscriptionEvents(ctx context.Context, feedType string) ([]*SubscriptionEvent, error) {
	out := &listSubscriptionEventsResponse{}

	_, err := h.Get(ctx, fmt.Sprintf("/%s/changed/subscription/events", feedType), nil, &out)
	if err != nil {
		return nil, err
	}

	return out.Subscriptions, nil
}

type SubscribeOptions struct {
	EventName string
}

// Subscribe - Handles subscriptions for changed appointment slots.
// POST /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-6
func (h *HTTPClient) Subscribe(ctx context.Context, feedType string, opts *SubscribeOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}
		form.Add("eventname", opts.EventName)
	}

	_, err := h.PostForm(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type UnsubscribeOptions struct {
	EventName string
}

// Unsubscribe - Handles subscriptions for changed appointment slots.
// POST /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-6
func (h *HTTPClient) Unsubscribe(ctx context.Context, feedType string, opts *UnsubscribeOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}
		form.Add("eventname", opts.EventName)
	}

	_, err := h.DeleteForm(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}

package validate

import (
	"errors"
	"testing"
)

func TestWebhookValid(t *testing.T) {
	err := Webhook("wh-1", "http://localhost:9000/hook")
	if err != nil {
		t.Errorf("Webhook() = %v, want nil", err)
	}
}

func TestWebhook(t *testing.T) {
	tests := []struct {
		name string
		id   string
		url  string
		want error
	}{
		{name: "valid http", id: "wh-1", url: "http://example.com", want: nil},
		{name: "valid https", id: "wh-1", url: "https://example.com", want: nil},
		{name: "empty id", id: "", url: "http://example.com", want: ErrEmptyID},
		{name: "empty url", id: "wh-1", url: "", want: ErrEmptyURL},
		{name: "bad scheme", id: "wh-1", url: "ftp://example.com", want: ErrBadScheme},
		{name: "empty id wins over empty url", id: "", url: "", want: ErrEmptyID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Webhook(tt.id, tt.url)
			if !errors.Is(got, tt.want) {
				t.Errorf("Webhook(%q, %q) = %v, want %v", tt.id, tt.url, got, tt.want)
			}
		})
	}
}

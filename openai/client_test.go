package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("unexpected request")
		}
		_, _ = w.Write([]byte(`{"status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"hello"}]}]}`))
	}))
	defer s.Close()
	c.endpoint = s.URL
	got, err := c.Generate(context.Background(), "test-model", "hi")
	if err != nil || got != "hello" {
		t.Fatalf("got %q, error %v", got, err)
	}
}

func TestMissingKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	if _, err := NewClient(); err == nil {
		t.Fatal("expected missing key error")
	}
}

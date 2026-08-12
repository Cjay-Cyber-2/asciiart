package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAsciiAllBanners(t *testing.T) {
	for _, banner := range []string{"standard", "shadow", "thinkertoy"} {
		out, err := ascii("Hello", banner)
		if err != nil {
			t.Fatalf("%s: %v", banner, err)
		}
		if strings.TrimSpace(out) == "" {
			t.Fatalf("%s: empty output", banner)
		}
		if !strings.Contains(out, "\n") {
			t.Fatalf("%s: expected multiline output", banner)
		}
	}
}

func TestAsciiNewlines(t *testing.T) {
	out, err := ascii("A\\nB", "standard")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 16 {
		t.Fatalf("expected at least 16 art rows for two chars, got %d", len(lines))
	}

	out2, err := ascii("A\nB", "shadow")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out2) == "" {
		t.Fatal("real newline input produced empty output")
	}
}

func TestAsciiRejectsBadInput(t *testing.T) {
	if _, err := ascii("héllo", "standard"); err == nil {
		t.Fatal("expected error for non-ascii")
	}
	if _, err := ascii("hi", "missing"); err == nil {
		t.Fatal("expected error for bad banner")
	}
}

func TestHandlerGenerate(t *testing.T) {
	body := strings.NewReader("text=Hi&banner=thinkertoy")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	Handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "OUTPUT STREAM") {
		t.Fatal("missing output panel")
	}
	if !strings.Contains(rr.Body.String(), "value=\"thinkertoy\"") {
		t.Fatal("banner not preserved")
	}
	if !strings.Contains(rr.Body.String(), ">Hi<") && !strings.Contains(rr.Body.String(), ">Hi</textarea>") {
		// textarea should keep text
		if !strings.Contains(rr.Body.String(), "Hi") {
			t.Fatal("text not preserved in form")
		}
	}
}

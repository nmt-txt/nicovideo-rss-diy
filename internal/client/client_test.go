package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSearchVideo_ReturnsVideos(t *testing.T) {
	// create test server that serves files from testdata/client based on q param
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, "missing q", http.StatusBadRequest)
			return
		}

		p := filepath.Join("..", "..", "testdata", "client", q+".json")
		b, err := os.ReadFile(p)
		if err != nil {
			http.Error(w, fmt.Sprintf("not found: %s", q), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(json.RawMessage(b))
	}))
	defer srv.Close()

	c := NewVideoClient(srv.URL, "niconico-rss-diy/0.1 test")
	ctx := context.Background()

	for _, q := range []string{"vocaloid", "software_talk", "vocaloidd"} {
		resp, err := c.SearchVideo(ctx, q, nil)
		if err != nil {
			t.Fatalf("SearchVideo error: %v", err)
		}

		if resp == nil {
			t.Fatalf("expected non-nil response")
		}

		// load fixture to assert equality
		fixturePath := filepath.Join("..", "..", "testdata", "client", q+".json")
		fb, err := os.ReadFile(fixturePath)
		if err != nil {
			t.Fatalf("read fixture: %v", err)
		}
		var expect SearchVideoResponse
		if err := json.Unmarshal(fb, &expect); err != nil {
			t.Fatalf("unmarshal fixture: %v", err)
		}

		if len(resp.Videos) != len(expect.Videos) {
			t.Fatalf("expected %d videos, got %d", len(expect.Videos), len(resp.Videos))
		}

		// check first ID matches
		if len(expect.Videos) > 0 {
			if resp.Videos[0].ID != expect.Videos[0].ID {
				t.Fatalf("expected first video id %s, got %s", expect.Videos[0].ID, resp.Videos[0].ID)

			}
		}
	}

}

func TestSearchVideo_ErrorStatusMapping(t *testing.T) {
	cases := []struct {
		name         string
		status       int
		code         string
		errorMessage string
		want         error
	}{
		{name: "bad_request", status: 400, code: "QUERY_PARSE_ERROR", errorMessage: "query parse error", want: ErrRespQueryParse},
		{name: "internal", status: 500, code: "INTERNAL_SERVER_ERROR", errorMessage: "please retry later", want: ErrRespInternal},
		{name: "maint", status: 503, code: "MAINTENANCE", errorMessage: "please retry later.", want: ErrRespMaintainance},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				meta := map[string]interface{}{
					"status":       tc.status,
					"errorCode":    tc.code,
					"errorMessage": tc.errorMessage,
				}
				resp := map[string]interface{}{
					"meta": meta,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}))
			defer srv.Close()

			c := NewVideoClient(srv.URL, "niconico-rss-diy/0.1 test")
			ctx := context.Background()
			_, err := c.SearchVideo(ctx, "vocaloid", nil)
			if err == nil {
				t.Fatalf("expected error for status %d", tc.status)
			}
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
			t.Logf("received error(expect: %d): %v", tc.status, err)
		})
	}
}

func TestFetchThumbnailMeta_ReturnsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Fatalf("expected HEAD request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", "6337")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewThumbnailClient("niconico-rss-diy/0.1 test")
	ctx := context.Background()

	meta, err := c.FetchThumbnailMeta(ctx, srv.URL+"/thumb.jpg")
	if err != nil {
		t.Fatalf("FetchThumbnailMeta error: %v", err)
	}

	if meta == nil {
		t.Fatalf("expected non-nil metadata")
	}

	if meta.Type != "image/jpeg" {
		t.Fatalf("expected Content-Type image/jpeg, got %q", meta.Type)
	}

	if meta.Length != 6337 {
		t.Fatalf("expected Content-Length 6337, got %d", meta.Length)
	}

	if meta.URL == "" {
		t.Fatalf("expected non-empty URL in metadata")
	}
}

func TestFetchThumbnailMeta_ErrorStatusMapping(t *testing.T) {
	cases := []struct {
		name   string
		status int
		want   error
	}{
		{name: "internal", status: http.StatusInternalServerError, want: ErrRespInternal},
		{name: "maint", status: http.StatusServiceUnavailable, want: ErrRespMaintainance},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "HEAD" {
					t.Fatalf("expected HEAD request, got %s", r.Method)
				}
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()

			c := NewThumbnailClient("niconico-rss-diy/0.1 test")
			ctx := context.Background()

			_, err := c.FetchThumbnailMeta(ctx, srv.URL+"/thumb.jpg")
			if err == nil {
				t.Fatalf("expected error for status %d", tc.status)
			}
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

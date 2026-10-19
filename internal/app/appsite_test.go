package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// The file is byte for byte what write-config.mjs buildAppSiteAssociation writes on the
// Cloud for the same app ids.
func TestAppSiteAssociation(t *testing.T) {
	ids := []string{"RX9M3639SF.com.entropy88.app", "ABCDE12345.sh.opendesk.app"}
	w := httptest.NewRecorder()
	appSiteAssociationHandler(ids)(w, httptest.NewRequest(http.MethodGet, appSiteAssociationPath, nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type %q", ct)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=3600" {
		t.Fatalf("Cache-Control %q", cc)
	}
	want := `{"applinks":{"details":[{"appIDs":["RX9M3639SF.com.entropy88.app","ABCDE12345.sh.opendesk.app"],` +
		`"components":[{"/":"/inbox/*/?*"},{"/":"/views/*/?*"}]}]},` +
		`"webcredentials":{"apps":["RX9M3639SF.com.entropy88.app","ABCDE12345.sh.opendesk.app"]}}` + "\n"
	if got := w.Body.String(); got != want {
		t.Fatalf("body\n got %s\nwant %s", got, want)
	}
}

// Without APPLE_APP_IDS the path is not found, and the web app's catch-all never
// answers it with a page.
func TestAppSiteAssociationUnset(t *testing.T) {
	r := chi.NewRouter()
	r.Get(appSiteAssociationPath, appSiteAssociationHandler(nil))
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("<!doctype html>")) }))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, appSiteAssociationPath, nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("status %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("Content-Type %q", ct)
	}
}

package app

import (
	"encoding/json"
	"net/http"
)

// appSiteAssociationPath associates iPhone app builds with this host: universal links
// and passkeys (contracts/native.md "Passkeys").
const appSiteAssociationPath = "/.well-known/apple-app-site-association"

// appLinkPaths are the web app paths the iPhone app opens, as APP_LINK_PATHS in
// deploy/cloudflare/pages/app/write-config.mjs.
var appLinkPaths = []string{"/inbox/*/?*", "/views/*/?*"}

type appSiteAssociation struct {
	AppLinks struct {
		Details []appLinkDetail `json:"details"`
	} `json:"applinks"`
	WebCredentials struct {
		Apps []string `json:"apps"`
	} `json:"webcredentials"`
}

type appLinkDetail struct {
	AppIDs     []string            `json:"appIDs"`
	Components []map[string]string `json:"components"`
}

// appSiteAssociationJSON is the file write-config.mjs buildAppSiteAssociation writes for
// the same app ids.
func appSiteAssociationJSON(ids []string) []byte {
	var a appSiteAssociation
	d := appLinkDetail{AppIDs: ids}
	for _, p := range appLinkPaths {
		d.Components = append(d.Components, map[string]string{"/": p})
	}
	a.AppLinks.Details = []appLinkDetail{d}
	a.WebCredentials.Apps = ids
	b, _ := json.Marshal(a)
	return append(b, '\n')
}

// appSiteAssociationHandler serves the association for APPLE_APP_IDS, and not found
// without them so the web app never answers the path.
func appSiteAssociationHandler(ids []string) http.HandlerFunc {
	if len(ids) == 0 {
		return notFound
	}
	body := appSiteAssociationJSON(ids)
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(body)
	}
}

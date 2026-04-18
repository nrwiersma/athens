package download

import (
	"net/http"

	"github.com/gomods/athens/pkg/download/mode"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/log"
)

// PathVersionInfo URL.
const PathVersionInfo = "/{module:.+}/@v/{version}.info"

// InfoHandler implements GET baseURL/module/@v/version.info.
func InfoHandler(dp Protocol, lggr log.Entry, df *mode.DownloadFile) http.Handler {
	const op apierrors.Op = "download.InfoHandler"
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		mod, ver, err := getModuleParams(r, op)
		if err != nil {
			lggr.SystemErr(err)
			w.WriteHeader(apierrors.Kind(err))
			return
		}
		info, err := dp.Info(r.Context(), mod, ver)
		if err != nil {
			severityLevel := apierrors.Expect(err, apierrors.KindNotFound, apierrors.KindRedirect)
			lggr.SystemErr(apierrors.E(op, err, apierrors.M(mod), apierrors.V(ver), severityLevel))
			if apierrors.Kind(err) == apierrors.KindRedirect {
				url, err := getRedirectURL(df.URL(mod), r.URL.Path)
				if err != nil {
					lggr.SystemErr(err)
					w.WriteHeader(apierrors.Kind(err))
					return
				}
				http.Redirect(w, r, url, apierrors.KindRedirect)
				return
			}
			w.WriteHeader(apierrors.Kind(err))
		}

		_, _ = w.Write(info)
	}
	return http.HandlerFunc(f)
}

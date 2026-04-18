package download

import (
	"net/http"

	"github.com/gomods/athens/pkg/download/mode"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/log"
)

// PathVersionModule URL.
const PathVersionModule = "/{module:.+}/@v/{version}.mod"

// ModuleHandler implements GET baseURL/module/@v/version.mod.
func ModuleHandler(dp Protocol, lggr log.Entry, df *mode.DownloadFile) http.Handler {
	const op apierrors.Op = "download.VersionModuleHandler"
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		mod, ver, err := getModuleParams(r, op)
		if err != nil {
			err = apierrors.E(op, apierrors.M(mod), apierrors.V(ver), err)
			lggr.SystemErr(err)
			w.WriteHeader(apierrors.Kind(err))
			return
		}
		modBts, err := dp.GoMod(r.Context(), mod, ver)
		if err != nil {
			severityLevel := apierrors.Expect(err, apierrors.KindNotFound, apierrors.KindRedirect)
			err = apierrors.E(op, err, severityLevel)
			lggr.SystemErr(err)
			if apierrors.Kind(err) == apierrors.KindRedirect {
				url, err := getRedirectURL(df.URL(mod), r.URL.Path)
				if err != nil {
					err = apierrors.E(op, apierrors.M(mod), apierrors.V(ver), err)
					lggr.SystemErr(err)
					w.WriteHeader(apierrors.Kind(err))
					return
				}
				http.Redirect(w, r, url, apierrors.KindRedirect)
				return
			}
			w.WriteHeader(apierrors.Kind(err))
			return
		}

		_, _ = w.Write(modBts)
	}
	return http.HandlerFunc(f)
}

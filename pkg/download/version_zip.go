package download

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gomods/athens/pkg/download/mode"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/log"
)

// PathVersionZip URL.
const PathVersionZip = "/{module:.+}/@v/{version}.zip"

// ZipHandler implements GET baseURL/module/@v/version.zip.
func ZipHandler(dp Protocol, lggr log.Entry, df *mode.DownloadFile) http.Handler {
	const op apierrors.Op = "download.ZipHandler"
	f := func(w http.ResponseWriter, r *http.Request) {
		mod, ver, err := getModuleParams(r, op)
		if err != nil {
			lggr.SystemErr(err)
			w.WriteHeader(apierrors.Kind(err))
			return
		}
		zip, err := dp.Zip(r.Context(), mod, ver)
		if err != nil {
			severityLevel := apierrors.Expect(err, apierrors.KindNotFound, apierrors.KindRedirect)
			err = apierrors.E(op, err, severityLevel)
			lggr.SystemErr(err)
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
			return
		}
		defer func() { _ = zip.Close() }()

		w.Header().Set("Content-Type", "application/zip")
		size := zip.Size()
		if size > 0 {
			w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		}
		if r.Method == http.MethodHead {
			return
		}
		_, err = io.Copy(w, zip)
		if err != nil {
			lggr.SystemErr(apierrors.E(op, apierrors.M(mod), apierrors.V(ver), err))
		}
	}
	return http.HandlerFunc(f)
}

package download

import (
	"net/http"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/paths"
)

func getModuleParams(r *http.Request, op apierrors.Op) (mod, ver string, err error) {
	params, err := paths.GetAllParams(r)
	if err != nil {
		return "", "", apierrors.E(op, err, apierrors.KindBadRequest)
	}

	return params.Module, params.Version, nil
}

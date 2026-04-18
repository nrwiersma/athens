package fs

import (
	"context"
	"os"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/spf13/afero"
)

func (s *storageImpl) Exists(ctx context.Context, module, version string) (bool, error) {
	const op apierrors.Op = "fs.Exists"
	_, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	versionedPath := s.versionLocation(module, version)

	files, err := afero.ReadDir(s.filesystem, versionedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, apierrors.E(op, apierrors.M(module), apierrors.V(version), err)
	}

	return len(files) == 3, nil
}

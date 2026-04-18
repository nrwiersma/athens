package fs

import (
	"context"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
)

// Delete removes a specific version of a module.
func (s *storageImpl) Delete(ctx context.Context, module, version string) error {
	const op apierrors.Op = "fs.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	versionedPath := s.versionLocation(module, version)
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}
	return s.filesystem.RemoveAll(versionedPath)
}

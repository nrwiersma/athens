package gcp

import (
	"context"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	modupl "github.com/gomods/athens/pkg/storage/module"
)

// Delete implements the (./pkg/storage).Deleter interface and
// removes a version of a module from storage. Returning ErrNotFound
// if the version does not exist.
func (s *Storage) Delete(ctx context.Context, module, version string) error {
	const op apierrors.Op = "gcp.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}
	del := func(ctx context.Context, path string) error {
		return s.bucket.Object(path).Delete(ctx)
	}
	err = modupl.Delete(ctx, module, version, del, s.timeout)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	return nil
}

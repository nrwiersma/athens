package azureblob

import (
	"context"

	"github.com/gomods/athens/pkg/config"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
)

// Exists implements the (./pkg/storage).Checker interface
// returning true if the module at version exists in storage.
func (s *Storage) Exists(ctx context.Context, module, version string) (bool, error) {
	const op apierrors.Op = "azureblob.Exists"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	px := config.PackageVersionedName(module, version, "")
	paths, err := s.client.ListBlobs(ctx, px)
	if err != nil {
		return false, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	var count int
	for _, p := range paths {
		// sane assumption: no duplicate keys.
		switch p {
		case config.PackageVersionedName(module, version, "info"):
			count++
		case config.PackageVersionedName(module, version, "mod"):
			count++
		case config.PackageVersionedName(module, version, "zip"):
			count++
		}
	}
	return count == 3, nil
}

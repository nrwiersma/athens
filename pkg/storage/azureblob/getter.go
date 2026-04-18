package azureblob

import (
	"context"
	"fmt"
	"io"

	"github.com/gomods/athens/pkg/config"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gomods/athens/pkg/storage"
)

// Info implements the (./pkg/storage).Getter interface.
func (s *Storage) Info(ctx context.Context, module, version string) ([]byte, error) {
	const op apierrors.Op = "azureblob.Info"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}

	infoReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "info"))
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	infoBytes, err := io.ReadAll(infoReader)
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	err = infoReader.Close()
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	return infoBytes, nil
}

// GoMod implements the (./pkg/storage).Getter interface.
func (s *Storage) GoMod(ctx context.Context, module, version string) ([]byte, error) {
	const op apierrors.Op = "azureblob.GoMod"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}

	modReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "mod"))
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	modBytes, err := io.ReadAll(modReader)
	if err != nil {
		return nil, apierrors.E(op, fmt.Errorf("could not get new reader for mod file: %w", err), apierrors.M(module), apierrors.V(version))
	}

	err = modReader.Close()
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	return modBytes, nil
}

// Zip implements the (./pkg/storage).Getter interface.
func (s *Storage) Zip(ctx context.Context, module, version string) (storage.SizeReadCloser, error) {
	const op apierrors.Op = "azureblob.Zip"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}
	zipReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "zip"))
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	return zipReader, nil
}

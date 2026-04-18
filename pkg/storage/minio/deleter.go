package minio

import (
	"context"
	"fmt"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
)

func (s *storageImpl) Delete(ctx context.Context, module, version string) error {
	const op apierrors.Op = "minio.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	if !exists {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}

	versionedPath := s.versionLocation(module, version)

	modPath := fmt.Sprintf("%s/go.mod", versionedPath)
	if err := s.minioClient.RemoveObject(s.bucketName, modPath); err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	zipPath := fmt.Sprintf("%s/source.zip", versionedPath)
	if err := s.minioClient.RemoveObject(s.bucketName, zipPath); err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	infoPath := fmt.Sprintf("%s/%s.info", versionedPath, version)
	err = s.minioClient.RemoveObject(s.bucketName, infoPath)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	return nil
}

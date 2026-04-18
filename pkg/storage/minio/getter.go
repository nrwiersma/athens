package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gomods/athens/pkg/storage"
	"github.com/minio/minio-go/v6"
)

func (s *storageImpl) Info(ctx context.Context, module, vsn string) ([]byte, error) {
	const op apierrors.Op = "minio.Info"
	_, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	infoPath := fmt.Sprintf("%s/%s.info", s.versionLocation(module, vsn), vsn)
	infoReader, err := s.minioClient.GetObject(s.bucketName, infoPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, apierrors.E(op, err)
	}
	defer func() { _ = infoReader.Close() }()
	info, err := io.ReadAll(infoReader)
	if err != nil {
		return nil, transformNotFoundErr(op, module, vsn, err)
	}

	return info, nil
}

func (s *storageImpl) GoMod(ctx context.Context, module, vsn string) ([]byte, error) {
	const op apierrors.Op = "minio.GoMod"
	_, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	modPath := fmt.Sprintf("%s/go.mod", s.versionLocation(module, vsn))
	modReader, err := s.minioClient.GetObject(s.bucketName, modPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, apierrors.E(op, err)
	}
	defer func() { _ = modReader.Close() }()
	mod, err := io.ReadAll(modReader)
	if err != nil {
		return nil, transformNotFoundErr(op, module, vsn, err)
	}

	return mod, nil
}

func (s *storageImpl) Zip(ctx context.Context, module, vsn string) (storage.SizeReadCloser, error) {
	const op apierrors.Op = "minio.Zip"
	_, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	zipPath := fmt.Sprintf("%s/source.zip", s.versionLocation(module, vsn))
	_, err := s.minioClient.StatObject(s.bucketName, zipPath, minio.StatObjectOptions{})
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.KindNotFound, apierrors.M(module), apierrors.V(vsn))
	}

	zipReader, err := s.minioClient.GetObject(s.bucketName, zipPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, apierrors.E(op, err)
	}
	oi, err := zipReader.Stat()
	if err != nil {
		_ = zipReader.Close()
		return nil, apierrors.E(op, err)
	}
	return storage.NewSizer(zipReader, oi.Size), nil
}

func transformNotFoundErr(op apierrors.Op, module, version string, err error) error {
	if respErr, ok := errors.AsType[minio.ErrorResponse](err); ok {
		if respErr.StatusCode == http.StatusNotFound {
			return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
		}
	}
	return err
}

package s3

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gomods/athens/pkg/config"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gomods/athens/pkg/storage"
)

// Info implements the (./pkg/storage).Getter interface.
func (s *Storage) Info(ctx context.Context, module, version string) ([]byte, error) {
	const op apierrors.Op = "s3.Info"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	infoReader, err := s.open(ctx, config.PackageVersionedName(module, version, "info"))
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
		}
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	defer func() { _ = infoReader.Close() }()

	infoBytes, err := io.ReadAll(infoReader)
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	return infoBytes, nil
}

// GoMod implements the (./pkg/storage).Getter interface.
func (s *Storage) GoMod(ctx context.Context, module, version string) ([]byte, error) {
	const op apierrors.Op = "s3.GoMod"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	modReader, err := s.open(ctx, config.PackageVersionedName(module, version, "mod"))
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
		}
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	defer func() { _ = modReader.Close() }()

	modBytes, err := io.ReadAll(modReader)
	if err != nil {
		return nil, apierrors.E(op, fmt.Errorf("could not get new reader for mod file: %w", err), apierrors.M(module), apierrors.V(version))
	}

	return modBytes, nil
}

// Zip implements the (./pkg/storage).Getter interface.
func (s *Storage) Zip(ctx context.Context, module, version string) (storage.SizeReadCloser, error) {
	const op apierrors.Op = "s3.Zip"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	zipReader, err := s.open(ctx, config.PackageVersionedName(module, version, "zip"))
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return nil, apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
		}
		return nil, apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	return zipReader, nil
}

func (s *Storage) open(ctx context.Context, path string) (storage.SizeReadCloser, error) {
	const op apierrors.Op = "s3.open"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	getParams := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	goo, err := s.s3API.GetObject(ctx, getParams)
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return nil, apierrors.E(op, apierrors.KindNotFound)
		}
		return nil, apierrors.E(op, err)
	}
	var size int64
	if goo.ContentLength != nil {
		size = *goo.ContentLength
	}
	return storage.NewSizer(goo.Body, size), nil
}

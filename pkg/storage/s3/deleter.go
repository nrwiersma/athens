package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	modupl "github.com/gomods/athens/pkg/storage/module"
)

// Delete implements the (./pkg/storage).Deleter interface and
// removes a version of a module from storage. Returning ErrNotFound
// if the version does not exist.
func (s *Storage) Delete(ctx context.Context, module, version string) error {
	const op apierrors.Op = "s3.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}
	if !exists {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}

	return modupl.Delete(ctx, module, version, s.remove, s.timeout)
}

func (s *Storage) remove(ctx context.Context, path string) error {
	const op apierrors.Op = "s3.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	delParams := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	if _, err := s.s3API.DeleteObject(ctx, delParams); err != nil {
		return apierrors.E(op, err)
	}

	return nil
}

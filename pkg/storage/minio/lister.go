package minio

import (
	"context"
	"fmt"
	"sort"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
)

func (s *storageImpl) List(ctx context.Context, module string) ([]string, error) {
	const op apierrors.Op = "minio.List"
	_, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	doneCh := make(chan struct{})
	defer close(doneCh)
	searchPrefix := module + "/"
	objectCh, err := s.minioCore.ListObjectsV2(s.bucketName, searchPrefix, "", false, "", 0, "")
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(module))
	}
	var ret []string
	for _, object := range objectCh.Contents {
		if object.Err != nil {
			return nil, apierrors.E(op, object.Err, apierrors.M(module))
		}

		key, _, ver := extractKey(object.Key)
		goModKey := fmt.Sprintf("%s/go.mod", s.versionLocation(module, ver))
		if goModKey == key {
			ret = append(ret, ver)
		}
	}
	sort.Strings(ret)
	return ret, nil
}

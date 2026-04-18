package azureblob

import (
	"bytes"
	"context"
	"io"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	moduploader "github.com/gomods/athens/pkg/storage/module"
)

// Save implements the (./pkg/storage).Saver interface.
func (s *Storage) Save(ctx context.Context, module, version string, mod []byte, zip io.Reader, zipMD5, info []byte) error {
	const op apierrors.Op = "azureblob.Save"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	err := moduploader.Upload(ctx, module, version, bytes.NewReader(info), bytes.NewReader(mod), zip, s.client.UploadWithContext, s.timeout)
	if err != nil {
		return apierrors.E(op, err, apierrors.M(module), apierrors.V(version))
	}

	return nil
}

package mongo

import (
	"context"
	"errors"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gomods/athens/pkg/storage"
	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// List lists all versions of a module.
func (s *ModuleStore) List(ctx context.Context, moduleName string) ([]string, error) {
	const op apierrors.Op = "mongo.List"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	c := s.client.Database(s.db).Collection(s.coll)
	projection := bson.M{"version": 1, "_id": 0}
	query := bson.M{"module": moduleName}
	tctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cursor, err := c.Find(tctx, query, options.Find().SetProjection(projection))
	if err != nil {
		return nil, apierrors.E(op, err, apierrors.M(moduleName))
	}
	result := make([]storage.Module, 0)
	var errs error
	for cursor.Next(ctx) {
		var module storage.Module
		if err = cursor.Decode(&module); err != nil {
			kind := apierrors.KindUnexpected
			if errors.Is(err, mongo.ErrNoDocuments) {
				kind = apierrors.KindNotFound
			}
			errs = multierror.Append(errs, apierrors.E(op, err, kind))
		} else {
			result = append(result, module)
		}
	}

	versions := make([]string, len(result))
	for i, r := range result {
		versions[i] = r.Version
	}

	return versions, nil
}

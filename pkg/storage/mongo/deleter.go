package mongo

import (
	"context"
	"errors"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Delete removes a specific version of a module.
func (s *ModuleStore) Delete(ctx context.Context, module, version string) error {
	const op apierrors.Op = "mongo.Delete"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}
	if !exists {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), apierrors.KindNotFound)
	}

	db := s.client.Database(s.db)
	c := db.Collection(s.coll)
	bucket := db.GridFSBucket()

	filter := bson.D{bson.E{Key: "filename", Value: s.gridFileName(module, version)}}

	cursor, err := bucket.Find(ctx, filter)
	if err != nil {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), err)
	}

	var x bson.D
	for cursor.Next(ctx) {
		_ = cursor.Decode(&x)
	}
	b, err := bson.Marshal(x)
	if err != nil {
		return apierrors.E(op, apierrors.M(module), apierrors.V(version), err)
	}

	if err = bucket.Delete(ctx, bson.Raw(b).Lookup("_id").ObjectID()); err != nil {
		kind := apierrors.KindUnexpected
		if errors.Is(err, mongo.ErrFileNotFound) {
			kind = apierrors.KindNotFound
		}
		return apierrors.E(op, err, kind, apierrors.M(module), apierrors.V(version))
	}

	tctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	_, err = c.DeleteOne(tctx, bson.M{"module": module, "version": version})
	if err != nil {
		return apierrors.E(op, err, apierrors.KindNotFound, apierrors.M(module), apierrors.V(version))
	}
	return nil
}

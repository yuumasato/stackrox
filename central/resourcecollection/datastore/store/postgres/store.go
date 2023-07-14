// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "collections"
	storeName = "ResourceCollection"

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.CollectionsSchema
	targetResource = resources.WorkflowAdministration
)

type storeType = storage.ResourceCollection

// Store is the interface to interact with the storage for storage.ResourceCollection
type Store interface {
	Upsert(ctx context.Context, obj *storeType) error
	UpsertMany(ctx context.Context, objs []*storeType) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storeType, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storeType, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storeType, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storeType) error) error
}

type storeImpl struct {
	*pgSearch.GenericStore[storeType, *storeType]
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		GenericStore: pgSearch.NewGenericStore[storeType, *storeType](
			db,
			schema,
			pkGetter,
			insertIntoCollections,
			copyFromCollections,
			metricsSetAcquireDBConnDuration,
			metricsSetPostgresOperationDurationTime,
			pgSearch.GloballyScopedUpsertChecker[storeType, *storeType](targetResource),
			targetResource,
		),
	}
}

// region Helper functions

func pkGetter(obj *storeType) string {
	return obj.GetId()
}

func metricsSetPostgresOperationDurationTime(start time.Time, op ops.Op) {
	metrics.SetPostgresOperationDurationTime(start, op, storeName)
}

func metricsSetAcquireDBConnDuration(start time.Time, op ops.Op) {
	metrics.SetAcquireDBConnDuration(start, op, storeName)
}

func insertIntoCollections(ctx context.Context, batch *pgx.Batch, obj *storage.ResourceCollection) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetId(),
		obj.GetName(),
		obj.GetCreatedBy().GetName(),
		obj.GetUpdatedBy().GetName(),
		serialized,
	}

	finalStr := "INSERT INTO collections (Id, Name, CreatedBy_Name, UpdatedBy_Name, serialized) VALUES($1, $2, $3, $4, $5) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, CreatedBy_Name = EXCLUDED.CreatedBy_Name, UpdatedBy_Name = EXCLUDED.UpdatedBy_Name, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetEmbeddedCollections() {
		if err := insertIntoCollectionsEmbeddedCollections(ctx, batch, child, obj.GetId(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from collections_embedded_collections where collections_Id = $1 AND idx >= $2"
	batch.Queue(query, obj.GetId(), len(obj.GetEmbeddedCollections()))
	return nil
}

func insertIntoCollectionsEmbeddedCollections(_ context.Context, batch *pgx.Batch, obj *storage.ResourceCollection_EmbeddedResourceCollection, collectionID string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		collectionID,
		idx,
		obj.GetId(),
	}

	finalStr := "INSERT INTO collections_embedded_collections (collections_Id, idx, Id) VALUES($1, $2, $3) ON CONFLICT(collections_Id, idx) DO UPDATE SET collections_Id = EXCLUDED.collections_Id, idx = EXCLUDED.idx, Id = EXCLUDED.Id"
	batch.Queue(finalStr, values...)

	return nil
}

func copyFromCollections(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, objs ...*storage.ResourceCollection) error {
	inputRows := make([][]interface{}, 0, batchSize)

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	deletes := make([]string, 0, batchSize)

	copyCols := []string{
		"id",
		"name",
		"createdby_name",
		"updatedby_name",
		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{
			obj.GetId(),
			obj.GetName(),
			obj.GetCreatedBy().GetName(),
			obj.GetUpdatedBy().GetName(),
			serialized,
		})

		// Add the ID to be deleted.
		deletes = append(deletes, obj.GetId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = deletes[:0]

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"collections"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err := copyFromCollectionsEmbeddedCollections(ctx, s, tx, obj.GetId(), obj.GetEmbeddedCollections()...); err != nil {
			return err
		}
	}

	return nil
}

func copyFromCollectionsEmbeddedCollections(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, collectionID string, objs ...*storage.ResourceCollection_EmbeddedResourceCollection) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"collections_id",
		"idx",
		"id",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			collectionID,
			idx,
			obj.GetId(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"collections_embedded_collections"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

// endregion Helper functions

// region Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableCollections(ctx, db)
}

func dropTableCollections(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS collections CASCADE")
	dropTableCollectionsEmbeddedCollections(ctx, db)

}

func dropTableCollectionsEmbeddedCollections(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS collections_embedded_collections CASCADE")

}

// endregion Used for testing

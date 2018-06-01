package tsm1

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/pkg/metrics"
	"github.com/influxdata/influxdb/query"
	"github.com/influxdata/influxdb/tsdb"
	"github.com/influxdata/influxql"
)

type cursorIterator struct {
	e   *Engine
	key []byte

	asc struct {
		Float    *floatAscendingBatchCursor
		Integer  *integerAscendingBatchCursor
		Unsigned *unsignedAscendingBatchCursor
		Boolean  *booleanAscendingBatchCursor
		String   *stringAscendingBatchCursor
	}

	desc struct {
		Float    *floatDescendingBatchCursor
		Integer  *integerDescendingBatchCursor
		Unsigned *unsignedDescendingBatchCursor
		Boolean  *booleanDescendingBatchCursor
		String   *stringDescendingBatchCursor
	}
}

func (q *cursorIterator) Next(ctx context.Context, r *tsdb.CursorRequest) (tsdb.Cursor, error) {
	// Look up fields for measurement.
	mf := q.e.fieldset.Fields(r.Name)
	if mf == nil {
		return nil, nil
	}

	// Find individual field.
	f := mf.Field(r.Field)
	if f == nil {
		// field doesn't exist for this measurement
		return nil, nil
	}

	if grp := metrics.GroupFromContext(ctx); grp != nil {
		grp.GetCounter(numberOfRefCursorsCounter).Add(1)
	}

	var opt query.IteratorOptions
	opt.Ascending = r.Ascending
	opt.StartTime = r.StartTime
	opt.EndTime = r.EndTime

	// Return appropriate cursor based on type.
	switch f.Type {
	case influxql.Float:
		return q.buildFloatBatchCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Integer:
		return q.buildIntegerBatchCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Unsigned:
		return q.buildUnsignedBatchCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.String:
		return q.buildStringBatchCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Boolean:
		return q.buildStringBatchCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	default:
		panic(fmt.Sprintf("unreachable: %T", f.Type))
	}
}

func (q *cursorIterator) seriesFieldKeyBytes(name []byte, tags models.Tags, field string) []byte {
	q.key = models.AppendMakeKey(q.key[:0], name, tags)
	q.key = append(q.key, keyFieldSeparatorBytes...)
	q.key = append(q.key, field...)
	return q.key
}

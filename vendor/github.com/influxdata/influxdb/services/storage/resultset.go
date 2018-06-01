package storage

import (
	"context"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/tsdb"
)

type readRequest struct {
	ctx        context.Context
	start, end int64
	asc        bool
	limit      int64
	aggregate  *Aggregate
}

type ResultSet struct {
	req readRequest
	cur seriesCursor
	row seriesRow
	mb  *multiShardBatchCursors
}

func (r *ResultSet) Close() {
	r.row.query = nil
	r.cur.Close()
}

func (r *ResultSet) Next() bool {
	row := r.cur.Next()
	if row == nil {
		return false
	}

	r.row = *row

	return true
}

func (r *ResultSet) Cursor() tsdb.Cursor {
	cur := r.mb.createCursor(r.row)
	if r.req.aggregate != nil {
		cur = newAggregateBatchCursor(r.req.ctx, r.req.aggregate, cur)
	}
	return cur
}

func (r *ResultSet) Tags() models.Tags {
	return r.row.tags
}

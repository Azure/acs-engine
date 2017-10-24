// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: file_store.gen.go.tmpl

package tsm1

// ReadFloatBlock reads the next block as a set of float values.
func (c *KeyCursor) ReadFloatBlock(buf *[]FloatValue) ([]FloatValue, error) {
	// No matching blocks to decode
	if len(c.current) == 0 {
		return nil, nil
	}

	// First block is the oldest block containing the points we're searching for.
	first := c.current[0]
	*buf = (*buf)[:0]
	values, err := first.r.ReadFloatBlockAt(&first.entry, buf)
	if err != nil {
		return nil, err
	}

	// Remove values we already read
	values = FloatValues(values).Exclude(first.readMin, first.readMax)

	// Remove any tombstones
	tombstones := first.r.TombstoneRange(c.key)
	values = c.filterFloatValues(tombstones, values)

	// Check we have remaining values.
	if len(values) == 0 {
		return nil, nil
	}

	// Only one block with this key and time range so return it
	if len(c.current) == 1 {
		if len(values) > 0 {
			first.markRead(values[0].UnixNano(), values[len(values)-1].UnixNano())
		}
		return values, nil
	}

	// Use the current block time range as our overlapping window
	minT, maxT := first.readMin, first.readMax
	if len(values) > 0 {
		minT, maxT = values[0].UnixNano(), values[len(values)-1].UnixNano()
	}
	if c.ascending {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the min time range to ensure values are returned in ascending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MinTime < minT && !cur.read() {
				minT = cur.entry.MinTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MaxTime > maxT {
					maxT = cur.entry.MaxTime
				}
				values = FloatValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)
			var a []FloatValue
			v, err := cur.r.ReadFloatBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterFloatValues(tombstones, v)

			// Remove values we already read
			v = FloatValues(v).Exclude(cur.readMin, cur.readMax)

			if len(v) > 0 {
				// Only use values in the overlapping window
				v = FloatValues(v).Include(minT, maxT)

				// Merge the remaing values with the existing
				values = FloatValues(values).Merge(v)
			}
			cur.markRead(minT, maxT)
		}

	} else {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the max time range to ensure values are returned in descending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MaxTime > maxT && !cur.read() {
				maxT = cur.entry.MaxTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MinTime < minT {
					minT = cur.entry.MinTime
				}
				values = FloatValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)

			var a []FloatValue
			v, err := cur.r.ReadFloatBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterFloatValues(tombstones, v)

			// Remove values we already read
			v = FloatValues(v).Exclude(cur.readMin, cur.readMax)

			// If the block we decoded should have all of it's values included, mark it as read so we
			// don't use it again.
			if len(v) > 0 {
				v = FloatValues(v).Include(minT, maxT)
				// Merge the remaing values with the existing
				values = FloatValues(v).Merge(values)
			}
			cur.markRead(minT, maxT)
		}
	}

	first.markRead(minT, maxT)

	return values, err
}

// ReadIntegerBlock reads the next block as a set of integer values.
func (c *KeyCursor) ReadIntegerBlock(buf *[]IntegerValue) ([]IntegerValue, error) {
	// No matching blocks to decode
	if len(c.current) == 0 {
		return nil, nil
	}

	// First block is the oldest block containing the points we're searching for.
	first := c.current[0]
	*buf = (*buf)[:0]
	values, err := first.r.ReadIntegerBlockAt(&first.entry, buf)
	if err != nil {
		return nil, err
	}

	// Remove values we already read
	values = IntegerValues(values).Exclude(first.readMin, first.readMax)

	// Remove any tombstones
	tombstones := first.r.TombstoneRange(c.key)
	values = c.filterIntegerValues(tombstones, values)

	// Check we have remaining values.
	if len(values) == 0 {
		return nil, nil
	}

	// Only one block with this key and time range so return it
	if len(c.current) == 1 {
		if len(values) > 0 {
			first.markRead(values[0].UnixNano(), values[len(values)-1].UnixNano())
		}
		return values, nil
	}

	// Use the current block time range as our overlapping window
	minT, maxT := first.readMin, first.readMax
	if len(values) > 0 {
		minT, maxT = values[0].UnixNano(), values[len(values)-1].UnixNano()
	}
	if c.ascending {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the min time range to ensure values are returned in ascending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MinTime < minT && !cur.read() {
				minT = cur.entry.MinTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MaxTime > maxT {
					maxT = cur.entry.MaxTime
				}
				values = IntegerValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)
			var a []IntegerValue
			v, err := cur.r.ReadIntegerBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterIntegerValues(tombstones, v)

			// Remove values we already read
			v = IntegerValues(v).Exclude(cur.readMin, cur.readMax)

			if len(v) > 0 {
				// Only use values in the overlapping window
				v = IntegerValues(v).Include(minT, maxT)

				// Merge the remaing values with the existing
				values = IntegerValues(values).Merge(v)
			}
			cur.markRead(minT, maxT)
		}

	} else {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the max time range to ensure values are returned in descending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MaxTime > maxT && !cur.read() {
				maxT = cur.entry.MaxTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MinTime < minT {
					minT = cur.entry.MinTime
				}
				values = IntegerValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)

			var a []IntegerValue
			v, err := cur.r.ReadIntegerBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterIntegerValues(tombstones, v)

			// Remove values we already read
			v = IntegerValues(v).Exclude(cur.readMin, cur.readMax)

			// If the block we decoded should have all of it's values included, mark it as read so we
			// don't use it again.
			if len(v) > 0 {
				v = IntegerValues(v).Include(minT, maxT)
				// Merge the remaing values with the existing
				values = IntegerValues(v).Merge(values)
			}
			cur.markRead(minT, maxT)
		}
	}

	first.markRead(minT, maxT)

	return values, err
}

// ReadUnsignedBlock reads the next block as a set of unsigned values.
func (c *KeyCursor) ReadUnsignedBlock(buf *[]UnsignedValue) ([]UnsignedValue, error) {
	// No matching blocks to decode
	if len(c.current) == 0 {
		return nil, nil
	}

	// First block is the oldest block containing the points we're searching for.
	first := c.current[0]
	*buf = (*buf)[:0]
	values, err := first.r.ReadUnsignedBlockAt(&first.entry, buf)
	if err != nil {
		return nil, err
	}

	// Remove values we already read
	values = UnsignedValues(values).Exclude(first.readMin, first.readMax)

	// Remove any tombstones
	tombstones := first.r.TombstoneRange(c.key)
	values = c.filterUnsignedValues(tombstones, values)

	// Check we have remaining values.
	if len(values) == 0 {
		return nil, nil
	}

	// Only one block with this key and time range so return it
	if len(c.current) == 1 {
		if len(values) > 0 {
			first.markRead(values[0].UnixNano(), values[len(values)-1].UnixNano())
		}
		return values, nil
	}

	// Use the current block time range as our overlapping window
	minT, maxT := first.readMin, first.readMax
	if len(values) > 0 {
		minT, maxT = values[0].UnixNano(), values[len(values)-1].UnixNano()
	}
	if c.ascending {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the min time range to ensure values are returned in ascending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MinTime < minT && !cur.read() {
				minT = cur.entry.MinTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MaxTime > maxT {
					maxT = cur.entry.MaxTime
				}
				values = UnsignedValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)
			var a []UnsignedValue
			v, err := cur.r.ReadUnsignedBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterUnsignedValues(tombstones, v)

			// Remove values we already read
			v = UnsignedValues(v).Exclude(cur.readMin, cur.readMax)

			if len(v) > 0 {
				// Only use values in the overlapping window
				v = UnsignedValues(v).Include(minT, maxT)

				// Merge the remaing values with the existing
				values = UnsignedValues(values).Merge(v)
			}
			cur.markRead(minT, maxT)
		}

	} else {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the max time range to ensure values are returned in descending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MaxTime > maxT && !cur.read() {
				maxT = cur.entry.MaxTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MinTime < minT {
					minT = cur.entry.MinTime
				}
				values = UnsignedValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)

			var a []UnsignedValue
			v, err := cur.r.ReadUnsignedBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterUnsignedValues(tombstones, v)

			// Remove values we already read
			v = UnsignedValues(v).Exclude(cur.readMin, cur.readMax)

			// If the block we decoded should have all of it's values included, mark it as read so we
			// don't use it again.
			if len(v) > 0 {
				v = UnsignedValues(v).Include(minT, maxT)
				// Merge the remaing values with the existing
				values = UnsignedValues(v).Merge(values)
			}
			cur.markRead(minT, maxT)
		}
	}

	first.markRead(minT, maxT)

	return values, err
}

// ReadStringBlock reads the next block as a set of string values.
func (c *KeyCursor) ReadStringBlock(buf *[]StringValue) ([]StringValue, error) {
	// No matching blocks to decode
	if len(c.current) == 0 {
		return nil, nil
	}

	// First block is the oldest block containing the points we're searching for.
	first := c.current[0]
	*buf = (*buf)[:0]
	values, err := first.r.ReadStringBlockAt(&first.entry, buf)
	if err != nil {
		return nil, err
	}

	// Remove values we already read
	values = StringValues(values).Exclude(first.readMin, first.readMax)

	// Remove any tombstones
	tombstones := first.r.TombstoneRange(c.key)
	values = c.filterStringValues(tombstones, values)

	// Check we have remaining values.
	if len(values) == 0 {
		return nil, nil
	}

	// Only one block with this key and time range so return it
	if len(c.current) == 1 {
		if len(values) > 0 {
			first.markRead(values[0].UnixNano(), values[len(values)-1].UnixNano())
		}
		return values, nil
	}

	// Use the current block time range as our overlapping window
	minT, maxT := first.readMin, first.readMax
	if len(values) > 0 {
		minT, maxT = values[0].UnixNano(), values[len(values)-1].UnixNano()
	}
	if c.ascending {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the min time range to ensure values are returned in ascending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MinTime < minT && !cur.read() {
				minT = cur.entry.MinTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MaxTime > maxT {
					maxT = cur.entry.MaxTime
				}
				values = StringValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)
			var a []StringValue
			v, err := cur.r.ReadStringBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterStringValues(tombstones, v)

			// Remove values we already read
			v = StringValues(v).Exclude(cur.readMin, cur.readMax)

			if len(v) > 0 {
				// Only use values in the overlapping window
				v = StringValues(v).Include(minT, maxT)

				// Merge the remaing values with the existing
				values = StringValues(values).Merge(v)
			}
			cur.markRead(minT, maxT)
		}

	} else {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the max time range to ensure values are returned in descending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MaxTime > maxT && !cur.read() {
				maxT = cur.entry.MaxTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MinTime < minT {
					minT = cur.entry.MinTime
				}
				values = StringValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)

			var a []StringValue
			v, err := cur.r.ReadStringBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterStringValues(tombstones, v)

			// Remove values we already read
			v = StringValues(v).Exclude(cur.readMin, cur.readMax)

			// If the block we decoded should have all of it's values included, mark it as read so we
			// don't use it again.
			if len(v) > 0 {
				v = StringValues(v).Include(minT, maxT)
				// Merge the remaing values with the existing
				values = StringValues(v).Merge(values)
			}
			cur.markRead(minT, maxT)
		}
	}

	first.markRead(minT, maxT)

	return values, err
}

// ReadBooleanBlock reads the next block as a set of boolean values.
func (c *KeyCursor) ReadBooleanBlock(buf *[]BooleanValue) ([]BooleanValue, error) {
	// No matching blocks to decode
	if len(c.current) == 0 {
		return nil, nil
	}

	// First block is the oldest block containing the points we're searching for.
	first := c.current[0]
	*buf = (*buf)[:0]
	values, err := first.r.ReadBooleanBlockAt(&first.entry, buf)
	if err != nil {
		return nil, err
	}

	// Remove values we already read
	values = BooleanValues(values).Exclude(first.readMin, first.readMax)

	// Remove any tombstones
	tombstones := first.r.TombstoneRange(c.key)
	values = c.filterBooleanValues(tombstones, values)

	// Check we have remaining values.
	if len(values) == 0 {
		return nil, nil
	}

	// Only one block with this key and time range so return it
	if len(c.current) == 1 {
		if len(values) > 0 {
			first.markRead(values[0].UnixNano(), values[len(values)-1].UnixNano())
		}
		return values, nil
	}

	// Use the current block time range as our overlapping window
	minT, maxT := first.readMin, first.readMax
	if len(values) > 0 {
		minT, maxT = values[0].UnixNano(), values[len(values)-1].UnixNano()
	}
	if c.ascending {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the min time range to ensure values are returned in ascending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MinTime < minT && !cur.read() {
				minT = cur.entry.MinTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MaxTime > maxT {
					maxT = cur.entry.MaxTime
				}
				values = BooleanValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)
			var a []BooleanValue
			v, err := cur.r.ReadBooleanBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterBooleanValues(tombstones, v)

			// Remove values we already read
			v = BooleanValues(v).Exclude(cur.readMin, cur.readMax)

			if len(v) > 0 {
				// Only use values in the overlapping window
				v = BooleanValues(v).Include(minT, maxT)

				// Merge the remaing values with the existing
				values = BooleanValues(values).Merge(v)
			}
			cur.markRead(minT, maxT)
		}

	} else {
		// Blocks are ordered by generation, we may have values in the past in later blocks, if so,
		// expand the window to include the max time range to ensure values are returned in descending
		// order
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.MaxTime > maxT && !cur.read() {
				maxT = cur.entry.MaxTime
			}
		}

		// Find first block that overlaps our window
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			if cur.entry.OverlapsTimeRange(minT, maxT) && !cur.read() {
				// Shrink our window so it's the intersection of the first overlapping block and the
				// first block.  We do this to minimize the region that overlaps and needs to
				// be merged.
				if cur.entry.MinTime < minT {
					minT = cur.entry.MinTime
				}
				values = BooleanValues(values).Include(minT, maxT)
				break
			}
		}

		// Search the remaining blocks that overlap our window and append their values so we can
		// merge them.
		for i := 1; i < len(c.current); i++ {
			cur := c.current[i]
			// Skip this block if it doesn't contain points we looking for or they have already been read
			if !cur.entry.OverlapsTimeRange(minT, maxT) || cur.read() {
				cur.markRead(minT, maxT)
				continue
			}

			tombstones := cur.r.TombstoneRange(c.key)

			var a []BooleanValue
			v, err := cur.r.ReadBooleanBlockAt(&cur.entry, &a)
			if err != nil {
				return nil, err
			}
			// Remove any tombstoned values
			v = c.filterBooleanValues(tombstones, v)

			// Remove values we already read
			v = BooleanValues(v).Exclude(cur.readMin, cur.readMax)

			// If the block we decoded should have all of it's values included, mark it as read so we
			// don't use it again.
			if len(v) > 0 {
				v = BooleanValues(v).Include(minT, maxT)
				// Merge the remaing values with the existing
				values = BooleanValues(v).Merge(values)
			}
			cur.markRead(minT, maxT)
		}
	}

	first.markRead(minT, maxT)

	return values, err
}

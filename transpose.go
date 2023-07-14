package datatable

type _cell struct {
	val   any
	index int
}

type inject struct {
	index int
	cells Cells
}

func (rows *Rows) Transpose() {
	var cellsToInject = make([][]*inject, len(*rows))

	for i, row := range *rows {
		var colSpannedCells Cells
		for i, c := range row {
			c.Val = &_cell{index: i, val: c.Val}
			if c.Colspan > 1 {
				colSpannedCells = append(colSpannedCells, c)
			}
		}

		cellsToInject[i] = make([]*inject, len(colSpannedCells))

		for j, cel := range colSpannedCells {
			cellsToInject[i][j] = &inject{
				index: cel.Val.(*_cell).index,
				cells: make(Cells, cel.Colspan-1),
			}
		}
	}

	var r int
	for _, row := range cellsToInject {
		if len(row) > 0 {
			var injectIndex, injectCount int
			for _, col := range row {
				injectIndex = col.index + injectCount + 1
				(*rows)[r] = (*rows)[r].splice(injectIndex, 0, col.cells)
				injectCount += len(col.cells)
			}
		}
		r++
	}

	*rows = (*rows).transpose()
}

func (rows Rows) transpose() (transposed Rows) {
	transposed = make(Rows, len(rows[0]))
	for i := range rows[0] {
		if transposed[i] == nil {
			transposed[i] = Cells{}
		}

		for j := range rows {
			cell := rows[j][i]
			if cell != nil {
				cell.Val = cell.Val.(*_cell).val
				cell.Colspan, cell.Rowspan = cell.Rowspan, cell.Colspan
				transposed[i] = append(transposed[i], cell)
			}
		}
	}
	return transposed
}

func (input Cells) splice(start, deleteCount int, item Cells) Cells {
	cpy := make(Cells, len(input))
	copy(cpy, input)
	if start > len(cpy) {
		return append(cpy, item...)
	}
	ret := append(cpy[:start], item...)
	if start+deleteCount > len(cpy) {
		return ret
	}
	return append(ret, input[start+deleteCount:]...)
}

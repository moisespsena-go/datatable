package datatable

import "encoding/json"

type (
	Cells []Cell
	Rows  []Cells

	Cell interface {
		Val() any
		SetVal(any)
		Rowspan() int
		SetRowspan(int)
		Colspan() int
		SetColspan(int)
	}

	CellImpl struct {
		ValValue     any `json:"val,omitempty"`
		RowspanValue int `json:"rowspan,omitempty"`
		ColspanValue int `json:"colspan,omitempty"`
	}
)

func (c *Cells) UnmarshalJSON(data []byte) (err error) {
	var items []*CellImpl
	if err = json.Unmarshal(data, &items); err != nil {
		return
	}
	cells := make(Cells, len(items))
	for i, item := range items {
		cells[i] = item
	}
	*c = cells
	return nil
}

func (c *CellImpl) Val() any {
	return c.ValValue
}

func (c *CellImpl) SetVal(a any) {
	c.ValValue = a
}

func (c *CellImpl) Rowspan() int {
	return c.RowspanValue
}

func (c *CellImpl) SetRowspan(i int) {
	c.RowspanValue = i
}

func (c *CellImpl) Colspan() int {
	return c.ColspanValue
}

func (c *CellImpl) SetColspan(i int) {
	c.ColspanValue = i
}

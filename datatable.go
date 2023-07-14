package datatable

type (
	Cells []*Cell
	Rows  []Cells

	Cell struct {
		Val     any `json:"val,omitempty"`
		Rowspan int `json:"rowspan,omitempty"`
		Colspan int `json:"colspan,omitempty"`
	}
)

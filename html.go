package datatable

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type HtmlBuilder struct {
	Begin func(w io.Writer) (err error)
	End   func(w io.Writer) (err error)
	Row   func(w io.Writer, index int) (close func() error, err error)
	Cell  func(w io.Writer, rowIndex, cellIndex int, cell Cell, attrs []string) (err error)
}

var DefaultHtmlBuilder = &HtmlBuilder{
	Begin: func(w io.Writer) (err error) {
		_, err = fmt.Fprintf(w, "<table cellspacing=\"0\" border=\"1\">\n\t<tbody>")
		return
	},
	End: func(w io.Writer) (err error) {
		_, err = fmt.Fprintf(w, "\t</tbody>\n</table>\n")
		return
	},
	Row: func(w io.Writer, index int) (close func() error, err error) {
		if _, err = w.Write([]byte("\t\t<tr>\n")); err != nil {
			return
		}
		return func() (err error) {
			_, err = w.Write([]byte("\t\t</tr>\n"))
			return
		}, nil
	},
	Cell: func(w io.Writer, rowIndex, cellIndex int, cell Cell, attrs []string) (err error) {
		if _, err = w.Write([]byte("\t\t\t<td" + strings.Join(attrs, " ") + ">")); err != nil {
			return
		}
		if val := cell.Val(); val != nil {
			if _, err = fmt.Fprint(w, val); err != nil {
				return
			}
		}
		_, err = w.Write([]byte("</td>\n"))
		return
	},
}

func (b *HtmlBuilder) Build(w io.Writer, rows Rows) (err error) {
	var (
		attrs    []string
		rowClose func() error
	)

	if err = b.Begin(w); err != nil {
		return
	}

	for rowIndex, r := range rows {
		if rowClose, err = b.Row(w, rowIndex); err != nil {
			return
		}
		for cellIndex, cell := range r {
			attrs = []string{""}
			if cell.Rowspan() > 1 {
				attrs = append(attrs, `rowspan="`+strconv.Itoa(cell.Rowspan())+`"`)
			}
			if cell.Colspan() > 1 {
				attrs = append(attrs, `colspan="`+strconv.Itoa(cell.Colspan())+`"`)
			}
			if len(attrs) == 1 {
				attrs = nil
			}
			if err = b.Cell(w, rowIndex, cellIndex, cell, attrs); err != nil {
				return
			}
		}
		if err = rowClose(); err != nil {
			return
		}
	}

	return b.End(w)
}

package helpers

import (
	"github.com/asaskevich/govalidator"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
)

func BuildTableRows(rows [][]string, tBody *gr.Element) {
	for _, row := range rows {
		var rowEls = make([]gr.Modifier, len(row))

		for ri, rowName := range row {

			if govalidator.IsURL(rowName) {
				rowEls[ri] = el.TableData(el.Anchor(attr.HRef(rowName), attr.Target("_blank"), gr.Text(rowName)))
			} else {
				rowEls[ri] = el.TableData(gr.Text(rowName))
			}
		}

		tr := el.TableRow(rowEls...)
		tr.Modify(tBody)
	}
}

func BuildTableHeader(header []string) []gr.Modifier {
	var tableHeader = make([]gr.Modifier, len(header))
	for i, head := range header {
		tableHeader[i] = el.TableHeader(gr.Text(head))
	}
	return tableHeader
}

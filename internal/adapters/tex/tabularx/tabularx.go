package tabularx

import (
	"fmt"
	"strings"
)

type Tabularx struct {
	Width fmt.Stringer
	Rows  Rows
}

type Row struct {
	Cells Cells
}

type Rows []Row

type Cell struct {
	Text fmt.Stringer
}

type Cells []Cell

func New(width fmt.Stringer) *Tabularx {
	return &Tabularx{
		Width: width,
	}
}

func (r *Tabularx) AddRow(cells ...Cell) {
	r.Rows = append(r.Rows, Row{Cells: cells})
}

func (r *Tabularx) Render() string {
	return fmt.Sprintf(`
			\begin{tabularx}{%s}{%s}
				%s
			\end{tabularx}
		`,
		r.Width,
		r.columnFormat(),
		r.rows(),
	)
}

func (r *Tabularx) columnFormat() string {
	return strings.Repeat("c", len(r.Rows[0].Cells))
}

func (r *Tabularx) rows() string {
	var rows []string

	for _, row := range r.Rows {
		var cells []string
		for _, cell := range row.Cells {
			cells = append(cells, cell.Text.String())
		}

		rows = append(rows, strings.Join(cells, " & "))
	}

	return strings.Join(rows, `\\`+"\n")
}

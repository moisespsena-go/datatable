package main

import (
	"encoding/json"
	"os"

	"github.com/moisespsena-go/datatable"
)

func main() {
	var rows datatable.Rows
	f, err := os.Open("table.json")
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(f).Decode(&rows); err != nil {
		panic(err)
	}
	f.Close()

	f, err = os.Create("table.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>DataTable Example</title>
</head>
<body>
<h1>DataTable Example</h1>
<p>Project Page: <a href="https://github.com/moisespsena-go/datatable">github.com/moisespsena-go/datatable</a></p>
<h2>Original</h2>
`); err != nil {
		panic(err)
	}
	// build html table
	if err = datatable.DefaultHtmlBuilder.Build(f, rows); err != nil {
		panic(err)
	}
	if _, err = f.WriteString("<h2>Transposed</h2>"); err != nil {
		panic(err)
	}

	// transposes
	rows.Transpose()

	// build transposed html table
	if err = datatable.DefaultHtmlBuilder.Build(f, rows); err != nil {
		panic(err)
	}
	if _, err = f.WriteString("</body>\n</html>"); err != nil {
		panic(err)
	}
}

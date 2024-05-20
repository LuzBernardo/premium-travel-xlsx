package main

import (
	"errors"

	"github.com/guark/guark/app"
)

func init() {
	app.Funcs["generate_xlsx"] = generateXLSX
	/* .Funcs["generate_xlsx"] = generateXLSX */
}

var _ app.Func = generateXLSX

func generateXLSX(ctx app.Context) (interface{}, error) {
	data := ctx.Params["data"].(string)
	if data == "" {
		return nil, errors.New("data is required")
	}

	// Create the XLSX file with the user input data
	buffer, err := Create(data)
	if err != nil {
		return nil, err
	}

	// Return the XLSX file as a byte slice
	return buffer.Bytes(), nil
}

package testhelper

import "github.com/google/go-cmp/cmp"

var ErrorTransformer = cmp.Transformer("Error", func(e error) string {
	if e == nil {
		return "<nil>"
	}

	return e.Error()
})

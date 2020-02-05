package controller

import (
	"github.com/govargo/foo-controller-operatorsdk/pkg/controller/foo"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, foo.Add)
}

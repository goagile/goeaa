package effectivity

import (
	"time"

	"github.com/khardi/gotterns/eaa/daterange"
)

//
// Person with Name
//
type Person struct {
	Name string
}

//
// Employment
//
type Employment struct {
	Company     string
	Effectivity daterange.DateRange
}

//
// IsEffective
//
func (e *Employment) IsEffective(d time.Time) bool {
	return e.Effectivity.Includes(d)
}

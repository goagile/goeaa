package effectivity

import (
	"testing"

	"github.com/khardi/gotterns/eaa/daterange"
)

//
// Person
//
func TestPersonName(t *testing.T) {
	want := "Ted"
	ted := &Person{Name: "Ted"}

	got := ted.Name

	if got != want {
		t.Fatalf("\ngot:%v\nwant:%v\n", got, want)
	}
}

//
// Employment
//
func TestEmploymentEffectivityTrue(t *testing.T) {
	want := true
	ted := &Employment{
		Company: "Google",
		Effectivity: daterange.DateRange{
			Start: daterange.April(1, 2020),
			End:   daterange.April(30, 2020),
		},
	}

	got := ted.IsEffective(daterange.April(15, 2020))

	if got != want {
		t.Fatalf("\ngot:%v\nwant:%v\n", got, want)
	}
}

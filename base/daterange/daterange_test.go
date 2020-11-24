package daterange

import (
	"testing"
)

//
// Includes
//
func TestIncludesTrue(t *testing.T) {
	want := true
	r := &DateRange{
		Start: April(15, 2020),
		End:   April(20, 2020),
	}

	got := r.Includes(April(16, 2020))

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIncludesBeforeStartFalse(t *testing.T) {
	want := false
	r := &DateRange{
		Start: April(16, 2020),
		End:   April(20, 2020),
	}

	got := r.Includes(April(14, 2020))

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIncludesAfterEndFalse(t *testing.T) {
	want := false
	r := &DateRange{
		Start: April(15, 2020),
		End:   April(20, 2020),
	}

	got := r.Includes(April(21, 2020))

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

//
// IsEmpty
//
func TestIsEmptyTrue(t *testing.T) {
	want := true
	r := &DateRange{
		Start: April(10, 2020),
		End:   April(9, 2020),
	}

	got := r.IsEmpty()

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIsEmptyFalse(t *testing.T) {
	want := false
	r := &DateRange{
		Start: April(9, 2020),
		End:   April(10, 2020),
	}

	got := r.IsEmpty()

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

//
// Equals
//
func TestEqualsTrue(t *testing.T) {
	want := true
	a := &DateRange{
		Start: April(9, 2020),
		End:   April(10, 2020),
	}
	b := &DateRange{
		Start: April(9, 2020),
		End:   April(10, 2020),
	}

	got := a.Equals(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestEqualsFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(9, 2020),
		End:   April(10, 2020),
	}
	b := &DateRange{
		Start: April(1, 2020),
		End:   April(2, 2020),
	}

	got := a.Equals(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

//
// IncludesRange
//
func TestIncludesRangeTrue(t *testing.T) {
	want := true
	a := &DateRange{
		Start: April(1, 2020),
		End:   April(10, 2020),
	}
	b := &DateRange{
		Start: April(2, 2020),
		End:   April(9, 2020),
	}

	got := a.IncludesRange(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIncludesRangeAfterFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(3, 2020),
		End:   April(5, 2020),
	}
	b := &DateRange{
		Start: April(1, 2020),
		End:   April(2, 2020),
	}

	got := a.IncludesRange(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIncludesRangeBeforeFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(3, 2020),
		End:   April(5, 2020),
	}
	b := &DateRange{
		Start: April(6, 2020),
		End:   April(7, 2020),
	}

	got := a.IncludesRange(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestIncludesRangeOverlapsFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(1, 2020),
		End:   April(9, 2020),
	}
	b := &DateRange{
		Start: April(2, 2020),
		End:   April(11, 2020),
	}

	got := a.IncludesRange(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

//
// Overlaps
//
func TestOverlapsTrue(t *testing.T) {
	want := true
	a := &DateRange{
		Start: April(1, 2020),
		End:   April(9, 2020),
	}
	b := &DateRange{
		Start: April(2, 2020),
		End:   April(11, 2020),
	}

	got := a.Overlaps(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestOverlapsBeforeFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(3, 2020),
		End:   April(4, 2020),
	}
	b := &DateRange{
		Start: April(1, 2020),
		End:   April(2, 2020),
	}

	got := a.Overlaps(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestOverlapsAfterFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(1, 2020),
		End:   April(2, 2020),
	}
	b := &DateRange{
		Start: April(3, 2020),
		End:   April(4, 2020),
	}

	got := a.Overlaps(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestOverlapsStartsBeforeEndsAfterFalse(t *testing.T) {
	want := false
	a := &DateRange{
		Start: April(2, 2020),
		End:   April(3, 2020),
	}
	b := &DateRange{
		Start: April(1, 2020),
		End:   April(4, 2020),
	}

	got := a.Overlaps(b)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

func TestOverlapsStartsBeforeEndsAfterTrue(t *testing.T) {
	want := true
	a := &DateRange{
		Start: April(2, 2020),
		End:   April(3, 2020),
	}
	b := &DateRange{
		Start: April(1, 2020),
		End:   April(4, 2020),
	}

	got := b.Overlaps(a)

	if got != want {
		t.Fatalf("got:%v, want:%v", got, want)
	}
}

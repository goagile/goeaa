package daterange

import "time"

type DateRange struct {
	Start time.Time
	End   time.Time
}

func (r *DateRange) Includes(t time.Time) bool {
	return t.After(r.Start) && t.Before(r.End)
}

func (r *DateRange) IsEmpty() bool {
	return r.Start.After(r.End)
}

func (r *DateRange) Equals(o *DateRange) bool {
	return r.Start.Equal(o.Start) && r.End.Equal(o.End)
}

func (r *DateRange) IncludesRange(o *DateRange) bool {
	return r.Includes(o.Start) && r.Includes(o.End)
}

func (r *DateRange) Overlaps(o *DateRange) bool {
	return r.Includes(o.Start) || r.Includes(o.End) || r.IncludesRange(o)
}

func April(day int, year int) time.Time {
	return time.Date(year, time.April, day, 0, 0, 0, 0, time.UTC)
}

package octav

type Date struct {
	Year  int
	Month int
	Day   int
}

type WallClock struct {
	Hour   int
	Minute int
}

// YYYY-MM-DD[HH:MM-HH:MM]
type ConferenceDate struct {
	Date  Date
	Open  WallClock
	Close WallClock
}


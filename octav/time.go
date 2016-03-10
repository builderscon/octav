package octav

import "fmt"

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (c WallClock) String() string {
	return fmt.Sprintf("%02d:%02d", c.Hour, c.Minute)
}

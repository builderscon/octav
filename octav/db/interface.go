package db

import "time"

type Venue struct {
	OID        uint64 // intenral id, used for sorting and what not
	EID        string // ID that is visible to the outside
	Name       string // Name of the venue (English)
	CreatedOn  time.Time
	ModifiedOn time.Time
}

type Room struct {
	OID        uint64 // intenral id, used for sorting and what not
	EID        string // ID that is visible to the outside
	VenueID    string // ID of the venue that this room belongs to
	Name       string // Name of the room (English)
	Capacity   uint   // How many people fit in this room? Approximation.
	CreatedOn  time.Time
	ModifiedOn time.Time
}

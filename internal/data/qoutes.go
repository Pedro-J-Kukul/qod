package data

import "time"

type Qoute struct {
	ID        int64
	Type      string
	Qoute     string
	Author    string
	CreatedAt time.Time
	Version   int32
}

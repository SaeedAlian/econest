package role_store

import "time"

type Role struct {
	Id        int
	Name      string
	CreatedAt time.Time
}

package role_db_manager

import "time"

type Role struct {
	Id        int
	Name      string
	CreatedAt time.Time
}

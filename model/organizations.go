package model

// Org model
type Org struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	RoleID string `json:"-"`
	Role   Role   `json:"-"`
	Type   string `json:"-"`
}

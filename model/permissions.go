package model

import (
	"encoding/json"
)

type Permission struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Name)
}
func (p *Permission) ConvertPermissionNew() PermissionNew {
	return PermissionNew{p.ID, p.Name}
}

type ListPermission struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (p *Permission) MapPermissionListModel() ListPermission {
	return ListPermission{p.ID, p.Name}
}

package domain

import "time"

type Folder struct {
	Id        string    `json:"ID"`
	Name      string    `json:"Name"`
	ParentID  string    `json:"ParentID"`
	OwnerID   string    `json:"ownerID"`
	CreatedAt time.Time `json:"created_at"`
}

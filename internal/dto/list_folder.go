package dto

type ListFoldersRequest struct {
	OwnerID string
}

type ListFolderResponse struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	ParentId string   `json:"parent_id,omitempty"`
	ChildIds []string `json:"child_ids,omitempty"`
}

type ListFoldersResponse struct {
	Folders []ListFolderResponse `json:"folders"`
}

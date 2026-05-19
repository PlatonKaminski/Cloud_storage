package dto

type CreateFolderRequest struct {
	Name     string
	OwnerId  string
	ParentId string
}

type CreateFolderResponse struct {
	Id string `json:"id"`
}

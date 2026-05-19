package dto

type UploadFileResponse struct {
	Id string `json:"id"`
}

type UploadFileRequest struct {
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Size      int64  `json:"size"`
	Owner_id  string `json:"owner_Id"`
	Folder_id string `json:"folder_Id"`
}

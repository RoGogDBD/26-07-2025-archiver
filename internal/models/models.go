package models

type CreateTaskRequest struct {
	Name string `json:"name"`
}

type AddURLRequest struct {
	URL string `json:"url"`
}

type Task struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	URLs    []string `json:"urls"`
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Archive string   `json:"archive_url"`
}

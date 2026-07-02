package dto

import "time"

// FileResponse mirrors a FileDto / view model returned from controllers.
// The internal storage `location` is intentionally not exposed to clients;
// file content is retrieved via the dedicated download endpoint instead.
type FileResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Extension   string     `json:"extension"`
	Size        int64      `json:"size"`
	ContentType *string    `json:"contentType"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

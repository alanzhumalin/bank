package dto

type CursorMeta struct {
	Limit      int    `json:"limit"`
	HasNext    bool   `json:"has_next"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type CursorResponse[T any] struct {
	Data []T        `json:"data"`
	Meta CursorMeta `json:"meta"`
}

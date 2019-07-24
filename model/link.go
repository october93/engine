package model

import "github.com/october93/engine/kit/globalid"

// A link is a URI saved by a user on the client for later retrieval.
type Link struct {
	ID       globalid.ID `json:"id"`
	URL      string      `json:"url"`
	Username string      `json:"username"`
}

// NewLink returns a new link instance.
func NewLink(url, username string) *Link {
	return &Link{ID: globalid.Next(), URL: url, Username: username}
}

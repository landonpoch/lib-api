package domain

import (
	"time"

	"github.com/google/uuid"
)

type Library struct {
	Books      []Book
	Page       int
	TotalBooks int
}

// Usually json tags aren't necessary on domain objects like this
// However in this case there is no difference between the raw domain
// object and the returned response for an individual book.  To avoid
// unnecessary mapping logic, I've just place the tags here for simplicity.
type Book struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	Publisher        string    `json:"publisher,omitempty"`
	CreatedDate      time.Time `json:"created_date"`
	LastModifiedDate time.Time `json:"last_modified_date"`
}

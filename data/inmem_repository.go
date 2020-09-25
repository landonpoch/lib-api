package data

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/landonpoch/lib-api/domain"
)

type InMemBookRepository struct {
	lock  sync.RWMutex
	index map[string]*int // In this case, index just helps for faster lookups
	books []domain.Book
}

var ErrDoesNotExist = errors.New("repo: DoesNotExist")

func NewInMemBookRepository() *InMemBookRepository {
	return &InMemBookRepository{
		index: make(map[string]*int),
		books: make([]domain.Book, 0),
	}
}

func (r *InMemBookRepository) GetBooks(count, page int) (domain.Library, error) {
	start := count * page
	end := start + count
	if end > len(r.books) {
		end = len(r.books)
	}
	library := domain.Library{
		Books:      r.books[start:end],
		Page:       page,
		TotalBooks: len(r.books),
	}
	return library, nil
}

func (r *InMemBookRepository) GetBook(id uuid.UUID) (*domain.Book, error) {
	_, book := r.lookup(id)
	if book == nil {
		return nil, ErrDoesNotExist
	}
	return book, nil
}

func (r *InMemBookRepository) CreateBook(book domain.Book) error {
	book.CreatedDate = time.Now().UTC()
	book.LastModifiedDate = book.CreatedDate
	index := len(r.books)
	r.lock.Lock()
	r.index[book.ID.String()] = &index
	r.books = append(r.books, book)
	r.lock.Unlock()
	return nil
}

func (r *InMemBookRepository) UpdateBook(book domain.Book) error {
	_, currentBook := r.lookup(book.ID)
	if currentBook == nil {
		return ErrDoesNotExist
	}
	currentBook.Title = book.Title
	currentBook.Author = book.Author
	currentBook.Publisher = book.Publisher
	currentBook.LastModifiedDate = time.Now().UTC()
	return nil
}

func (r *InMemBookRepository) DeleteBook(id uuid.UUID) error {
	index, _ := r.lookup(id)
	if index == nil {
		return ErrDoesNotExist
	}
	r.lock.Lock()
	// Cut out the book at the index
	copy(r.books[*index:], r.books[*index+1:])
	r.books = r.books[:len(r.books)-1] // reduce the slice by size of 1
	// Update the index
	delete(r.index, id.String())
	for _, v := range r.index {
		if *v > *index {
			*v--
		}
	}
	r.lock.Unlock()
	return nil
}

func (r *InMemBookRepository) lookup(id uuid.UUID) (*int, *domain.Book) {
	if index := r.index[id.String()]; index != nil {
		return index, &r.books[*index]
	}
	return nil, nil
}

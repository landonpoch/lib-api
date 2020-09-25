package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/landonpoch/lib-api/data"
	"github.com/landonpoch/lib-api/domain"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/mock_routes.go -package=application_mock github.com/landonpoch/lib-api/application BookRepository
type BookRepository interface {
	GetBooks(count, page int) (domain.Library, error)
	CreateBook(domain.Book) error
	GetBook(id uuid.UUID) (*domain.Book, error)
	UpdateBook(domain.Book) error
	DeleteBook(id uuid.UUID) error
}

type Routes struct {
	log    *zap.SugaredLogger
	repo   BookRepository
	Router *mux.Router
}

func NewRoutes(repository BookRepository) *Routes {
	routes := &Routes{
		log:    zap.S(),
		repo:   repository,
		Router: mux.NewRouter(),
	}
	routes.initialize()
	return routes
}

func (r *Routes) initialize() {
	// Middleware would typically be configured here
	// for logging, tracing, instrumentation and metrics,
	// cors or anything else you need.

	// Utility routes
	r.Router.HandleFunc("/", r.GetVersion).Methods("GET")
	r.Router.HandleFunc("/live", r.Liveness).Methods("GET")
	r.Router.HandleFunc("/ready", r.Readiness).Methods("GET")

	// Library routes
	// Depending on how you want to design your service, you may
	// decide to have non-idempotent creates and use a POST method
	// without providing a guaranteed unique ID from the client.
	// In this case I used PUT assuming the client would always
	// generate a unique id (UUID in this case) so that we could
	// have idempotent creates.  I used a patch for updates.  The
	// only difference here is that updates will return 404s if
	// the object doesn't already exist, whereas creates would just
	// create a new one.  Create can also act as an update in this case.
	// These of course are design decisions that would be discussed
	// with the team and the correct paradigm would be chosen for
	// the use cases at hand.
	r.Router.HandleFunc("/books", r.GetBooks).Methods("GET")
	r.Router.HandleFunc("/books/{id}", r.GetBook).Methods("GET")
	r.Router.HandleFunc("/books", r.CreateBook).Methods("PUT")
	r.Router.HandleFunc("/books", r.UpdateBook).Methods("PATCH")
	r.Router.HandleFunc("/books/{id}", r.DeleteBook).Methods("DELETE")
}

// Would usually be baked into the binary and auto incremented on new builds
var Version = "1.0.0"

func (r *Routes) GetVersion(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(Version))
}

func (r *Routes) Liveness(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func (r *Routes) Readiness(res http.ResponseWriter, req *http.Request) {
	// Readiness check should verify all necessary dependencies are available
	res.WriteHeader(http.StatusOK)
}

func (r *Routes) GetBooks(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	count, err := strconv.Atoi(params.Get("count"))
	if err != nil {
		count = 10
	}
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		page = 0
	}
	library, err := r.repo.GetBooks(count, page)
	if err != nil {
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
		return
	}
	nextHref := ""
	if ((page + 1) * count) < library.TotalBooks {
		nextHref = "http://" + req.Host + req.URL.Path +
			"?count=" + strconv.Itoa(count) +
			"&page=" + strconv.Itoa(page+1)
	}
	prevHref := ""
	if page > 0 {
		prevHref = "http://" + req.Host + req.URL.Path +
			"?count=" + strconv.Itoa(count) +
			"&page=" + strconv.Itoa(page-1)
	}
	rootURL := "http://" + req.Host + req.URL.Path + "/"
	books := make([]BookSummary, 0)
	for _, book := range library.Books {
		books = append(books, mapBook(book, rootURL))
	}
	response := QueryResponse{
		Books:      books,
		TotalBooks: library.TotalBooks,
		NextHref:   nextHref,
		PrevHref:   prevHref,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		zap.S().Errorw("An error occurred while marshaling library", "err", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(res, bytes)
}

func mapBook(book domain.Book, rootURL string) BookSummary {
	return BookSummary{
		Title: book.Title,
		Href:  rootURL + book.ID.String(),
	}
}

type QueryResponse struct {
	Books      []BookSummary `json:"books"`
	TotalBooks int           `json:"total_books"`
	NextHref   string        `json:"next_href,omitempty"`
	PrevHref   string        `json:"prev_href,omitempty"`
}

type BookSummary struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func (r *Routes) GetBook(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	parsedID, err := uuid.Parse(id)
	if err != nil {
		r.log.Errorw("An error occurred while parsing identifier", "err", err)
		writeError(http.StatusBadRequest, res, "Invalid Identifier: input format unparsable")
		return
	}
	book, err := r.repo.GetBook(parsedID)
	if err != nil {
		if err == data.ErrDoesNotExist {
			writeError(http.StatusNotFound, res, "Book not found")
			return
		}
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
		return
	}
	bytes, err := json.Marshal(book)
	if err != nil {
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
		return
	}
	writeJSON(res, bytes)
}

func (r *Routes) CreateBook(res http.ResponseWriter, req *http.Request) {
	book := domain.Book{}
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		r.log.Errorw("An error occurred while reading body", "err", err)
		writeError(http.StatusBadRequest, res, "Invalid Book: input format unparsable")
		return
	}
	if !isValid(res, book) {
		return
	}
	err = r.repo.CreateBook(book)
	if err != nil {
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (r *Routes) UpdateBook(res http.ResponseWriter, req *http.Request) {
	book := domain.Book{}
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		r.log.Errorw("An error occurred while reading body", "err", err)
		writeError(http.StatusBadRequest, res, "Invalid Book: input format unparsable")
		return
	}
	err = r.repo.UpdateBook(book)
	if err != nil {
		if err == data.ErrDoesNotExist {
			writeError(http.StatusNotFound, res, "Book not found")
			return
		}
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
	}
	res.WriteHeader(http.StatusOK)
}

func (r *Routes) DeleteBook(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	parsedID, err := uuid.Parse(id)
	if err != nil {
		r.log.Errorw("An error occurred while parsing identifier", "err", err)
		writeError(http.StatusBadRequest, res, "Invalid Identifier: input format unparsable")
		return
	}
	err = r.repo.DeleteBook(parsedID)
	if err != nil {
		if err == data.ErrDoesNotExist {
			writeError(http.StatusNotFound, res, "Book not found")
			return
		}
		writeError(http.StatusInternalServerError, res, "Unknown error occurred")
	}
}

func isValid(res http.ResponseWriter, book domain.Book) bool {
	invalidFields := make([]string, 0)
	if book.ID == uuid.Nil {
		invalidFields = append(invalidFields, "id")
	}
	if book.Title == "" {
		invalidFields = append(invalidFields, "title")
	}
	if book.Author == "" {
		invalidFields = append(invalidFields, "author")
	}
	if len(invalidFields) > 0 {
		msg := fmt.Sprintf("Invalid Book, the following fields are required: %s", strings.Join(invalidFields, ", "))
		writeError(http.StatusBadRequest, res, msg)
		return false
	}
	return true
}

func writeError(statusCode int, res http.ResponseWriter, msg string) {
	res.WriteHeader(statusCode)
	res.Write([]byte(msg))
}

func writeJSON(res http.ResponseWriter, body []byte) {
	res.Header().Add("Content-Type", "application/json")
	res.Write(body)
}

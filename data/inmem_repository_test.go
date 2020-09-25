package data_test // blackbox testing

import (
	"testing"

	"github.com/google/uuid"
	"github.com/landonpoch/lib-api/data"
	"github.com/landonpoch/lib-api/domain"
	. "github.com/smartystreets/goconvey/convey"
)

var DonQuixote = domain.Book{
	ID:     uuid.New(),
	Title:  "Don Quixote",
	Author: "Miguel de Cervantes",
}

var WarAndPeace = domain.Book{
	ID:        uuid.New(),
	Title:     "War and Peace",
	Author:    "Leo Tolstoy",
	Publisher: "The Russian Messenger",
}

var Catcher = domain.Book{
	ID:        uuid.New(),
	Title:     "The Catcher in the Rye",
	Author:    "J. D. Salinger",
	Publisher: "Little, Brown and Company",
}

var Guide = domain.Book{
	ID:        uuid.New(),
	Title:     "The Hitchhiker's Guide to the Galaxy",
	Author:    "Douglas Adams",
	Publisher: "Pan Books",
}

func TestInMemRepository(t *testing.T) {
	Convey("Can create a book", t, func() {
		repo := data.NewInMemBookRepository()
		library, err := repo.GetBooks(10, 0)
		So(err, ShouldBeNil)
		So(library, ShouldNotBeNil)
		So(library.TotalBooks, ShouldEqual, 0)
		repo.CreateBook(DonQuixote)
		book, err := repo.GetBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		So(book, ShouldNotBeNil)
		So(book.ID, ShouldEqual, DonQuixote.ID)
		So(book.Title, ShouldEqual, DonQuixote.Title)
		So(book.Author, ShouldEqual, DonQuixote.Author)
	})

	Convey("Can get books", t, func() {
		repo := data.NewInMemBookRepository()
		repo.CreateBook(DonQuixote)
		repo.CreateBook(WarAndPeace)
		repo.CreateBook(Catcher)
		library, err := repo.GetBooks(10, 0)
		So(err, ShouldBeNil)
		So(library, ShouldNotBeNil)
		So(len(library.Books), ShouldEqual, 3)
		So(library.TotalBooks, ShouldEqual, 3)
		So(library.Page, ShouldEqual, 0)
	})

	Convey("Can get range of books", t, func() {
		repo := data.NewInMemBookRepository()
		repo.CreateBook(DonQuixote)
		repo.CreateBook(WarAndPeace)
		repo.CreateBook(Catcher)
		repo.CreateBook(Guide)
		library, err := repo.GetBooks(2, 1)
		So(err, ShouldBeNil)
		So(library, ShouldNotBeNil)
		So(len(library.Books), ShouldEqual, 2)
		So(library.TotalBooks, ShouldEqual, 4)
		So(library.Page, ShouldEqual, 1)
	})

	Convey("Can get a single book", t, func() {
		repo := data.NewInMemBookRepository()
		repo.CreateBook(DonQuixote)
		book, err := repo.GetBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		So(book, ShouldNotBeNil)
		So(book.ID, ShouldEqual, DonQuixote.ID)
		So(book.Title, ShouldEqual, DonQuixote.Title)
		So(book.Author, ShouldEqual, DonQuixote.Author)
	})

	Convey("Can delete a book", t, func() {
		repo := data.NewInMemBookRepository()
		repo.CreateBook(DonQuixote)
		book, err := repo.GetBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		So(book, ShouldNotBeNil)
		So(book.ID, ShouldEqual, DonQuixote.ID)
		So(book.Title, ShouldEqual, DonQuixote.Title)
		So(book.Author, ShouldEqual, DonQuixote.Author)
		err = repo.DeleteBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		book, err = repo.GetBook(DonQuixote.ID)
		So(book, ShouldBeNil)
		So(err, ShouldEqual, data.ErrDoesNotExist)
	})

	Convey("Can update a book", t, func() {
		repo := data.NewInMemBookRepository()
		repo.CreateBook(DonQuixote)
		book, err := repo.GetBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		So(book, ShouldNotBeNil)
		So(book.ID, ShouldEqual, DonQuixote.ID)
		So(book.Title, ShouldEqual, DonQuixote.Title)
		So(book.Author, ShouldEqual, DonQuixote.Author)
		newTitle := "Don Quixote de La Mancha"
		book.Title = newTitle
		err = repo.UpdateBook(*book)
		So(err, ShouldBeNil)
		book, err = repo.GetBook(DonQuixote.ID)
		So(err, ShouldBeNil)
		So(book.Title, ShouldEqual, newTitle)
	})
}

package application_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/landonpoch/lib-api/application"
	application_mock "github.com/landonpoch/lib-api/application/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	Convey("Can create a book", t, func() {
		repo := application_mock.NewMockBookRepository(ctrl)
		repo.EXPECT().CreateBook(gomock.Any()).Return(nil)
		routes := application.NewRoutes(repo)

		reqBody := `{
			"id": "b19537af-7997-422b-a3ff-9cac51b4d59e",
			"title": "Don Quixote",
			"author": "Miguel de Cervantes"
		}`
		readerBody := strings.NewReader(reqBody)
		req := httptest.NewRequest("PUT", `/books`, readerBody).
			WithContext(ctx)
		resp := httptest.NewRecorder()

		routes.CreateBook(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
	})

	Convey("Create book returns HTTP 400 on invalid json", t, func() {
		repo := application_mock.NewMockBookRepository(ctrl)
		routes := application.NewRoutes(repo)

		reqBody := `{
			"id": "b19537af-7997-422b-a3ff-9cac51b4d59e",
			"title": "Don Quixote",
			"author": "Miguel de Cervantes",
		}`
		readerBody := strings.NewReader(reqBody)
		req := httptest.NewRequest("PUT", `/books`, readerBody).
			WithContext(ctx)
		resp := httptest.NewRecorder()

		routes.CreateBook(resp, req)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
	})
}

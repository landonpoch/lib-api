# LIB-API
A simple library REST API which allows you to perform CRUD operations on books.

## Running API
Go 1.15 is recommended.  Clone repo and run the following commands:
 - go install
 - lib-api (you may need to restart your terminal for this to be in your path)

You can import the postman collection included in the library (Library.postman_collection.json).  This was a v2.1 collection export.

## Available endpoints
There are 5 supported operations:
 - GET /books (returns the summary collection of books, supports paging and linking)
 - GET /books/{ID} (returns a single specified book)
 - PUT /books (Creates a book, is idempotent and assumes client will generate UUID)
 - PATCH /books (Updates a book, is also idempotent.  Will not update or create if ID is not found)
 - DELETE /books/{ID} (Deletes the book)

## Design
Code contains additional comments explaining rationale behind most decisions made. Also, when building production grade microservices I recommend referencing https://12factor.net/.  Some items included in this sample project such as structured logging and suggestions of a config service fall in line with 12 factor app concepts.

## Unit testing
Unit tests are included.  While coverage is not 100% the tests that are included show a good range of scenarios.

Included scenarios:
 - Mocked HTTP request / response
 - Mocked dependency
 - Happy path and non-happy path scenario in routes_test.go
 - More complete coverage in inmem_repository_test.go

This means that enough scenarios were provided to illustrate how unit testing would work if complete coverage of routes.go and inmem_repository.go were necessary.  Typically I do not test main.go and bootstrap logic as it is simply wiring up dependencies.  Dependencies are written to an interface so that mocks can be generated for easy testing.

## Integration testing
Typically end to end testing would be written as a separate process, maybe even in a separate scripting language for quick turnaround.  This would usually hook into a CI/CD pipeline and merges so that all code can be thouroughly tested prior to merging to master and going through the pipeline to production.
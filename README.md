## LIBRERIA

### Golang/REST Code Challenge

Using Go as your language, create a CRUD API to manage a list of Books, fulfilling the following requirements:

1. Books should have the following Attributes:
- Title
- Author
- Publisher
- Publish Date
- Rating (1-3)
- Status (CheckedIn, CheckedOut)
2. Each endpoint should have test coverage of both successful and failed (due to user error) requests.
3. Use a data store of your choice.
4. Include unit tests that exercise the endpoints such that both 200-level and 400-level responses are induced.
5. The app should be stood up in Docker, and client code (such as cURL requests, and your unit tests) 
   should be executed on the host machine, against the containerized app.
6. Send the project along as a git repository.
7. Please do not use go-swagger to generate the server side code.

### Run

To run in Docker:

```
make run
```

To run as a binary:
```
docker-compose up -d postgres 
make build
./bin/libreria
```

### Tests

To run unit tests:
```
make test-unit
```

To run integration tests:
```
make test-integration
```


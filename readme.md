### Learn Golang using "Building Modern Web Applications with Go (Golang) Trevor Sawler, Ph.D." on Udemy

#### To run project
> ./run.sh

website available on http://localhost:8080

#### To run tests
> go test ./...


#### Check tests coverage
From folder /cmd/web/
> go test -coverprofile=coverage.out && go tool cover -html=coverage.out


#### Run migrations
> soda migrate



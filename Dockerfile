FROM golang:1.15.5 as base

WORKDIR app
COPY . .

# Run the build
FROM base AS build
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -o /bin/libreria .
RUN ls -l

# Build the target runtime layer
FROM alpine:3.12.0 as runtime
RUN ls -l
COPY --from=build /bin/libreria /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/libreria"]

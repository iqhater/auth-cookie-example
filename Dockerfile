FROM golang:alpine3.12 as builder

# enable go modules
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

# run only if go.mod or go.sum will be changed (cache)
RUN go mod download

COPY . .

#run test with coverage and goes to test_data folder. Must be empty 
CMD go test -race -v -cover ./...

# build binary without debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w'

# generate clean, final image
FROM scratch

# copy golang binary into container
COPY --from=builder /app/auth_sign_in /app/

# Specify the container's entrypoint as the action
ENTRYPOINT ["/app/auth_sign_in"]
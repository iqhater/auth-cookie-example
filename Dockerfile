# docker build --rm -t iqhater/auth_sign_in_local:latest .
# docker run --env-file .env -p 8080:8080 --rm -it iqhater/auth_sign_in_local:latest

FROM golang:alpine3.12 as builder

# enable go modules
ENV GO111MODULE=on
ENV PORT=5050

WORKDIR /app
# RUN apk update && apk upgrade && apk add --no-cache ca-certificates
# RUN update-ca-certificates

COPY go.mod .
COPY go.sum .

# run only if go.mod or go.sum will be changed (cache)
RUN go mod download

COPY . .

#run test with coverage and goes to test_data folder. Must be empty 
CMD go test -race -v -coverprofile=coverage.txt -covermode=atomic ./...

# build binary without debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w'

# generate clean, final image
FROM scratch

WORKDIR /app

# copy golang binary into container
COPY --from=builder /app/auth_sign_in /app/
COPY --from=builder /app/public /app/public

# Specify the container's entrypoint as the action
ENTRYPOINT ["/app/auth_sign_in"]
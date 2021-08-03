FROM golang:1.16.6-alpine3.14 as build

WORKDIR /app

# Fetch dependencies to cache them for subsequent builds
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

# Copy only the build output into the final image
# so we don't pull all the build tooling in
FROM alpine:3.14

WORKDIR /app
COPY --from=build /app/main ./

EXPOSE 4000

CMD ["./main"]
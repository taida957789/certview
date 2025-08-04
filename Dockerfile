FROM golang:1.24.3-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM scratch
COPY --from=build /app/app /app/app
EXPOSE 8080
CMD ["/app/app"]

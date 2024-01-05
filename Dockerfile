FROM golang:1.20-bullseye

COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o pos main.go
# expose port 8080
EXPOSE 8080
CMD ["./pos"]
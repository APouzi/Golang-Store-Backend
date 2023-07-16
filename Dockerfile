
FROM golang:1.20

# RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY /Golang-Store-Backend/go.mod /Golang-Store-Backend/go.sum ./
RUN go mod download && go mod verify

# COPY . .
COPY /Golang-Store-Backend .

RUN go build -o golang-shop .
# RUN chmod +x app

EXPOSE 8000
CMD ["./golang-shop", "-initdb", "t", "-initView", "t"]
# CMD ["./golang-shop", "-initdb", "t", "-initView", "f"]
# CMD ["./golang-shop"]

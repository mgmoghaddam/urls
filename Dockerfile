FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8090
CMD ["/app/main"]
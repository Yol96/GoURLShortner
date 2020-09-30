FROM golang:latest

FROM golang:latest 
RUN mkdir /app 
ADD . /app/
WORKDIR /app/cmd/apiserver
RUN go build -o ./cmd/apiserver
CMD ["/app/apiserver"]

EXPOSE 8080
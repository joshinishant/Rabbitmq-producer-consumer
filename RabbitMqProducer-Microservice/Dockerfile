FROM golang:1.17.5-alpine3.15
WORKDIR /app
COPY . .
RUN go get ./...
CMD [ "go" ,"run","Restendpoints.go" ]
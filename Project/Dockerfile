FROM golang
COPY . /app
WORKDIR /app
RUN go mod download
ENTRYPOINT [ "go", "run", "./cmd/api" ]
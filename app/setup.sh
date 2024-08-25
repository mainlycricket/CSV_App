go mod init app.com/app
go mod tidy
go fmt
goimports -w *.go
go build .
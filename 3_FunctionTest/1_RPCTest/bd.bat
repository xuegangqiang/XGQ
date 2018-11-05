cls
@go install .\Model
@go build -ldflags="-s -w" .\Svr
@go build -ldflags="-s -w" .\Clt

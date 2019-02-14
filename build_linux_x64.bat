set GOARCH=amd64
set GOOS=linux
go build -o xiniu_api_64 main.go ilog.go

echo "编译完成，任意键退出"

pause
#/bin/sh

echo build命令 yop + cmd/yomail/main.go
go build -o yop cmd/yomail/main.go

echo build命令 yop + cmd/server/main.go
go build -o yop cmd/server/main.go

echo 执行运行
./yop start


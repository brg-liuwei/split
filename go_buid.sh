cur_path=`pwd`
export GOPATH=$cur_path:$GOPATH
cd src
go test config
go test util
cd ..
go build -o bin/splitService src/main.go

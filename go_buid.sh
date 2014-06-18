cur_path=`pwd`
export GOPATH=$cur_path:$GOPATH
export CGO_CFLAGS="-I/Users/wliu/Code/split/lib/leveldb/include"
export CGO_LDFLAGS="-L/Users/wliu/Code/split/lib/leveldb"

cd lib/leveldb/
make

cd ../../
cd src/
go test config
go test util
go test kv
cd ..
go build -o bin/splitService src/main.go

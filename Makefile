export CGO_CFLAGS=-I/home/vagrant/rocksdb/include
export CGO_LDFLAGS=-L/home/vagrant/rocksdb -lrocksdb -lstdc++ -lm -lsnappy

.PHONY: server

server:
	go build -o bin/rkv-server cmd/rkv-server/*

test:
	go test -race -v server/*

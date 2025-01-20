package main

import "github.com/rocksun/hellogo/zinxx/znet"

func main() {
	s := znet.NewServer("[zinxV0.1]")
	s.Serve()
}

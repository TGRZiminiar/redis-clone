package main

import (
	"tgrziminiar/redisclone/cache"
	"tgrziminiar/redisclone/server"
)

func main() {
	cfg := server.ServerConfigs{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	server := server.NewServer(&cfg, cache.NewCache())
	server.Start()
}

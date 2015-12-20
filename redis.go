package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

type redisInfos struct {
	path string
	host string
	port uint
	auth string
}

func (infos *redisInfos) parse() {
	redisURL, err := url.Parse(infos.path)
	if err != nil {
		log.Println(err)
		return
	}

	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			infos.auth = password
		}
	}

	infos.host = redisURL.Host
	infos.host = strings.Split(redisURL.Host, ":")[0]
	port, err := strconv.Atoi(strings.Split(redisURL.Host, ":")[1])
	if err != nil {
		log.Println(err)
		return
	}
	infos.port = uint(port)
}

package main

import (
	"log"
	"time"
)

func startWorkers() {
	ticker := time.NewTicker(time.Second * 5)
	for _ = range ticker.C {
		go followUsers()
	}
}

func followUsers() {
	ids, err := redis.ZRangeByScore("queue:users", 0, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return
	}

	for _, id := range ids {
		user, err := intercomAPI.Users.FindByID(id)
		if err != nil {
			log.Println(err)

			redis.ZRem("queue:users", id)
			return
		}

		username := ""
		for _, profile := range user.SocialProfiles.SocialProfiles {
			if profile.Name == "Twitter" {
				username = profile.Username
			}
		}

		if username != "" {
			_, err = twitterAPI.FollowUser(username)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Followed user " + username)
		}

		redis.ZRem("queue:users", id)
	}
}

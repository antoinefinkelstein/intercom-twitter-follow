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
	ids, err := redis.ZRange("queue:users", 0, -1)
	if err != nil {
		log.Println(err)
		return
	}

	for _, id := range ids {
		user, err := intercomAPI.Users.FindByID(id)
		if err != nil {
			log.Println(err)
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

package redis_helper

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
)

type User struct {
	Name  string
	Age   int
	Score int
}

func TestRC_SetIfFailRetry(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rc := NewRC(client)

	ctx := context.Background()

	u := User{
		Name:  "mtg",
		Age:   18,
		Score: 100,
	}
	userJson, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	client.Set(ctx, "test", userJson, 0)
	err = rc.SetIfFailRetry(ctx, "test", func(currentValue interface{}) (interface{}, error) {
		u := &User{}
		err = json.Unmarshal(currentValue.([]byte), u)
		if err != nil {
			log.Fatal(err)
		}
		u.Score = u.Score + 1
		userJson, err := json.Marshal(u)
		if err != nil {
			log.Fatal(err)
		}
		return userJson, nil
	}, 0)
	if err != nil {
		log.Fatal(err)
	}
	newJson := client.Get(ctx, "test").Val()
	newUser := &User{}
	err = json.Unmarshal([]byte(newJson), newUser)
	if err != nil {
		log.Fatal(err)
	}
	if newUser.Score != 101 {
		log.Fatal("fail")
	}
}

package cache

import (
	"github.com/go-redis/redis"
	"lambda/crawlers"
	"lambda/models"
	"strconv"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

const GameViewKey = "GameView"

func GameView(id int64) {
	gameId := strconv.FormatInt(id, 10)

	var keys []string
	keys, _, err := RedisClient.ZScan(GameViewKey, 0, gameId, 1).Result()
	if len(keys) != 0 {
		RedisClient.ZIncrBy(GameViewKey, float64(1), gameId)
	} else {
		RedisClient.ZAdd(GameViewKey, redis.Z{
			Score:  float64(1),
			Member: gameId,
		})
	}
	if err != nil {
		panic(err)
	}
}

func GameViewCollect() {
	keys, _, err := RedisClient.ZScan(GameViewKey, 0, "", 10000).Result()

	if err != nil {
		panic(err)
	}
	if len(keys) != 0 {
		for i := 0; i < len(keys)/2; i++ {
			key, _ := strconv.ParseInt(keys[i*2], 10, 64)
			value, _ := strconv.Atoi(keys[(i*2)+1])
			var total int
			v, err := models.GetGameById(key)
			if err != nil {
				panic(err)
			} else {
				total = v.View + value
				crawlers.UpdateSteamGame(v, total)
			}
		}
		RedisClient.Del(GameViewKey)
	}
}

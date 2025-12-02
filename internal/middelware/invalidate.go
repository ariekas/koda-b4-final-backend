package middelware

import (
	"context"
	"fmt"
	"shortlink/internal/config"
)

func InvalidateRedis(pattern string) {
	rdb := config.Redis()

	iter := rdb.Scan(context.Background(), 0, pattern, 0).Iterator()
	for iter.Next(context.Background()) {
		rdb.Del(context.Background(), iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Printf("error: failed to delete cache: %s\n", err)
	}
}
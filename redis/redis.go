package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type RC struct {
	maxRetry int
	client   *redis.Client
}

func NewRC(client *redis.Client, maxRetry ...int) *RC {
	m := 20
	if len(maxRetry) == 0 {
		m = maxRetry[0]
	}
	return &RC{
		maxRetry: m,
		client:   client,
	}
}

type FuncValue func(currentValue interface{}) (newValue interface{}, err error)

func (rc *RC) SetIfFailRetry(ctx context.Context, key string, funcValue FuncValue, expiration time.Duration) error {

	txf := func(tx *redis.Tx) error {

		currentValue, err := tx.Get(ctx, key).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			newValue, err := funcValue(currentValue)
			if err != nil {
				return err
			}
			pipe.Set(ctx, key, newValue, expiration)
			return nil
		})
		return err
	}
	for i := 0; i < rc.maxRetry; i++ {
		err := rc.client.Watch(ctx, txf, key)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}
	return errors.New("increment reached maximum number of retries")
}

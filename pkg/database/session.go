package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

type Session interface {
	Set(ctx context.Context, data SessionData) (string, error)
	GetAndUpdate(ctx context.Context, refreshToken string) (string, SessionData, error)
}

type SessionData struct {
	UserID   int
	UserType string
}

type RedisSession struct {
	rdb               *redis.Client
	sessionExpiration time.Duration
	dbResponseTime    time.Duration
}

func InitRedisSession() Session {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString(config.SessionHost),
			viper.GetInt(config.SessionPort),
		),
		Password: viper.GetString(config.SessionPassword),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %s", err.Error()))
	}

	return RedisSession{
		rdb:               rdb,
		sessionExpiration: time.Duration(viper.GetInt(config.SessionSaveTime)) * time.Hour * 24,
		dbResponseTime:    time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
	}
}

func (r RedisSession) Set(ctx context.Context, data SessionData) (string, error) {
	ctx, cansel := context.WithTimeout(ctx, r.dbResponseTime)
	defer cansel()

	sessionDataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	uuidBytes, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	key := uuidBytes.String()

	err = r.rdb.Set(ctx, key, sessionDataJSON, r.sessionExpiration).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

func (r RedisSession) GetAndUpdate(ctx context.Context, refreshToken string) (string, SessionData, error) {
	var userData SessionData

	ctx, cansel := context.WithTimeout(ctx, r.dbResponseTime)
	defer cansel()

	data, err := r.rdb.Get(ctx, refreshToken).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			return "", SessionData{}, utils.NeedToAuthorize
		default:
			return "", SessionData{}, err
		}
	}

	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return "", SessionData{}, err
	}

	key, err := r.Set(ctx, userData)
	if err != nil {
		return "", SessionData{}, err
	}

	err = r.rdb.Del(ctx, refreshToken).Err()
	if err != nil {
		return "", SessionData{}, err
	}

	return key, userData, nil
}

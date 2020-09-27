package store

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Yol96/GoURLShortner/internal/app"
	"github.com/catinello/base62"
	"github.com/go-redis/redis"
)

// UserRepository struct contains pointer to configured store
type UserRepository struct {
	store *Store
}

// LinkInfo contains response struct fields
type LinkInfo struct {
	Address        string `json:"address"`
	CreatedAt      string `json:"created_at"`
	ExpirationTime int64  `json:"expiration_time"`
}

const (

	// RedisURLKey the redis auto increment key
	RedisURLKey = "next.url.id"

	// ShortLinkKey the key for shortlink to url
	ShortLinkKey = "urlshortlink:%s:url"

	// URLHashKey the key for urlhash to url
	URLHashKey = "urlhash:%s:url"

	// ShortLinkDetailKey the key for urlshortlink to urlshortlink detail
	ShortLinkDetailKey = "urlshortlink:%s:detail"
)

// Create creates a new key-pair values in redis db
func (r *UserRepository) Create(url string, expirationTime int64) (string, error) {
	// Hashing the string
	sha := sha1.New()
	sha.Write([]byte(url))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Getting a value by key (returns error when key does not exist)
	res, err := r.store.Cli.Get(fmt.Sprintf(URLHashKey, hash)).Result()

	// Ignore, if doesn't exist or url expired
	if err == redis.Nil {

	} else if err != nil {
		return "", err
	} else {
		if res == "{}" {

		} else {
			return res, nil
		}
	}

	if err := r.store.Cli.Incr(RedisURLKey).Err(); err != nil {
		return "", err
	}

	urlID, err := r.store.Cli.Get(RedisURLKey).Int()
	if err != nil {
		return "", err
	}

	encryptedID := base62.Encode(urlID)

	err = r.store.Cli.Set(
		fmt.Sprintf(ShortLinkKey, encryptedID),
		url,
		time.Minute*time.Duration(expirationTime)).Err()

	if err != nil {
		return "", err
	}

	err = r.store.Cli.Set(
		fmt.Sprintf(URLHashKey, hash),
		encryptedID,
		time.Minute*time.Duration(expirationTime)).Err()

	if err != nil {
		return "", err
	}

	shortLinkInfo := &LinkInfo{
		Address:        url,
		CreatedAt:      time.Now().String(),
		ExpirationTime: expirationTime,
	}

	json, err := json.Marshal(shortLinkInfo)
	if err != nil {
		return "", err
	}

	err = r.store.Cli.Set(
		fmt.Sprintf(ShortLinkDetailKey, encryptedID),
		json,
		time.Minute*time.Duration(expirationTime)).Err()

	if err != nil {
		return "", nil
	}

	return encryptedID, err
}

// Info returns info originial url by shprten url (if exist)
func (r *UserRepository) Info(encryptedID string) (interface{}, error) {
	json, err := r.store.Cli.Get(fmt.Sprintf(ShortLinkDetailKey, encryptedID)).Result()

	// Return 404, if doesn't exist. Error, if smth went wrong. Json if exist.
	if err == redis.Nil {
		return nil, app.StatusError{
			Code: 404,
			Err:  errors.New("Unknown url"),
		}
	} else if err != nil {
		return nil, err
	} else {
		return json, nil
	}
}

// Get returns original url, if exist
func (r *UserRepository) Get(encryptedID string) (string, error) {
	json, err := r.store.Cli.Get(fmt.Sprintf(ShortLinkKey, encryptedID)).Result()
	if err == redis.Nil {
		return "", app.StatusError{
			Code: 404,
			Err:  errors.New("Unknown url"),
		}
	} else if err != nil {
		return "", err
	} else {
		return json, nil
	}
}

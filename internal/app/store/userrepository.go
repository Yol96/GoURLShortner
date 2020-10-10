package store

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Yol96/GoURLShortner/internal/app"

	"github.com/Yol96/GoURLShortner/internal/app/model"
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
	ExpirationTime int    `json:"expiration_time"`
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
func (r *UserRepository) Create(sl *model.Link) error {
	// Hashing the string
	sha := sha1.New()
	sha.Write([]byte(sl.OriginalAddress))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Getting a value by key (returns error when key does not exist)
	res, err := r.store.Cli.Get(fmt.Sprintf(URLHashKey, hash)).Result()

	// Ignore, if doesn't exist or url expired
	if err == redis.Nil {

	} else if err != nil {
		return err
	} else {
		if res == "{}" {

		} else {
			sl.ShortLink = res
			return nil
		}
	}

	if err := r.store.Cli.Incr(RedisURLKey).Err(); err != nil {
		return err
	}

	urlID, err := r.store.Cli.Get(RedisURLKey).Int()
	if err != nil {
		return err
	}

	sl.ShortLink = base62.Encode(urlID)

	if err = r.store.Cli.Set(
		fmt.Sprintf(ShortLinkKey, sl.ShortLink),
		sl.OriginalAddress,
		time.Minute*time.Duration(sl.ExpirationTime)).Err(); err != nil {
		return err
	}

	if err = r.store.Cli.Set(
		fmt.Sprintf(URLHashKey, hash),
		sl.ShortLink,
		time.Minute*time.Duration(sl.ExpirationTime)).Err(); err != nil {
		return err
	}

	json, err := json.Marshal(sl)
	if err != nil {
		return err
	}

	if err = r.store.Cli.Set(
		fmt.Sprintf(ShortLinkDetailKey, sl.ShortLink),
		json,
		time.Minute*time.Duration(sl.ExpirationTime)).Err(); err != nil {
		return err
	}

	return nil
}

// Info returns info originial url by shprten url (if exist)
func (r *UserRepository) Info(sl *model.Link) error {
	json, err := r.store.Cli.Get(fmt.Sprintf(ShortLinkDetailKey, sl.ShortLink)).Result()
	if err != nil {
		return app.ErrRecordNotFound
	}

	if err := sl.ParseStringIntoStruct(json); err != nil {
		return err
	}

	// Return 404, if doesn't exist. Error, if smth went wrong. Json if exist.
	if err == redis.Nil {
		return app.ErrRecordNotFound
	} else if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

// Get returns original url, if exist
func (r *UserRepository) Get(encryptedID string) (string, error) {
	url, err := r.store.Cli.Get(fmt.Sprintf(ShortLinkKey, encryptedID)).Result()
	if err == redis.Nil {
		return "", app.ErrRecordNotFound
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

package service

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"
	"urls/configs"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// URL is the struct for storing URLs in the MySQL database
type URL struct {
	ID        int       `gorm:"primary_key"`
	Original  string    `gorm:"not null"`
	Short     string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"not null"`
	Hit       int       `gorm:"not null"`
}

var db *gorm.DB
var cache *redis.Client

var mu sync.Mutex

func init() {
	configs.InitViper()
	db = configs.InitDB()
	db.AutoMigrate(&URL{})
	cache = configs.InitRedis()
	fmt.Println(cache)
}

// Shorten generates a short URL and stores the original and short URL in the MySQL database and Redis cache.
func Shorten(original string) (string, error) {
	shortURL, err := cache.Get(original).Result()
	if err == nil {
		// Return the short URL if it exists in the cache
		return shortURL, nil
	}
	var url URL
	if err := db.Model(&URL{}).Where("original = ?", original).First(&url).Error; err == nil {
		return url.Short, nil
	}
	shortURL = generateShortURL(original)
	err = cache.Set(original, shortURL, 0).Err()
	if err != nil {
		println(err)
		return "", err
	}
	err = db.Create(&URL{Original: original, Short: shortURL, CreatedAt: time.Now(), Hit: 0}).Error
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// Expand retrieves the original URL from the MySQL database and Redis cache, and increases the click ratio (hit rate)
func Expand(short string) (string, error) {
	original, err := cache.Get(short).Result()
	if err == redis.Nil {
		var url URL
		if err := db.Model(&URL{}).Where("short = ?", short).First(&url).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return "", errors.New("Short URL not found")
			}
			return "", err
		}
		db.Model(&url).Update("hit", gorm.Expr("hit + ?", 1))
		cache.Set(short, url.Original, 0)
		return url.Original, nil
	} else if err != nil {
		return "", err
	}
	mu.Lock()
	db.Model(&URL{}).Where("short = ?", short).Update("hit", gorm.Expr("hit + ?", 1))
	mu.Unlock()
	return original, nil
}

func GetURLHits(shortURL string) (int, error) {
	var hits int
	err := db.Model(&URL{}).Where("short = ?", shortURL).Select("hit").Row().Scan(&hits)
	if err != nil {
		return 0, err
	}
	return hits, nil
}

func generateShortURL(originalURL string) string {
	id := bijectiveFunction(originalURL)
	u, _ := url.Parse(originalURL)
	var mainDomain string
	if u.Host != "" {
		mainDomain = u.Host
	} else {
		mainDomain = u.Path
	}
	if strings.HasPrefix(mainDomain, "www.") {
		mainDomain = strings.ReplaceAll(mainDomain, "www.", "")
	}

	//removing vowels I believe can help it to be more human-friendly
	uHost := removeVowels(mainDomain)

	// For human-readable function add first 3 characters of originalURL
	var shortUrl string
	if len(uHost) > 2 {
		shortUrl = uHost[:3] + id[:3]
	} else {
		shortUrl = mainDomain[:3] + id[:3]
	}

	return shortUrl
}

func removeVowels(u string) string {
	for _, c := range []string{"a", "e", "i", "o", "u"} {

		u = strings.ReplaceAll(u, c, "")
	}
	return u
}

func bijectiveFunction(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	num := int(h.Sum32())
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base := 62
	if num == 0 {
		return string(alphabet[0])
	}

	chars := make([]byte, 0)
	length := int(math.Floor(math.Log(float64(num))/math.Log(float64(base))) + 1)

	for i := length - 1; i >= 0; i-- {
		pow := int(math.Pow(float64(base), float64(i)))
		c := num / pow
		chars = append(chars, alphabet[c])
		num -= c * pow
	}

	return string(chars)
}

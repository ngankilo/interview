package request_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-streaming/model"
	"go-streaming/utils"
	"io"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

var rdb *redis.Client

func InitRedisCache(config *model.ConfigRedisCache) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Redis cache")
		return err
	}
	return nil
}

func RequestApi(path, method string, data interface{}, ttl int) string {
	prvKey, err := utils.LoadFile("config-jwtrs256-key")
	if err != nil {
		log.Error().Err(err).Msg("Failed to load private key")
		return ""
	}
	apiRequest, err := utils.LoadConfigRequestApi("config-request-api")
	if err != nil {
		log.Error().Err(err).Msg("Failed to load API config")
		return ""
	}

	loginId := apiRequest.Login.Username
	pwd := apiRequest.Login.Password
	domain := apiRequest.Domain
	slpitPath := strings.Split(path, "/")
	service := slpitPath[0]
	apiVersion := apiRequest.Version[service]
	timeout := time.Now().Unix() + 30

	tokenInternal, err := model.GenToken(loginId, pwd, prvKey, timeout)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return ""
	}
	tokenSent := fmt.Sprintf("Bearer %s", tokenInternal)

	var sendData *bytes.Buffer
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal request data")
			return ""
		}
		sendData = bytes.NewBuffer(jsonData)
	} else {
		sendData = bytes.NewBuffer([]byte("{}"))
	}

	uri := fmt.Sprintf("%s/%s/%s", domain, apiVersion, path)
	log.Info().Str("uri", uri).Msg("Sending request to internal API")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = time.Second
	client.RetryWaitMax = time.Second * 5

	req, err := retryablehttp.NewRequest(method, uri, sendData)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Failed to create request")
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", tokenSent)

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Request failed")
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Failed to read response body")
		return ""
	}

	if res.StatusCode != 200 {
		log.Warn().Int("status", res.StatusCode).Str("uri", uri).Msg("Non-200 response")
		return ""
	}

	log.Info().Str("uri", uri).Str("response", string(body)).Msg("Received response")
	return string(body)
}

func GetDataRedis(key string) (string, error) {
	if rdb == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to get Redis value")
		return "", err
	}
	return val, nil
}

func SetDataRedis(key, body string, ttl int) error {
	if rdb == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	duration := time.Duration(ttl) * time.Second
	if err := rdb.Set(context.Background(), key, body, duration).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to set Redis value")
		return err
	}
	return nil
}

func CheckRealtimeMarketData(userId string, ttl int) bool {
	keyRedis := fmt.Sprintf("market-data:%s", userId)
	val, err := GetDataRedis(keyRedis)
	if err == nil && val != "" {
		return parseRealtimeMarketData(val)
	}

	path := fmt.Sprintf("user/market-data?user_id=%s", userId)
	log.Info().Str("user_id", userId).Str("path", path).Msg("Checking realtime market data")

	result := RequestApi(path, "GET", nil, ttl)
	if result == "" {
		return false
	}

 _

System: **Refactored Code Continuation**

Below, I continue the refactoring of the remaining files (`request_api.go`, `sse.go`, `option.go`, and `main.go`) to complete the full refactor of the project. Each file incorporates the optimizations outlined in the analysis, including concurrency safety, graceful shutdown, configurable settings, enhanced error handling, and structured logging with `zerolog`. All artifacts are wrapped in `<xaiArtifact>` tags with unique UUIDs, proper titles, and content types.

---

#### **7. request_api.go (Continued)**
**Changes** (continued from previous):
- Completed the `CheckRealtimeMarketData` and `CheckLatestDeviceId` functions with proper error handling.
- Added JSON response validation.
- Used configurable TTLs from the `Config` struct.
- Ensured Redis client reuse with initialization in `main.go`.

<xaiArtifact artifact_id="69813342-4b7f-4f5c-9092-d270892a32f3" artifact_version_id="97f39fa1-effd-4217-8c11-a370aa8fca35" title="request_api.go" contentType="text/go">
package request_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-streaming/model"
	"go-streaming/utils"
	"io"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

var rdb *redis.Client

func InitRedisCache(config *model.ConfigRedisCache) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Redis cache")
		return err
	}
	log.Info().Msg("Redis cache connected")
	return nil
}

func RequestApi(path, method string, data interface{}, ttl int) string {
	prvKey, err := utils.LoadFile("config-jwtrs256-key")
	if err != nil {
		log.Error().Err(err).Msg("Failed to load private key")
		return ""
	}
	apiRequest, err := utils.LoadConfigRequestApi("config-request-api")
	if err != nil {
		log.Error().Err(err).Msg("Failed to load API config")
		return ""
	}

	loginId := apiRequest.Login.Username
	pwd := apiRequest.Login.Password
	domain := apiRequest.Domain
	splitPath := strings.Split(path, "/")
	service := splitPath[0]
	apiVersion := apiRequest.Version[service]
	timeout := time.Now().Unix() + 30

	tokenInternal, err := model.GenToken(loginId, pwd, prvKey, timeout)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return ""
	}
	tokenSent := fmt.Sprintf("Bearer %s", tokenInternal)

	var sendData *bytes.Buffer
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal request data")
			return ""
		}
		sendData = bytes.NewBuffer(jsonData)
	} else {
		sendData = bytes.NewBuffer([]byte("{}"))
	}

	uri := fmt.Sprintf("%s/%s/%s", domain, apiVersion, path)
	log.Info().Str("uri", uri).Msg("Sending request to internal API")

	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = time.Second
	client.RetryWaitMax = time.Second * 5

	req, err := retryablehttp.NewRequest(method, uri, sendData)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Failed to create request")
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", tokenSent)

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Request failed")
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("Failed to read response body")
		return ""
	}

	if res.StatusCode != http.StatusOK {
		log.Warn().Int("status", res.StatusCode).Str("uri", uri).Msg("Non-200 response")
		return ""
	}

	log.Info().Str("uri", uri).Str("response", string(body)).Msg("Received response")
	return string(body)
}

func GetDataRedis(key string) (string, error) {
	if rdb == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	start := time.Now()
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to get Redis value")
		// metrics.Observe("redis.get_latency", time.Since(start).Seconds())
		return "", err
	}
	// metrics.Observe("redis.get_latency", time.Since(start).Seconds())
	return val, nil
}

func SetDataRedis(key, body string, ttl int) error {
	if rdb == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	start := time.Now()
	duration := time.Duration(ttl) * time.Second
	if err := rdb.Set(context.Background(), key, body, duration).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to set Redis value")
		// metrics.Observe("redis.set_latency", time.Since(start).Seconds())
		return err
	}
	// metrics.Observe("redis.set_latency", time.Since(start).Seconds())
	return nil
}

func parseRealtimeMarketData(result string) bool {
	var jsonResponse []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResponse); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal market data JSON")
		return false
	}

	for _, item := range jsonResponse {
		data, ok := item["data"].([]interface{})
		if !ok {
			log.Warn().Msg("Invalid market data format")
			continue
		}
		for _, element := range data {
			itemMap, ok := element.(map[string]interface{})
			if !ok {
				log.Warn().Msg("Invalid market data element")
				continue
			}
			marketDataType, ok := itemMap["market_data_type"].(float64)
			if !ok {
				log.Warn().Msg("Invalid market data type")
				continue
			}
			if int(marketDataType) == 3 { // userTypeRealtime
				return true
			}
		}
	}
	return false
}

func CheckRealtimeMarketData(userId string, ttl int) bool {
	keyRedis := fmt.Sprintf("market-data:%s", userId)
	val, err := GetDataRedis(keyRedis)
	if err == nil && val != "" {
		return parseRealtimeMarketData(val)
	}

	path := fmt.Sprintf("user/market-data?user_id=%s", userId)
	log.Info().Str("user_id", userId).Str("path", path).Msg("Checking realtime market data")

	result := RequestApi(path, "GET", nil, ttl)
	if result == "" {
		return false
	}

	if err := SetDataRedis(keyRedis, result, ttl); err != nil {
		log.Error().Err(err).Str("key", keyRedis).Msg("Failed to cache market data")
	}

	return parseRealtimeMarketData(result)
}

func CheckLatestDeviceId(userId, deviceId string, ttl int) bool {
	keyRedis := fmt.Sprintf("latest-device-id:%s:%s", userId, deviceId)
	val, err := GetDataRedis(keyRedis)
	if err == nil && val != "" {
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val), &jsonResponse); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal device ID JSON")
			return false
		}
		isLatest, ok := jsonResponse["isLatest"].(bool)
		return ok && isLatest
	}

	path := "auth/check-latest-device-id"
	log.Info().Str("user_id", userId).Str("device_id", deviceId).Str("path", path).Msg("Checking latest device ID")

	requestBody := map[string]interface{}{
		"data": map[string]string{
			"user_id":   userId,
			"device_id": deviceId,
		},
	}

	result := RequestApi(path, "POST", requestBody, ttl)
	if result == "" {
		return false
	}

	if err := SetDataRedis(keyRedis, result, ttl); err != nil {
		log.Error().Err(err).Str("key", keyRedis).Msg("Failed to cache device ID")
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResponse); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal device ID JSON")
		return false
	}

	isLatest, ok := jsonResponse["isLatest"].(bool)
	return ok && isLatest
}
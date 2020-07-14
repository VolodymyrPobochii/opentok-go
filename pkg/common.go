package pkg

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"
)

func Defaults(dst, src map[string]interface{}) map[string]interface{} {
	for key, value := range src {
		if dst[key] == nil {
			dst[key] = value
		}
	}
	return dst
}

func Pick(src map[string]interface{}, keys []string) map[string]interface{} {
	dst := make(map[string]interface{})
	for _, key := range keys {
		dst[key] = src[key]
	}
	return dst
}

func Includes(src []interface{}, target interface{}) bool {
	for _, value := range src {
		if value == target {
			return true
		}
	}
	return false
}

func Clone(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

var (
	last   int64
	repeat int64
)

func Nonce(length int) func() (int64, error) {

	if length == 0 {
		length = 15
	}
	return func() (int64, error) {
		t := time.Now()
		millis := t.UnixNano() / int64(time.Millisecond)
		now := int64(math.Pow(10, 2)) * millis
		if now == last {
			repeat++
		} else {
			repeat = 0
			last = now
		}
		s := fmt.Sprintf("%d", now+repeat)
		nonce, err := strconv.ParseInt(s[(len(s)-length):], 10, 64)
		if err != nil {
			return 0, err
		}
		return nonce, nil
	}
}

func QueryString(data map[string]interface{}) string {
	dataValues := url.Values{}
	for key, value := range data {
		if value != nil {
			dataValues.Add(key, fmt.Sprintf("%v", value))
		}
	}
	return dataValues.Encode()
}

func StripNils(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		if value == nil {
			delete(data, key)
		}
	}
	return data
}

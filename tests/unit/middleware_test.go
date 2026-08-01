package unit_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"makeDotApp/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestNewMaxBodyBytesMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.NewMaxBodyBytesMiddleware(10))

	r.POST("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := strings.NewReader("12345678901") // 11文字
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/octet-stream")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}

func TestNewGlobalRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.NewGlobalRateLimitMiddleware(rate.Every(time.Second), 1))

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 1回目は通る
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	// 2回目は即時制限される
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Code)
	}
}

func TestNewConcurrencyLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.NewConcurrencyLimitMiddleware(1))

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	var wg sync.WaitGroup
	statuses := make(chan int, 4)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			statuses <- w.Code
		}()
	}

	wg.Wait()
	close(statuses)

	var ok, busy int
	for s := range statuses {
		switch s {
		case http.StatusOK:
			ok++
		case http.StatusServiceUnavailable:
			busy++
		}
	}

	if ok != 1 || busy != 3 {
		t.Fatalf("expected one OK and three 503, got ok=%d busy=%d", ok, busy)
	}
}

package stats

import (
	"fmt"
	"os"
	"sync"
	"time"
	"github.com/valyala/fasthttp"
)

// Stats data structure
type Stats struct {
	mu                  sync.RWMutex
	Uptime              time.Time
	Pid                 int
	ResponseCounts      map[string]int
	TotalResponseCounts map[string]int
	TotalResponseTime   time.Time
}

// New constructs a new Stats structure
func New() *Stats {
	stats := &Stats{
		Uptime:              time.Now(),
		Pid:                 os.Getpid(),
		ResponseCounts:      map[string]int{},
		TotalResponseCounts: map[string]int{},
		TotalResponseTime:   time.Time{},
	}

	go func() {
		for {
			stats.ResetResponseCounts()

			time.Sleep(time.Second * 1)
		}
	}()

	return stats
}

// ResetResponseCounts reset the response counts
func (mw *Stats) ResetResponseCounts() {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.ResponseCounts = map[string]int{}
}

// Handler is a MiddlewareFunc makes Stats implement the Middleware interface.
func (mw *Stats) Handler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		beginning, recorder := mw.Begin(ctx)
		h(recorder)
		mw.End(beginning, recorder)
	})
}

// Begin starts a recorder
func (mw *Stats) Begin(ctx *fasthttp.RequestCtx) (time.Time, *fasthttp.RequestCtx) {
	start := time.Now()
	ctx.SetStatusCode(200)
	return start, ctx
}

// EndWithStatus closes the recorder with a specific status
func (mw *Stats) EndWithStatus(start time.Time, status int) {
	end := time.Now()

	responseTime := end.Sub(start)

	mw.mu.Lock()

	defer mw.mu.Unlock()

	statusCode := fmt.Sprintf("%d", status)

	mw.ResponseCounts[statusCode]++
	mw.TotalResponseCounts[statusCode]++
	mw.TotalResponseTime = mw.TotalResponseTime.Add(responseTime)
}

// End closes the recorder with the recorder status
func (mw *Stats) End(start time.Time, ctx *fasthttp.RequestCtx) {
	mw.EndWithStatus(start, ctx.Response.StatusCode())
}

// Data serializable structure
type Data struct {
	Pid                    int            `json:"pid"`
	UpTime                 string         `json:"uptime"`
	UpTimeSec              float64        `json:"uptime_sec"`
	Time                   string         `json:"time"`
	TimeUnix               int64          `json:"unixtime"`
	StatusCodeCount        map[string]int `json:"status_code_count"`
	TotalStatusCodeCount   map[string]int `json:"total_status_code_count"`
	Count                  int            `json:"count"`
	TotalCount             int            `json:"total_count"`
	TotalResponseTime      string         `json:"total_response_time"`
	TotalResponseTimeSec   float64        `json:"total_response_time_sec"`
	AverageResponseTime    string         `json:"average_response_time"`
	AverageResponseTimeSec float64        `json:"average_response_time_sec"`
}

// Data returns the data serializable structure
func (mw *Stats) Data() *Data {

	mw.mu.RLock()

	responseCounts := make(map[string]int, len(mw.ResponseCounts))
	totalResponseCounts := make(map[string]int, len(mw.TotalResponseCounts))

	now := time.Now()

	uptime := now.Sub(mw.Uptime)

	count := 0
	for code, current := range mw.ResponseCounts {
		responseCounts[code] = current
		count += current
	}

	totalCount := 0
	for code, count := range mw.TotalResponseCounts {
		totalResponseCounts[code] = count
		totalCount += count
	}

	totalResponseTime := mw.TotalResponseTime.Sub(time.Time{})

	averageResponseTime := time.Duration(0)
	if totalCount > 0 {
		avgNs := int64(totalResponseTime) / int64(totalCount)
		averageResponseTime = time.Duration(avgNs)
	}

	mw.mu.RUnlock()

	r := &Data{
		Pid:                    mw.Pid,
		UpTime:                 uptime.String(),
		UpTimeSec:              uptime.Seconds(),
		Time:                   now.String(),
		TimeUnix:               now.Unix(),
		StatusCodeCount:        responseCounts,
		TotalStatusCodeCount:   totalResponseCounts,
		Count:                  count,
		TotalCount:             totalCount,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
	}

	return r
}
package apis

import (
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/store"
)

const (
	DefaultRateLimitMiddlewareId       = "pbRateLimit"
	DefaultRateLimitMiddlewarePriority = -1000
)

const (
	rateLimitersStoreKey       = "__pbRateLimiters__"
	rateLimitersCronKey        = "__pbRateLimitersCleanup__"
	rateLimitersSettingsHookId = "__pbRateLimitersSettingsHook__"
)

// rateLimit defines the global rate limit middleware.
//
// This middleware is registered by default for all routes.
func rateLimit() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultRateLimitMiddlewareId,
		Priority: DefaultRateLimitMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			if skipRateLimit(e) {
				return e.Next()
			}

			rule, ok := e.App.Settings().RateLimits.FindRateLimitRule(
				defaultRateLimitLabels(e),
				defaultRateLimitAudience(e)...,
			)
			if ok {
				err := checkRateLimit(e, rule.Label+rule.Audience, rule)
				if err != nil {
					return err
				}
			}

			return e.Next()
		},
	}
}

// collectionPathRateLimit defines a rate limit middleware for the internal collection handlers.
func collectionPathRateLimit(collectionPathParam string, baseTags ...string) *hook.Handler[*core.RequestEvent] {
	if collectionPathParam == "" {
		collectionPathParam = "collection"
	}

	return &hook.Handler[*core.RequestEvent]{
		Id:       DefaultRateLimitMiddlewareId,
		Priority: DefaultRateLimitMiddlewarePriority,
		Func: func(e *core.RequestEvent) error {
			collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue(collectionPathParam))
			if err != nil {
				return e.NotFoundError("Missing or invalid collection context.", err)
			}

			if err := checkCollectionRateLimit(e, collection, baseTags...); err != nil {
				return err
			}

			return e.Next()
		},
	}
}

// checkCollectionRateLimit checks whether the current request satisfy the
// rate limit configuration for the specific collection.
//
// Each baseTags entry will be prefixed with the collection name and its wildcard variant.
func checkCollectionRateLimit(e *core.RequestEvent, collection *core.Collection, baseTags ...string) error {
	if skipRateLimit(e) {
		return nil
	}

	labels := make([]string, 0, 2+len(baseTags)*2)

	rtId := collection.Id + e.Request.Pattern

	// add first the primary labels (aka. ["collectionName:action1", "collectionName:action2"])
	for _, baseTag := range baseTags {
		rtId += baseTag
		labels = append(labels, collection.Name+":"+baseTag)
	}

	// add the wildcard labels (aka. [..., "*:action1","*:action2", "*"])
	for _, baseTag := range baseTags {
		labels = append(labels, "*:"+baseTag)
	}
	labels = append(labels, defaultRateLimitLabels(e)...)

	rule, ok := e.App.Settings().RateLimits.FindRateLimitRule(labels, defaultRateLimitAudience(e)...)
	if ok {
		return checkRateLimit(e, rtId+rule.Audience, rule)
	}

	return nil
}

// -------------------------------------------------------------------

// @todo consider exporting as helper?
//
//nolint:unused
func isClientRateLimited(e *core.RequestEvent, rtId string) bool {
	rateLimiters, ok := e.App.Store().Get(rateLimitersStoreKey).(*store.Store[string, *rateLimiter])
	if !ok || rateLimiters == nil {
		return false
	}

	rt, ok := rateLimiters.GetOk(rtId)
	if !ok || rt == nil {
		return false
	}

	client, ok := rt.getClient(e.RealIP())
	if !ok || client == nil {
		return false
	}

	return client.available <= 0 && time.Now().Unix()-client.lastConsume < client.interval
}

// @todo consider exporting as helper?
func checkRateLimit(e *core.RequestEvent, rtId string, rule core.RateLimitRule) error {
	switch rule.Audience {
	case core.RateLimitRuleAudienceAll:
		// valid for both guest and regular users
	case core.RateLimitRuleAudienceGuest:
		if e.Auth != nil {
			return nil
		}
	case core.RateLimitRuleAudienceAuth:
		if e.Auth == nil {
			return nil
		}
	}

	rateLimiters := e.App.Store().GetOrSet(rateLimitersStoreKey, func() any {
		return initRateLimitersStore(e.App)
	}).(*store.Store[string, *rateLimiter])
	if rateLimiters == nil {
		e.App.Logger().Warn("Failed to retrieve app rate limiters store")
		return nil
	}

	rt := rateLimiters.GetOrSet(rtId, func() *rateLimiter {
		return newRateLimiter(rule.MaxRequests, rule.Duration, rule.Duration+1800)
	})
	if rt == nil {
		e.App.Logger().Warn("Failed to retrieve app rate limiter", "id", rtId)
		return nil
	}

	key := e.RealIP()
	if key == "" {
		e.App.Logger().Warn("Empty rate limit client key")
		return nil
	}

	if !rt.isAllowed(key) {
		return e.TooManyRequestsError("", nil)
	}

	return nil
}

func skipRateLimit(e *core.RequestEvent) bool {
	return !e.App.Settings().RateLimits.Enabled || e.HasSuperuserAuth()
}

var defaultAuthAudience = []string{core.RateLimitRuleAudienceAll, core.RateLimitRuleAudienceAuth}
var defaultGuestAudience = []string{core.RateLimitRuleAudienceAll, core.RateLimitRuleAudienceGuest}

func defaultRateLimitAudience(e *core.RequestEvent) []string {
	if e.Auth != nil {
		return defaultAuthAudience
	}

	return defaultGuestAudience
}

func defaultRateLimitLabels(e *core.RequestEvent) []string {
	return []string{e.Request.Method + " " + e.Request.URL.Path, e.Request.URL.Path}
}

func destroyRateLimitersStore(app core.App) {
	app.OnSettingsReload().Unbind(rateLimitersSettingsHookId)
	app.Cron().Remove(rateLimitersCronKey)
	app.Store().Remove(rateLimitersStoreKey)
}

func initRateLimitersStore(app core.App) *store.Store[string, *rateLimiter] {
	app.Cron().Add(rateLimitersCronKey, "2 * * * *", func() { // offset a little since too many cleanup tasks execute at 00
		limitersStore, ok := app.Store().Get(rateLimitersStoreKey).(*store.Store[string, *rateLimiter])
		if !ok {
			return
		}
		limiters := limitersStore.GetAll()
		for _, limiter := range limiters {
			limiter.clean()
		}
	})

	app.OnSettingsReload().Bind(&hook.Handler[*core.SettingsReloadEvent]{
		Id: rateLimitersSettingsHookId,
		Func: func(e *core.SettingsReloadEvent) error {
			err := e.Next()
			if err != nil {
				return err
			}

			// reset
			destroyRateLimitersStore(e.App)

			return nil
		},
	})

	return store.New[string, *rateLimiter](nil)
}

func newRateLimiter(maxAllowed int, intervalInSec int64, minDeleteIntervalInSec int64) *rateLimiter {
	return &rateLimiter{
		maxAllowed:        maxAllowed,
		interval:          intervalInSec,
		minDeleteInterval: minDeleteIntervalInSec,
		clients:           map[string]*fixedWindow{},
	}
}

type rateLimiter struct {
	clients map[string]*fixedWindow

	maxAllowed        int
	interval          int64
	minDeleteInterval int64
	totalDeleted      int64

	sync.RWMutex
}

//nolint:unused
func (rt *rateLimiter) getClient(key string) (*fixedWindow, bool) {
	rt.RLock()
	client, ok := rt.clients[key]
	rt.RUnlock()

	return client, ok
}

func (rt *rateLimiter) isAllowed(key string) bool {
	// lock only reads to minimize locks contention
	rt.RLock()
	client, ok := rt.clients[key]
	rt.RUnlock()

	if !ok {
		rt.Lock()
		// check again in case the client was added by another request
		client, ok = rt.clients[key]
		if !ok {
			client = newFixedWindow(rt.maxAllowed, rt.interval)
			rt.clients[key] = client
		}
		rt.Unlock()
	}

	return client.consume()
}

func (rt *rateLimiter) clean() {
	rt.Lock()
	defer rt.Unlock()

	nowUnix := time.Now().Unix()

	for k, client := range rt.clients {
		if client.hasExpired(nowUnix, rt.minDeleteInterval) {
			delete(rt.clients, k)
			rt.totalDeleted++
		}
	}

	// "shrink" the map if too may items were deleted
	//
	// @todo remove after https://github.com/golang/go/issues/20135
	if rt.totalDeleted >= 300 {
		shrunk := make(map[string]*fixedWindow, len(rt.clients))
		for k, v := range rt.clients {
			shrunk[k] = v
		}
		rt.clients = shrunk
		rt.totalDeleted = 0
	}
}

func newFixedWindow(maxAllowed int, intervalInSec int64) *fixedWindow {
	return &fixedWindow{
		maxAllowed: maxAllowed,
		interval:   intervalInSec,
	}
}

type fixedWindow struct {
	// use plain Mutex instead of RWMutex since the operations are expected
	// to be mostly writes (e.g. consume()) and it should perform better
	sync.Mutex

	maxAllowed  int   // the max allowed tokens per interval
	available   int   // the total available tokens
	interval    int64 // in seconds
	lastConsume int64 // the time of the last consume
}

// hasExpired checks whether it has been at least minElapsed seconds since the lastConsume time.
// (usually used to perform periodic cleanup of staled instances).
func (l *fixedWindow) hasExpired(relativeNow int64, minElapsed int64) bool {
	l.Lock()
	defer l.Unlock()

	return relativeNow-l.lastConsume > minElapsed
}

// consume decrease the current window allowance with 1 (if not exhausted already).
//
// It returns false if the allowance has been already exhausted and the user
// has to wait until it resets back to its maxAllowed value.
func (l *fixedWindow) consume() bool {
	l.Lock()
	defer l.Unlock()

	nowUnix := time.Now().Unix()

	// reset consumed counter
	if nowUnix-l.lastConsume >= l.interval {
		l.available = l.maxAllowed
	}

	if l.available > 0 {
		l.available--
		l.lastConsume = nowUnix

		return true
	}

	return false
}

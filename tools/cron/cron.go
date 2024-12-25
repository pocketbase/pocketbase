// Package cron implements a crontab-like service to execute and schedule
// repeative tasks/jobs.
//
// Example:
//
//	c := cron.New()
//	c.MustAdd("dailyReport", "0 0 * * *", func() { ... })
//	c.Start()
package cron

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

// Cron is a crontab-like struct for tasks/jobs scheduling.
type Cron struct {
	timezone   *time.Location
	ticker     *time.Ticker
	startTimer *time.Timer
	tickerDone chan bool
	jobs       []*Job
	interval   time.Duration
	mux        sync.RWMutex
}

// New create a new Cron struct with default tick interval of 1 minute
// and timezone in UTC.
//
// You can change the default tick interval with Cron.SetInterval().
// You can change the default timezone with Cron.SetTimezone().
func New() *Cron {
	return &Cron{
		interval:   1 * time.Minute,
		timezone:   time.UTC,
		jobs:       []*Job{},
		tickerDone: make(chan bool),
	}
}

// SetInterval changes the current cron tick interval
// (it usually should be >= 1 minute).
func (c *Cron) SetInterval(d time.Duration) {
	// update interval
	c.mux.Lock()
	wasStarted := c.ticker != nil
	c.interval = d
	c.mux.Unlock()

	// restart the ticker
	if wasStarted {
		c.Start()
	}
}

// SetTimezone changes the current cron tick timezone.
func (c *Cron) SetTimezone(l *time.Location) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.timezone = l
}

// MustAdd is similar to Add() but panic on failure.
func (c *Cron) MustAdd(jobId string, cronExpr string, run func()) {
	if err := c.Add(jobId, cronExpr, run); err != nil {
		panic(err)
	}
}

// Add registers a single cron job.
//
// If there is already a job with the provided id, then the old job
// will be replaced with the new one.
//
// cronExpr is a regular cron expression, eg. "0 */3 * * *" (aka. at minute 0 past every 3rd hour).
// Check cron.NewSchedule() for the supported tokens.
func (c *Cron) Add(jobId string, cronExpr string, fn func()) error {
	if fn == nil {
		return errors.New("failed to add new cron job: fn must be non-nil function")
	}

	schedule, err := NewSchedule(cronExpr)
	if err != nil {
		return fmt.Errorf("failed to add new cron job: %w", err)
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// remove previous (if any)
	c.jobs = slices.DeleteFunc(c.jobs, func(j *Job) bool {
		return j.Id() == jobId
	})

	// add new
	c.jobs = append(c.jobs, &Job{
		id:       jobId,
		fn:       fn,
		schedule: schedule,
	})

	return nil
}

// Remove removes a single cron job by its id.
func (c *Cron) Remove(jobId string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.jobs == nil {
		return // nothing to remove
	}

	c.jobs = slices.DeleteFunc(c.jobs, func(j *Job) bool {
		return j.Id() == jobId
	})
}

// RemoveAll removes all registered cron jobs.
func (c *Cron) RemoveAll() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.jobs = []*Job{}
}

// Total returns the current total number of registered cron jobs.
func (c *Cron) Total() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return len(c.jobs)
}

// Jobs returns a shallow copy of the currently registered cron jobs.
func (c *Cron) Jobs() []*Job {
	c.mux.RLock()
	defer c.mux.RUnlock()

	copy := make([]*Job, len(c.jobs))
	for i, j := range c.jobs {
		copy[i] = j
	}

	return copy
}

// Stop stops the current cron ticker (if not already).
//
// You can resume the ticker by calling Start().
func (c *Cron) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.startTimer != nil {
		c.startTimer.Stop()
		c.startTimer = nil
	}

	if c.ticker == nil {
		return // already stopped
	}

	c.tickerDone <- true
	c.ticker.Stop()
	c.ticker = nil
}

// Start starts the cron ticker.
//
// Calling Start() on already started cron will restart the ticker.
func (c *Cron) Start() {
	c.Stop()

	// delay the ticker to start at 00 of 1 c.interval duration
	now := time.Now()
	next := now.Add(c.interval).Truncate(c.interval)
	delay := next.Sub(now)

	c.mux.Lock()
	c.startTimer = time.AfterFunc(delay, func() {
		c.mux.Lock()
		c.ticker = time.NewTicker(c.interval)
		c.mux.Unlock()

		// run immediately at 00
		c.runDue(time.Now())

		// run after each tick
		go func() {
			for {
				select {
				case <-c.tickerDone:
					return
				case t := <-c.ticker.C:
					c.runDue(t)
				}
			}
		}()
	})
	c.mux.Unlock()
}

// HasStarted checks whether the current Cron ticker has been started.
func (c *Cron) HasStarted() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.ticker != nil
}

// runDue runs all registered jobs that are scheduled for the provided time.
func (c *Cron) runDue(t time.Time) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	moment := NewMoment(t.In(c.timezone))

	for _, j := range c.jobs {
		if j.schedule.IsDue(moment) {
			go j.Run()
		}
	}
}

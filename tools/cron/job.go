package cron

import "encoding/json"

// Job defines a single registered cron job.
type Job struct {
	fn       func()
	schedule *Schedule
	id       string
	paused   bool
}

// Id returns the cron job id.
func (j *Job) Id() string {
	return j.id
}

// Expression returns the plain cron job schedule expression.
func (j *Job) Expression() string {
	return j.schedule.rawExpr
}

// Run runs the cron job function.
func (j *Job) Run() {
	if j.fn != nil {
		j.fn()
	}
}

// Pause pauses the cron job.
func (j *Job) Pause() {
	j.paused = true
}

// Resume resumes the cron job.
func (j *Job) Resume() {
	j.paused = false
}

// IsPaused returns whether the cron job is currently paused.
func (j *Job) IsPaused() bool {
	return j.paused
}

// MarshalJSON implements [json.Marshaler] and export the current
// jobs data into valid JSON.
func (j Job) MarshalJSON() ([]byte, error) {
	plain := struct {
		Id         string `json:"id"`
		Expression string `json:"expression"`
		Paused     bool   `json:"paused"`
	}{
		Id:         j.Id(),
		Expression: j.Expression(),
		Paused:     j.IsPaused(),
	}

	return json.Marshal(plain)
}

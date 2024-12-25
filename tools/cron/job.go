package cron

import "encoding/json"

// Job defines a single registered cron job.
type Job struct {
	fn       func()
	schedule *Schedule
	id       string
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

// MarshalJSON implements [json.Marshaler] and export the current
// jobs data into valid JSON.
func (j Job) MarshalJSON() ([]byte, error) {
	plain := struct {
		Id         string `json:"id"`
		Expression string `json:"expression"`
	}{
		Id:         j.Id(),
		Expression: j.Expression(),
	}

	return json.Marshal(plain)
}

package apis

import (
	"net/http"
	"slices"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/routine"
)

// bindCronApi registers the crons api endpoint.
func bindCronApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/crons").Bind(RequireSuperuserAuth())
	subGroup.GET("", cronsList)
	subGroup.POST("/{id}", cronRun)
	subGroup.PUT("/{id}/pause", cronPause)
	subGroup.PUT("/{id}/resume", cronResume)
}

func cronsList(e *core.RequestEvent) error {
	jobs := e.App.Cron().Jobs()

	slices.SortStableFunc(jobs, func(a, b *cron.Job) int {
		if strings.HasPrefix(a.Id(), "__pb") {
			return 1
		}
		if strings.HasPrefix(b.Id(), "__pb") {
			return -1
		}
		return strings.Compare(a.Id(), b.Id())
	})

	return e.JSON(http.StatusOK, jobs)
}

func cronRun(e *core.RequestEvent) error {
	cronId := e.Request.PathValue("id")

	var foundJob *cron.Job

	jobs := e.App.Cron().Jobs()
	for _, j := range jobs {
		if j.Id() == cronId {
			foundJob = j
			break
		}
	}

	if foundJob == nil {
		return e.NotFoundError("Missing or invalid cron job", nil)
	}

	routine.FireAndForget(func() {
		foundJob.Run()
	})

	return e.NoContent(http.StatusNoContent)
}

func cronPause(e *core.RequestEvent) error {
	cronId := e.Request.PathValue("id")

	err := e.App.Cron().PauseJob(cronId)
	if err != nil {
		if strings.Contains(err.Error(), "system jobs cannot be paused") {
			return e.BadRequestError("System jobs cannot be paused", err)
		}
		if strings.Contains(err.Error(), "job not found") {
			return e.NotFoundError("Missing or invalid cron job", err)
		}
		return e.InternalServerError("Failed to pause cron job", err)
	}

	return e.NoContent(http.StatusNoContent)
}

func cronResume(e *core.RequestEvent) error {
	cronId := e.Request.PathValue("id")

	err := e.App.Cron().ResumeJob(cronId)
	if err != nil {
		if strings.Contains(err.Error(), "system jobs cannot be resumed") {
			return e.BadRequestError("System jobs cannot be resumed", err)
		}
		if strings.Contains(err.Error(), "job not found") {
			return e.NotFoundError("Missing or invalid cron job", err)
		}
		return e.InternalServerError("Failed to resume cron job", err)
	}

	return e.NoContent(http.StatusNoContent)
}

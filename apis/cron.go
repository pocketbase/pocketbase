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

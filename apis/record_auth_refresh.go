package apis

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func recordAuthRefresh(e *core.RequestEvent) error {
	record := e.Auth
	if record == nil {
		return e.NotFoundError("Missing auth record context.", nil)
	}

	currentToken := getAuthTokenFromRequest(e)
	claims, _ := security.ParseUnverifiedJWT(currentToken)
	if v, ok := claims[core.TokenClaimRefreshable]; !ok || !cast.ToBool(v) {
		return e.ForbiddenError("The current auth token is not refreshable.", nil)
	}

	event := new(core.RecordAuthRefreshRequestEvent)
	event.RequestEvent = e
	event.Collection = record.Collection()
	event.Record = record

	return e.App.OnRecordAuthRefreshRequest().Trigger(event, func(e *core.RecordAuthRefreshRequestEvent) error {
		return RecordAuthResponse(e.RequestEvent, e.Record, "", nil)
	})
}

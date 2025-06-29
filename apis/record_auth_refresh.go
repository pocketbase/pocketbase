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

	event := new(core.RecordAuthRefreshRequestEvent)
	event.RequestEvent = e
	event.Collection = record.Collection()
	event.Record = record

	return e.App.OnRecordAuthRefreshRequest().Trigger(event, func(e *core.RecordAuthRefreshRequestEvent) error {
		token := getAuthTokenFromRequest(e.RequestEvent)

		// skip token renewal if the token's payload doesn't explicitly allow it (e.g. impersonate tokens)
		claims, _ := security.ParseUnverifiedJWT(token) //
		if v, ok := claims[core.TokenClaimRefreshable]; ok && cast.ToBool(v) {
			var tokenErr error
			token, tokenErr = e.Record.NewAuthToken()
			if tokenErr != nil {
				return e.InternalServerError("Failed to refresh auth token.", tokenErr)
			}
		}

		return recordAuthResponse(e.RequestEvent, e.Record, token, "", nil)
	})
}

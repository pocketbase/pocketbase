package apis

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// bindRecordAuthApi registers the auth record api endpoints and
// the corresponding handlers.
func bindRecordAuthApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	// global oauth2 subscription redirect handler
	rg.GET("/oauth2-redirect", oauth2SubscriptionRedirect).Bind(
		SkipSuccessActivityLog(), // skip success log as it could contain sensitive information in the url
	)
	// add again as POST in case of response_mode=form_post
	rg.POST("/oauth2-redirect", oauth2SubscriptionRedirect).Bind(
		SkipSuccessActivityLog(), // skip success log as it could contain sensitive information in the url
	)

	sub := rg.Group("/collections/{collection}")

	sub.GET("/auth-methods", recordAuthMethods).Bind(
		collectionPathRateLimit("", "listAuthMethods"),
	)

	sub.POST("/auth-refresh", recordAuthRefresh).Bind(
		collectionPathRateLimit("", "authRefresh"),
		RequireSameCollectionContextAuth(""),
	)

	sub.POST("/auth-with-password", recordAuthWithPassword).Bind(
		collectionPathRateLimit("", "authWithPassword", "auth"),
	)

	sub.POST("/auth-with-oauth2", recordAuthWithOAuth2).Bind(
		collectionPathRateLimit("", "authWithOAuth2", "auth"),
	)

	sub.POST("/request-otp", recordRequestOTP).Bind(
		collectionPathRateLimit("", "requestOTP"),
	)
	sub.POST("/auth-with-otp", recordAuthWithOTP).Bind(
		collectionPathRateLimit("", "authWithOTP", "auth"),
	)

	sub.POST("/request-password-reset", recordRequestPasswordReset).Bind(
		collectionPathRateLimit("", "requestPasswordReset"),
	)
	sub.POST("/confirm-password-reset", recordConfirmPasswordReset).Bind(
		collectionPathRateLimit("", "confirmPasswordReset"),
	)

	sub.POST("/request-verification", recordRequestVerification).Bind(
		collectionPathRateLimit("", "requestVerification"),
	)
	sub.POST("/confirm-verification", recordConfirmVerification).Bind(
		collectionPathRateLimit("", "confirmVerification"),
	)

	sub.POST("/request-email-change", recordRequestEmailChange).Bind(
		collectionPathRateLimit("", "requestEmailChange"),
		RequireSameCollectionContextAuth(""),
	)
	sub.POST("/confirm-email-change", recordConfirmEmailChange).Bind(
		collectionPathRateLimit("", "confirmEmailChange"),
	)

	sub.POST("/impersonate/{id}", recordAuthImpersonate).Bind(RequireSuperuserAuth())
}

func findAuthCollection(e *core.RequestEvent) (*core.Collection, error) {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))

	if err != nil || !collection.IsAuth() {
		return nil, e.NotFoundError("Missing or invalid auth collection context.", err)
	}

	return collection, nil
}

package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		return createWebAuthnCredentialsCollection(txApp)
	}, func(txApp core.App) error {
		col, err := txApp.FindCollectionByNameOrId(core.CollectionNameWebAuthnCredentials)
		if err != nil {
			return nil // already deleted or doesn't exist
		}
		return txApp.Delete(col)
	})
}

func createWebAuthnCredentialsCollection(txApp core.App) error {
	col := core.NewBaseCollection(core.CollectionNameWebAuthnCredentials)
	col.System = true

	ownerRule := "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
	col.ListRule = types.Pointer(ownerRule)
	col.ViewRule = types.Pointer(ownerRule)
	col.DeleteRule = types.Pointer(ownerRule)

	col.Fields.Add(&core.TextField{
		Name:     "collectionRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "recordRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "credentialId",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "publicKey",
		System:   true,
		Required: true,
		Hidden:   true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "attestationType",
		System:   true,
		Required: false,
	})
	col.Fields.Add(&core.JSONField{
		Name:   "transport",
		System: true,
	})
	col.Fields.Add(&core.JSONField{
		Name:   "flags",
		System: true,
	})
	col.Fields.Add(&core.NumberField{
		Name:   "signCount",
		System: true,
	})
	col.Fields.Add(&core.TextField{
		Name:   "name",
		System: true,
	})
	col.Fields.Add(&core.TextField{
		Name:   "aaguid",
		System: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	col.AddIndex("idx_webauthnCredentials_collectionRef_recordRef", false, "collectionRef, recordRef", "")
	col.AddIndex("idx_webauthnCredentials_credentialId", true, "credentialId", "")

	return txApp.Save(col)
}

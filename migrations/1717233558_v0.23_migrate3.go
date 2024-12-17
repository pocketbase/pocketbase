package migrations

import (
	"hash/crc32"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// note: this migration will be deleted in future version

func collectionIdChecksum(typ, name string) string {
	return "pbc_" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(typ+name))), 10)
}

func fieldIdChecksum(typ, name string) string {
	return typ + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(name))), 10)
}

// normalize system collection and field ids
func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		collections := []*core.Collection{}
		err := txApp.CollectionQuery().
			AndWhere(dbx.In(
				"name",
				core.CollectionNameMFAs,
				core.CollectionNameOTPs,
				core.CollectionNameExternalAuths,
				core.CollectionNameAuthOrigins,
				core.CollectionNameSuperusers,
			)).
			All(&collections)
		if err != nil {
			return err
		}

		for _, c := range collections {
			var needUpdate bool

			references, err := txApp.FindCollectionReferences(c, c.Id)
			if err != nil {
				return err
			}

			authOrigins, err := txApp.FindAllAuthOriginsByCollection(c)
			if err != nil {
				return err
			}

			mfas, err := txApp.FindAllMFAsByCollection(c)
			if err != nil {
				return err
			}

			otps, err := txApp.FindAllOTPsByCollection(c)
			if err != nil {
				return err
			}

			originalId := c.Id

			// normalize collection id
			if checksum := collectionIdChecksum(c.Type, c.Name); c.Id != checksum {
				c.Id = checksum
				needUpdate = true
			}

			// normalize system fields
			for _, f := range c.Fields {
				if !f.GetSystem() {
					continue
				}

				if checksum := fieldIdChecksum(f.Type(), f.GetName()); f.GetId() != checksum {
					f.SetId(checksum)
					needUpdate = true
				}
			}

			if !needUpdate {
				continue
			}

			rawExport, err := c.DBExport(txApp)
			if err != nil {
				return err
			}

			_, err = txApp.DB().Update("_collections", rawExport, dbx.HashExp{"id": originalId}).Execute()
			if err != nil {
				return err
			}

			// make sure that the cached collection id is also updated
			cached, err := txApp.FindCachedCollectionByNameOrId(c.Name)
			if err != nil {
				return err
			}
			cached.Id = c.Id

			// update collection references
			for refCollection, fields := range references {
				for _, f := range fields {
					relationField, ok := f.(*core.RelationField)
					if !ok || relationField.CollectionId == originalId {
						continue
					}

					relationField.CollectionId = c.Id
				}
				if err = txApp.SaveNoValidate(refCollection); err != nil {
					return err
				}
			}

			// update mfas references
			for _, item := range mfas {
				item.SetCollectionRef(c.Id)
				if err = txApp.SaveNoValidate(item); err != nil {
					return err
				}
			}

			// update otps references
			for _, item := range otps {
				item.SetCollectionRef(c.Id)
				if err = txApp.SaveNoValidate(item); err != nil {
					return err
				}
			}

			// update authOrigins references
			for _, item := range authOrigins {
				item.SetCollectionRef(c.Id)
				if err = txApp.SaveNoValidate(item); err != nil {
					return err
				}
			}
		}

		return nil
	}, nil)
}

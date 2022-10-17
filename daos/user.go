package daos

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
)

// UserQuery returns a new User model select query.
func (dao *Dao) UserQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.User{})
}

// LoadProfile loads the profile record associated to the provided user.
func (dao *Dao) LoadProfile(user *models.User) error {
	collection, err := dao.FindCollectionByNameOrId(models.ProfileCollectionName)
	if err != nil {
		return err
	}

	profile, err := dao.FindFirstRecordByData(collection, models.ProfileCollectionUserFieldName, user.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	user.Profile = profile

	return nil
}

// LoadProfiles loads the profile records associated to the provided users list.
func (dao *Dao) LoadProfiles(users []*models.User) error {
	collection, err := dao.FindCollectionByNameOrId(models.ProfileCollectionName)
	if err != nil {
		return err
	}

	// extract user ids
	ids := make([]string, len(users))
	usersMap := map[string]*models.User{}
	for i, user := range users {
		ids[i] = user.Id
		usersMap[user.Id] = user
	}

	profiles, err := dao.FindRecordsByExpr(collection, dbx.HashExp{
		models.ProfileCollectionUserFieldName: list.ToInterfaceSlice(ids),
	})
	if err != nil {
		return err
	}

	// populate each user.Profile member
	for _, profile := range profiles {
		userId := profile.GetStringDataValue(models.ProfileCollectionUserFieldName)
		user, ok := usersMap[userId]
		if !ok {
			continue
		}
		user.Profile = profile
	}

	return nil
}

// FindUserById finds a single User model by its id.
//
// This method also auto loads the related user profile record
// into the found model.
func (dao *Dao) FindUserById(id string) (*models.User, error) {
	model := &models.User{}

	err := dao.UserQuery().
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	// try to load the user profile (if exist)
	if err := dao.LoadProfile(model); err != nil {
		log.Println(err)
	}

	return model, nil
}

// FindUserByEmail finds a single User model by its non-empty email address.
//
// This method also auto loads the related user profile record
// into the found model.
func (dao *Dao) FindUserByEmail(email string) (*models.User, error) {
	model := &models.User{}

	err := dao.UserQuery().
		AndWhere(dbx.Not(dbx.HashExp{"email": ""})).
		AndWhere(dbx.HashExp{"email": email}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	// try to load the user profile (if exist)
	if err := dao.LoadProfile(model); err != nil {
		log.Println(err)
	}

	return model, nil
}

// FindUserByToken finds the user associated with the provided JWT token.
// Returns an error if the JWT token is invalid or expired.
//
// This method also auto loads the related user profile record
// into the found model.
func (dao *Dao) FindUserByToken(token string, baseTokenKey string) (*models.User, error) {
	unverifiedClaims, err := security.ParseUnverifiedJWT(token)
	if err != nil {
		return nil, err
	}

	// check required claims
	id, _ := unverifiedClaims["id"].(string)
	if id == "" {
		return nil, errors.New("Missing or invalid token claims.")
	}

	user, err := dao.FindUserById(id)
	if err != nil || user == nil {
		return nil, err
	}

	verificationKey := user.TokenKey + baseTokenKey

	// verify token signature
	if _, err := security.ParseJWT(token, verificationKey); err != nil {
		return nil, err
	}

	return user, nil
}

// IsUserEmailUnique checks if the provided email address is not
// already in use by other users.
func (dao *Dao) IsUserEmailUnique(email string, excludeId string) bool {
	if email == "" {
		return false
	}

	var exists bool
	err := dao.UserQuery().
		Select("count(*)").
		AndWhere(dbx.Not(dbx.HashExp{"id": excludeId})).
		AndWhere(dbx.HashExp{"email": email}).
		Limit(1).
		Row(&exists)

	return err == nil && !exists
}

// DeleteUser deletes the provided User model.
//
// This method will also cascade the delete operation to all
// Record models that references the provided User model
// (delete or set to NULL, depending on the related user shema field settings).
//
// The delete operation may fail if the user is part of a required
// reference in another Record model (aka. cannot be deleted or set to NULL).
func (dao *Dao) DeleteUser(user *models.User) error {
	// fetch related records
	// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
	relatedRecords, err := dao.FindUserRelatedRecords(user)
	if err != nil {
		return err
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// check if related records has to be deleted (if `CascadeDelete` is set)
		// OR
		// just unset the user related fields (if they are not required)
		// -----------------------------------------------------------
	recordsLoop:
		for _, record := range relatedRecords {
			var needSave bool

			for _, field := range record.Collection().Schema.Fields() {
				if field.Type != schema.FieldTypeUser {
					continue // not a user field
				}

				ids := record.GetStringSliceDataValue(field.Name)

				// unset the user id
				for i := len(ids) - 1; i >= 0; i-- {
					if ids[i] == user.Id {
						ids = append(ids[:i], ids[i+1:]...)
						break
					}
				}

				options, _ := field.Options.(*schema.UserOptions)

				// cascade delete
				// (only if there are no other user references in case of multiple select)
				if options.CascadeDelete && len(ids) == 0 {
					if err := txDao.DeleteRecord(record); err != nil {
						return err
					}
					// no need to further iterate the user fields (the record is deleted)
					continue recordsLoop
				}

				if field.Required && len(ids) == 0 {
					return fmt.Errorf("Failed delete the user because a record exist with required user reference to the current model (%q, %q).", record.Id, record.Collection().Name)
				}

				// apply the reference changes
				record.SetDataValue(field.Name, field.PrepareValue(ids))
				needSave = true
			}

			if needSave {
				if err := txDao.SaveRecord(record); err != nil {
					return err
				}
			}
		}
		// -----------------------------------------------------------

		return txDao.Delete(user)
	})
}

// SaveUser upserts the provided User model.
//
// An empty profile record will be created if the user
// doesn't have a profile record set yet.
func (dao *Dao) SaveUser(user *models.User) error {
	profileCollection, err := dao.FindCollectionByNameOrId(models.ProfileCollectionName)
	if err != nil {
		return err
	}

	// fetch the related user profile record (if exist)
	var userProfile *models.Record
	if user.HasId() {
		userProfile, _ = dao.FindFirstRecordByData(
			profileCollection,
			models.ProfileCollectionUserFieldName,
			user.Id,
		)
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		if err := txDao.Save(user); err != nil {
			return err
		}

		// create default/empty profile record if doesn't exist
		if userProfile == nil {
			userProfile = models.NewRecord(profileCollection)
			userProfile.SetDataValue(models.ProfileCollectionUserFieldName, user.Id)
			if err := txDao.Save(userProfile); err != nil {
				return err
			}
			user.Profile = userProfile
		}

		return nil
	})
}

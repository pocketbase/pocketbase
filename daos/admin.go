package daos

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
)

// AdminQuery returns a new Admin select query.
func (dao *Dao) AdminQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Admin{})
}

// FindAdminById finds the admin with the provided id.
func (dao *Dao) FindAdminById(id string) (*models.Admin, error) {
	model := &models.Admin{}

	err := dao.AdminQuery().
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindAdminByEmail finds the admin with the provided email address.
func (dao *Dao) FindAdminByEmail(email string) (*models.Admin, error) {
	model := &models.Admin{}

	err := dao.AdminQuery().
		AndWhere(dbx.HashExp{"email": email}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindAdminByToken finds the admin associated with the provided JWT.
//
// Returns an error if the JWT is invalid or expired.
func (dao *Dao) FindAdminByToken(token string, baseTokenKey string) (*models.Admin, error) {
	// @todo consider caching the unverified claims
	unverifiedClaims, err := security.ParseUnverifiedJWT(token)
	if err != nil {
		return nil, err
	}

	// check required claims
	id, _ := unverifiedClaims["id"].(string)
	if id == "" {
		return nil, errors.New("missing or invalid token claims")
	}

	admin, err := dao.FindAdminById(id)
	if err != nil || admin == nil {
		return nil, err
	}

	verificationKey := admin.TokenKey + baseTokenKey

	// verify token signature
	if _, err := security.ParseJWT(token, verificationKey); err != nil {
		return nil, err
	}

	return admin, nil
}

// TotalAdmins returns the number of existing admin records.
func (dao *Dao) TotalAdmins() (int, error) {
	var total int

	err := dao.AdminQuery().Select("count(*)").Row(&total)

	return total, err
}

// IsAdminEmailUnique checks if the provided email address is not
// already in use by other admins.
func (dao *Dao) IsAdminEmailUnique(email string, excludeIds ...string) bool {
	if email == "" {
		return false
	}

	query := dao.AdminQuery().Select("count(*)").
		AndWhere(dbx.HashExp{"email": email}).
		Limit(1)

	if uniqueExcludeIds := list.NonzeroUniques(excludeIds); len(uniqueExcludeIds) > 0 {
		query.AndWhere(dbx.NotIn("id", list.ToInterfaceSlice(uniqueExcludeIds)...))
	}

	var exists bool

	return query.Row(&exists) == nil && !exists
}

// DeleteAdmin deletes the provided Admin model.
//
// Returns an error if there is only 1 admin.
func (dao *Dao) DeleteAdmin(admin *models.Admin) error {
	total, err := dao.TotalAdmins()
	if err != nil {
		return err
	}

	if total == 1 {
		return errors.New("you cannot delete the only existing admin")
	}

	return dao.Delete(admin)
}

// SaveAdmin upserts the provided Admin model.
func (dao *Dao) SaveAdmin(admin *models.Admin) error {
	return dao.Save(admin)
}

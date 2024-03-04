package models

var (
	_ Model = (*Address)(nil)
)

type Address struct {
	BaseModel

	AdminID int    `db:"adminID" json:"adminID"`
	Street  string `db:"street" json:"street"`
	City    string `db:"city" json:"city"`
	State   string `db:"state" json:"state"`
	ZipCode string `db:"zipCode" json:"zipCode"`
	Country string `db:"country" json:"country"`
}

func (m *Address) TableName() string {
	return "_addresses"
}

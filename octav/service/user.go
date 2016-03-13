package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/tools"
)

func (v *User) populateRowForCreate(vdb *db.User, payload CreateUserRequest) error {
	vdb.EID = tools.UUID()

	vdb.FirstName = payload.FirstName
	vdb.LastName = payload.LastName
	vdb.Nickname = payload.Nickname
	vdb.Email = payload.Email
	vdb.TshirtSize = payload.TshirtSize

	return nil
}

func (v *User) populateRowForUpdate(vdb *db.User, payload UpdateUserRequest) error {
	return nil
}

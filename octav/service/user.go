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
	if payload.FirstName.Valid() {
		vdb.FirstName = payload.FirstName.String
	}

	if payload.LastName.Valid() {
		vdb.LastName = payload.LastName.String
	}

	if payload.Nickname.Valid() {
		vdb.Nickname = payload.Nickname.String
	}

	if payload.Email.Valid() {
		vdb.Email = payload.Email.String
	}

	if payload.TshirtSize.Valid() {
		vdb.TshirtSize = payload.TshirtSize.String
	}

	return nil
}

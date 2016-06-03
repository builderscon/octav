package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Client) populateRowForCreate(vdb *db.Client, payload model.CreateClientRequest) error {
	vdb.EID = payload.ID
	vdb.Secret = tools.UUID() // XXX Change this to generate a "password"
	vdb.Name = payload.Name
	return nil
}

func (v *Client) populateRowForUpdate(vdb *db.Client, payload model.UpdateClientRequest) error {
	vdb.Secret = payload.Secret
	vdb.Name = payload.Name
	return nil
}

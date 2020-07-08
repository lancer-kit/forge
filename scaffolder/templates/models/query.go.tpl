package models

import "github.com/lancer-kit/armory/db"

// QI is a top level interface for interaction with database.
type QI interface {
	db.Transactional

	/*
      TODO: Here is needed to add the query interfaces
    */
}

// Q implementation of the `QI` interface.
type Q struct {
	*db.SQLConn
}

// NewQ returns initialized instance of the `QI`.
func NewQ(dbConn *db.SQLConn) *Q {
	if dbConn == nil {
		dbConn = db.GetConnector()
	}
	return &Q{
		SQLConn: dbConn,
	}
}

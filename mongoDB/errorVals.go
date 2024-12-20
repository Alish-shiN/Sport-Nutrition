package mongoDB

import "errors"

var (
	ErrConnectionToMongo  = errors.New("error connecting to MongoDB")
	ErrDuringVerification = errors.New("connection verification error")
)

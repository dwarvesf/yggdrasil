package model

import uuid "github.com/satori/go.uuid"

var emptyUUID, _ = uuid.FromString("00000000-0000-0000-0000-000000000000")

//IsZero check uuid is not zero
func IsZero(u uuid.UUID) bool {
	return u == emptyUUID
}

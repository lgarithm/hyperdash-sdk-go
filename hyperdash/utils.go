package hyperdash

import uuid "github.com/satori/go.uuid"

func genUUID() string {
	u := uuid.NewV4()
	return u.String()
}

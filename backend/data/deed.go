package data

import "gorm.io/gorm"

type DeedOwnerKind string

const (
	DeedOwnerPerson DeedOwnerKind = "person"
)

type DeedObjectKind string

const (
	DeedObjectShip DeedObjectKind = "ship"
)

type Deed struct {
	gorm.Model

	OwnerID   uint
	OwnerType DeedOwnerKind

	ObjectID   uint
	ObjectType DeedObjectKind
}

package data

import "gorm.io/gorm"

type ContractPartyKind string

const (
	CrewContractParty   ContractPartyKind = "crew"
	PersonContractParty ContractPartyKind = "person"
)

type ContractTermsKind string

const (
	EmploymentContractTerms ContractTermsKind = "employment"
)

type Contract struct {
	gorm.Model

	OfferorKind ContractPartyKind
	OfferorID   uint

	OffereeKind ContractPartyKind
	OffereeID   uint // Optional, e.g. in the case of a bounty, or a letter of marque

	TermsKind ContractTermsKind
	TermsID   uint
}

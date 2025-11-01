package inbound

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
	"gorm.io/gorm"
)

type HireCrewReq struct {
	People []Person `json:"people"`
}

type HireCrewResp struct {
}

func HireCrew(r HireCrewReq, conn Connection) (*HireCrewResp, error) {

	slog.Info("HireCrew")

	ship, err := conn.Context().Ship()
	if err != nil {
		return nil, fmt.Errorf("failed to get ship: %w", err)
	}

	crew, err := ship.Crew()
	if err != nil {
		return nil, fmt.Errorf("failed to get crew: %w", err)
	}

	for _, p := range r.People {
		err = db.Conn().Transaction(func(tx *gorm.DB) error {
			personData := &data.Person{
				FirstName: p.FirstName,
				LastName:  p.LastName,
				NickName:  p.NickName,
				Age:       uint(p.Age),
				Morale:    1,
			}

			err := tx.Create(personData).Error
			if err != nil {
				return err
			}

			termsData := &data.EmploymentTerms{
				Title:     "Pyrate",
				StartDate: uint(conn.Context().Timeline.CurrentTick()),
			}

			err = tx.Create(termsData).Error
			if err != nil {
				return err
			}

			contractData := &data.Contract{
				OfferorID:   crew.Persistent.ID,
				OfferorKind: data.CrewContractParty,

				OffereeKind: data.PersonContractParty,
				OffereeID:   personData.ID,

				TermsKind: data.EmploymentContractTerms,
				TermsID:   termsData.ID,
			}

			return tx.Create(contractData).Error
		})
		if err != nil {
			return nil, err
		}
	}

	if err := crew.Save(); err != nil {
		return nil, fmt.Errorf("failed to save crew: %w", err)
	}

	return &HireCrewResp{}, nil

}

package crew

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
	"gorm.io/gorm"
)

type Crew struct {
	Persistent data.Crew
}

func (s *Crew) Create() (uint, error) {
	err := db.Conn().Create(&s.Persistent).Error
	if err != nil {
		return 0, err
	}

	return s.Persistent.ID, nil
}

func (s *Crew) Save() error {
	err := db.Conn().Save(s.Persistent).Error

	if err != nil {
		return err
	}

	onChangedRegistryMap.Invoke([]any{s.Persistent.ID}, OnChangedEvent{
		Crew: *s,
	})

	return nil
}

func Get(id uint) (*Crew, error) {
	var s data.Crew
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Crew{
		Persistent: s,
	}, nil
}

func Where(crewData data.Crew) (*Crew, error) {
	var s data.Crew
	err := db.Conn().Where(&crewData).First(&s).Error
	if err != nil {
		return nil, err
	}

	return &Crew{
		Persistent: s,
	}, nil
}

func (s *Crew) Fetch() (*Crew, error) {
	err := db.Conn().First(&s.Persistent, s.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}

type Summary struct {
	Size          int
	AverageMorale float64
}

func (s *Crew) Summary() (Summary, error) {
	averageMorale, err := s.AverageMorale()
	if err != nil {
		return Summary{}, fmt.Errorf("failed to get average morale: %w", err)
	}

	size, err := s.Size()
	if err != nil {
		return Summary{}, fmt.Errorf("failed to get size: %w", err)
	}

	return Summary{
		Size:          int(size),
		AverageMorale: averageMorale,
	}, nil
}

type Member struct {
	FirstName   string                 // through Person
	LastName    string                 // through Person
	NickName    string                 // through Person
	Age         uint                   // through Person
	Morale      float64                // through Person
	OfferorKind data.ContractPartyKind // through Contract
	OfferorID   uint                   // through Contract
	OffereeKind data.ContractPartyKind // through Contract
	OffereeID   uint                   // through Contract
	TermsKind   data.ContractTermsKind // through Contract
	TermsID     uint                   // through Contract
	Title       string                 // through EmploymentTerms
	StartDate   uint                   // through EmploymentTerms
	EndDate     uint                   // through EmploymentTerms
}

func (s *Crew) Members() ([]*Member, error) {
	membersData := []*Member{}
	err := db.Conn().
		Model(&data.Person{}).
		Select("*").
		Scopes(members(s)).
		Joins("join employment_terms on employment_terms.id=contracts.terms_id and contracts.terms_kind=\"employment\"").
		Find(&membersData).Error

	slog.Info("members", ",m", membersData)

	return membersData, err
}

func (s *Crew) Size() (uint, error) {
	var size int64
	err := db.Conn().
		Scopes(members(s)).
		Table("people").
		Count(&size).Error
	return uint(size), err
}

func (s *Crew) AverageMorale() (float64, error) {
	var averageMorale float64
	row := db.Conn().
		Scopes(members(s)).
		Table("people").
		Select("avg(people.morale)").
		Row()
	err := row.Scan(&averageMorale)
	return averageMorale, err
}

func members(crew *Crew) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("join contracts on contracts.offeree_id=people.id").
			Where("contracts.offeror_id=?", crew.Persistent.ID).
			Where("contracts.offeror_kind=?", data.CrewContractParty)
	}
}

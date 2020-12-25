package database

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

type Bolt struct {
	db *storm.DB
	s  storm.Node
}

func NewBolt(pathToDatabase string) (*Bolt, error) {
	db, err := storm.Open(pathToDatabase)
	if err != nil {
		return nil, errors.Wrap(err, "open bolt database")
	}

	return &Bolt{
		db: db,
		s:  db.From("systems"),
	}, nil
}

func (b *Bolt) AddSystem(system structs.System) error {
	return b.s.Save(&system)
}

func (b *Bolt) DisableSystem(id string) error {
	return b.s.UpdateField(&structs.System{ID: id}, "IsEnabled", false)
}

func (b *Bolt) GetSystems() ([]structs.System, error) {
	var systems []structs.System
	err := b.s.All(&systems)
	return systems, err
}

func (b *Bolt) Close() error {
	return b.db.Close()
}

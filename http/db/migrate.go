package db

import (
	"github.com/go-gormigrate/gormigrate"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	MigrateID201811092300 = "201811092300"
	MigrateID201811092315 = "201811092315"
	MigrateIDCallback     = "0"
)

func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// you migrations here
		// create city table
		{
			ID: MigrateID201811092300,
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type City struct {
					gorm.Model
					Name string
				}
				return tx.AutoMigrate(&City{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				log.WithFields(log.Fields{"ID": MigrateID201811092315}).Warn("callback")
				return tx.DropTable("cities").Error
			},
		},

		// create order table
		{
			ID: MigrateID201811092315,
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type Order struct {
					gorm.Model
					Name   string
					CityID uint
				}

				if err := tx.AutoMigrate(&Order{}).Error; err != nil {
					return err
				}
				if err := tx.Model(Order{}).AddForeignKey("city_id", "cities (id)", "RESTRICT", "RESTRICT").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				log.WithFields(log.Fields{"ID": MigrateID201811092315}).Warn("callback")
				return tx.DropTable("orders").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Debugf("error occurred, rolling back to: %s", MigrateIDCallback)
		m.RollbackTo(MigrateIDCallback)
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Infof("Migration did run successfully")
}

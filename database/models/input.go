package models

type Input struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	RegisterIndex uint16
	Loads         []Load `gorm:"many2many:input_loads;"`
}

type Load struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	RegisterIndex uint16
}

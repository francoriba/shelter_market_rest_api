package database

import "gorm.io/gorm"

// GormDB wraps *gorm.DB to satisfy the Database interface
type GormDB struct {
	*gorm.DB
}

// NewGormDB initializes a new GormDB instance
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{DB: db}
}

// Create inserts a new record into the database
func (gdb *GormDB) Create(model interface{}) error {
	return gdb.DB.Create(model).Error
}

// Begin starts a transaction
func (gdb *GormDB) Begin() *gorm.DB {
	return gdb.DB.Begin()
}

// Commit commits the transaction
func (gdb *GormDB) Commit() *gorm.DB {
	return gdb.DB.Commit()
}

// Rollback rolls back the transaction
func (gdb *GormDB) Rollback() *gorm.DB {
	return gdb.DB.Rollback()
}

// Close closes the database connection
func (gdb *GormDB) Close() error {
	sqlDB, err := gdb.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

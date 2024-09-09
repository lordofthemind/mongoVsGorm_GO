package pkgs

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// CheckAndEnableUUIDExtension checks if the uuid-ossp extension is enabled in PostgreSQL, and enables it if not.
func CheckAndEnableUUIDExtension(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection from GORM: %v", err)
	}

	// Check if the uuid-ossp extension is enabled
	row := sqlDB.QueryRow("SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'")
	var exists int
	err = row.Scan(&exists)

	if err == sql.ErrNoRows {
		// Extension does not exist, create it
		_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if err != nil {
			return fmt.Errorf("failed to create uuid-ossp extension: %v", err)
		}
		log.Println("uuid-ossp extension enabled successfully")
	} else if err != nil {
		return fmt.Errorf("failed to check for uuid-ossp extension: %v", err)
	} else {
		log.Println("uuid-ossp extension is already enabled")
	}

	return nil
}

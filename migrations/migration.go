package migrations

import (
	"log"
	db "webserver/library/database"
)

var migration_sql = `
CREATE TABLE IF NOT EXISTS user (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(500) NOT NULL,
	user VARCHAR(450) NOT NULL,
	email VARCHAR(450) NOT NULL,
	hashcode LONGTEXT NOT NULL,
	dtupdate datetime,
	dtcreate datetime default CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE INDEX id_UNIQUE (id ASC) VISIBLE);
`

func RunMigrations() {
	log.Println("[MIGRATION] Run migrations...")
	db.RunSQL(migration_sql)
	log.Println("[MIGRATION] Migrations runned!")
}

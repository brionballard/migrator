package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Migrator struct {
	Migrations []MigrationSchema
}

type MigrationSchema struct {
	Name    string
	Columns []string
}

func (m *Migrator) Migration(n string, c []string) {
	e := MigrationSchema{
		Name:    n,
		Columns: c,
	}
	m.Migrations = append(m.Migrations, e)
}

func Build(m MigrationSchema) string {
	query := "CREATE TABLE IF NOT EXISTS " + m.Name + "(" + getIDcolumn() + strings.Join(m.Columns[:], ",") + "," + getDateColumns() + ")"

	return query
}

func getIDcolumn() string {
	return "id int primary key auto_increment,"
}

func getDateColumns() string {
	return "created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP"
}

func (m *Migrator) Migrate() {
	db, err := sql.Open("mysql", getCreds())

	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	for _, e := range m.Migrations {
		migration := Build(e)

		fmt.Println("Migrating " + e.Name + " table.")
		_, err := db.Exec(migration)

		if err != nil {
			fmt.Println("Failed to migrate " + e.Name + " table.")
			panic(err)
		}
	}

	defer db.Close()
}

func getCreds() string {
	godotenv.Load()

	return os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
}

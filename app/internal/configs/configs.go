package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetConnString(dbname string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbhostsip := "ywts-db"
	dbusername := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	return dbusername + ":" + dbpassword + "@tcp(" + dbhostsip + ")/" + dbname
}

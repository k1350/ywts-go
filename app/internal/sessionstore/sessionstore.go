package sessionstore

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/srinathgs/mysqlstore"

	"app/internal/configs"
)

type sessionstore struct {
	Store *mysqlstore.MySQLStore
}

var sharedInstance *sessionstore = newSessionstore()

func newSessionstore() *sessionstore {
	sqlConnString := configs.GetConnString("ywts")
	store, err := mysqlstore.NewMySQLStore(sqlConnString+"?parseTime=true&loc=Local", "sessions", "/", 3600, []byte(getSessionKey()))
	if err != nil {
		panic(err)
	}
	return &sessionstore{
		Store: store,
	}
}

func GetInstance() *sessionstore {
	return sharedInstance
}

func CloseInstance() {
	sharedInstance.Store.Close()
}

func getSessionKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv("SESSION_KEY")
}

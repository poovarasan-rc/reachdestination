package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var GDBCon *sql.DB

func OpenConn() error {
	var lErr error

	GDBCon, lErr = BuildConn()
	if lErr != nil {
		log.Println("Open connection failed:", lErr)
		return lErr
	}
	return nil

}

func CloseConn() {
	if GDBCon != nil {
		GDBCon.Close()
	}
}

func BuildConn() (*sql.DB, error) {

	dbconfig := ReadTomlConfig("./toml/dbconfig.toml")

	Server := fmt.Sprintf("%v", dbconfig.(map[string]interface{})["Server"])
	Port, _ := strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["Port"]))
	User := fmt.Sprintf("%v", dbconfig.(map[string]interface{})["User"])
	Password := fmt.Sprintf("%v", dbconfig.(map[string]interface{})["Password"])
	Database := fmt.Sprintf("%v", dbconfig.(map[string]interface{})["Database"])

	MaxOpenConns, _ := strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MaxOpenConns"]))
	MaxIdleConns, _ := strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MaxIdleConns"]))
	MaxIdleTime, _ := strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MaxIdleTime"]))

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", User, Password, Server, Port, Database)

	lDB, lErr := sql.Open("mysql", connString)
	if lErr != nil {
		log.Println("Open connection failed:", lErr)
		return nil, lErr
	}

	// Set the maximum number of open connections (max pool size)
	lDB.SetMaxOpenConns(MaxOpenConns) // Adjust this value as needed   // val-:100
	// Set the maximum number of idle connections in the pool	// val-:5
	lDB.SetMaxIdleConns(MaxIdleConns)

	lDB.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Second)

	return lDB, nil
}

func ReadTomlConfig(filename string) interface{} {
	var f interface{}
	if _, err := toml.DecodeFile(filename, &f); err != nil {
		log.Println(err)
	}
	return f
}

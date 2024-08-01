package utils

import (
	"backend/adapters"
	"backend/config"
	"fmt"
	"time"
)

func Capture(action string, data string, username string) (string, error) {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return "", fmt.Errorf("[createlog error] %s", err)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return "", fmt.Errorf("[createlog error] %s", err)
	}

	datenow := time.Now()
	dateformatted := datenow.Format(config.App.DateFormat_Timestamp)

	query := fmt.Sprintf(`INSERT INTO log (data, query, date, user) VALUES ('%s', '%s', '%s', '%s');`, data, action, dateformatted, username)

	create, err := sql.Connection.Query(query)
	if err != nil {
		return "", fmt.Errorf("[createlog error] %s", err)
	}
	defer create.Close()

	return "", nil
}

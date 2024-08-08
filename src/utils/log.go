package utils

import (
	"backend/adapters"
	"backend/config"
	"fmt"
	"time"
)

func Capture(action string, data string, userid string) (string, error) {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return "", fmt.Errorf("[createlog error1] %s", err)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return "", fmt.Errorf("[createlog error2] %s", err)
	}

	// Get user name
	var id, name string
	query_user := fmt.Sprintf(`SELECT id, name FROM user WHERE id = '%s' LIMIT 1`, userid)
	if err = sql.Connection.QueryRow(query_user).Scan(&id, &name); err != nil {
		if err.Error() == `sql: no rows in result set` {
			name = userid
		} else {
			return "", fmt.Errorf("[createlog error3] %s", err)
		}
	}

	datenow := time.Now()
	dateformatted := datenow.Format(config.App.DateFormat_Timestamp)

	query := fmt.Sprintf(`INSERT INTO log (data, query, date, user) VALUES ('%s', '%s', '%s', '%s');`, data, action, dateformatted, name)

	create, err := sql.Connection.Query(query)
	if err != nil {
		return "", fmt.Errorf("[createlog error] %s", err)
	}
	defer create.Close()

	return "", nil
}

package user

import (
	"backend/adapters"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data            []User `json:"data"`
	RecordsTotal    string `json:"recordsTotal,omitempty"`
	RecordsFiltered string `json:"recordsFiltered,omitempty"`
}

type User struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Role    string `json:"role,omitempty"`
	Status  string `json:"status,omitempty"`
	Account string `json:"account,omitempty"`
	Picture int    `json:"picture,omitempty"`
	Hidden  int    `json:"hidden,omitempty"`
}

func Get(ctx *gin.Context) (Response, error) {
	var totalrows int
	var query, search, query_datatables string
	idParam := ctx.DefaultQuery("id", "0")

	// direct handling
	LimitParam := ctx.DefaultQuery("limit", "10")
	OffsetParam := ctx.DefaultQuery("offset", "0")

	// datatables handling
	SearchValue := ctx.DefaultQuery("search[value]", "")
	LengthParam := ctx.DefaultQuery("length", "")
	StartParam := ctx.DefaultQuery("start", "")
	if LengthParam != "" && StartParam != "" {
		LimitParam = LengthParam
		OffsetParam = StartParam
	}

	limit, err := strconv.Atoi(LimitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(OffsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return Response{}, err
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		// datatables total rows and filtered handling
		if SearchValue != "" {
			search = fmt.Sprintf(`WHERE (name LIKE '%%%s%%' OR email LIKE '%%%s%%' OR role LIKE '%%%s%%' OR status LIKE '%%%s%%' OR account LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue)
		}

		query = fmt.Sprintf(`SELECT id, name, email, role, status, account, hidden FROM user %s ORDER BY id DESC LIMIT %d OFFSET %d`, search, limit, offset)

		query_datatables = fmt.Sprintf(`SELECT COUNT(id) as totalrows FROM user %s ORDER BY id DESC`, search)
		if err = sql.Connection.QueryRow(query_datatables).Scan(&totalrows); err != nil {
			if err.Error() == `sql: no rows in result set` {
				totalrows = 0
			} else {
				return Response{}, err
			}
		}

	} else {
		query = fmt.Sprintf(`SELECT id, name, email, role, status, account, hidden FROM user WHERE id = %d LIMIT 1`, id)
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var id, hidden int
		var name, role, status, account, email string

		if err := rows.Scan(&id, &name, &email, &role, &status, &account, &hidden); err != nil {
			return Response{}, err
		}

		if role == "1" {
			role = "Root"
		} else if role == "2" {
			role = "Administrator"
		} else if role == "3" {
			role = "Sales Order"
		} else if role == "4" {
			role = "Finance"
		} else if role == "6" {
			role = "Production"
		} else {
			role = "Guest"
		}

		if status == "1" {
			status = "Verified"
		} else {
			status = "Not Verified"
		}

		if account == "1" {
			account = "Active"
		} else {
			account = "Inactive"
		}

		users = append(users, User{
			Id:      id,
			Name:    name,
			Email:   email,
			Role:    role,
			Status:  status,
			Account: account,
			Hidden:  hidden,
		})
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}

	response.Data = users

	return response, nil
}

func Create(Name string, Email string, Password string, Role int, Status int, Account int) ([]User, error) {
	if Password == "" {
		return nil, errors.New("password empty")
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id FROM user WHERE email = '%s' LIMIT 1", Email)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Unique email validation
	if rows.Next() {
		return nil, errors.New("email registered")
	}

	queryCreate := fmt.Sprintf("INSERT INTO user (name, email, password, role, status, account, picture, hidden) VALUES ('%s', '%s', md5('%s'), '%d', '%d', '%d', '', 0)", Name, Email, Password, Role, Status, Account)
	create, err := sql.Connection.Query(queryCreate)
	if err != nil {
		return nil, err
	}

	defer create.Close()

	return []User{}, nil
}

func Update(Id int, Name string, Email string, Role int, Status int, Account int, Hidden int) ([]User, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM user WHERE id = '%d' LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	// User ID validation
	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	// Unique email validation
	var id int
	query_email := fmt.Sprintf(`SELECT id FROM user WHERE email = '%s' LIMIT 1`, Email)
	if err = sql.Connection.QueryRow(query_email).Scan(&id); err != nil {
		if err.Error() == `sql: no rows in result set` {
		} else {
			return nil, err
		}
	}

	if Id != id {
		return nil, errors.New("email registered")
	}

	queryUpdate := fmt.Sprintf("UPDATE user SET name ='%s', email ='%s', role ='%d', status ='%d', account ='%d', hidden ='%d' WHERE id = %d", Name, Email, Role, Status, Account, Hidden, Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	return []User{}, nil
}

func Delete(Id int) ([]User, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM user WHERE id = %d LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	query := fmt.Sprintf("DELETE FROM user WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return []User{}, nil
}

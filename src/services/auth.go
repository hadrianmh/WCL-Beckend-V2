package services

import (
	"backend/adapters"
	"backend/utils"
	"errors"
	"fmt"
	"reflect"
)

type User struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Role    string `json:"role,omitempty"`
	Status  string `json:"status,omitempty"`
	Account string `json:"account,omitempty"`
}

type ResponseRequset struct {
	// Data    []User `json:"data,omitempty"`
	Data    utils.Tokens `json:"data,omitempty"`
	Message string       `json:"message,omitempty"`
}

// function list stored
var FunctionMap = map[string]interface{}{
	"Login":   Login,
	"Refresh": Refresh,
}

// CheckFunction dynamically calls a function
func Auth(action string, email string, password string, refresh_token string) (*ResponseRequset, error) {

	act := utils.Ucfirst(utils.StrReplaceAll(action, "/", ""))

	if check, err := CheckFunction(act, email, password, refresh_token); err != nil {
		return nil, err
	} else {
		return check, nil
	}
}

func CheckFunction(funcName string, em string, pw string, reftkn string) (*ResponseRequset, error) {
	if fn, exists := FunctionMap[funcName]; exists {

		if reflect.ValueOf(fn).Kind() == reflect.Func {

			if f, ok := fn.(func(string, string, string) *ResponseRequset); ok {
				return f(em, pw, reftkn), nil
			}

			return nil, errors.New("bad Request")
		}

		return nil, errors.New("bad Request")
	}

	return nil, errors.New("bad request")
}

func Login(email string, password string, refresh_token string) *ResponseRequset {
	sql, err := adapters.NewSql()
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	query := fmt.Sprintf("SELECT id, name, email, role, status, account FROM user WHERE email='%s' AND password=md5('%s') LIMIT 1", email, password)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var id int
		var name, email, role, status, account string

		if err := rows.Scan(&id, &name, &email, &role, &status, &account); err != nil {
			jsonResp := &ResponseRequset{
				Message: err.Error(),
				Data:    utils.Tokens{},
			}
			return jsonResp
		}

		users = append(users, User{
			Id:      id,
			Name:    name,
			Email:   email,
			Role:    role,
			Status:  status,
			Account: account,
		})
	}

	if err := rows.Err(); err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Fail // User not found
	if len(users) < 1 {
		jsonResp := &ResponseRequset{
			Message: "Wrong email or password",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Fail // Email user not verified
	if users[0].Status != "1" {
		jsonResp := &ResponseRequset{
			Message: "Please verify email to continue or contact the administrator",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Fail // User suspended
	if users[0].Account != "1" {
		jsonResp := &ResponseRequset{
			Message: "Account suspended. contact the administrator",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	tokens, err := utils.GenerateTokens(users[0].Id, email, users[0].Name, users[0].Role)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Success
	jsonResp := &ResponseRequset{
		Message: "",
		Data:    tokens,
	}
	return jsonResp
}

func Refresh(email string, password string, refresh_token string) *ResponseRequset {
	if len(refresh_token) < 1 {
		jsonResp := &ResponseRequset{
			Message: "bad request",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	username, err := utils.ValidateToken(refresh_token, "refresh_token")
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	sql, err := adapters.NewSql()
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	query := fmt.Sprintf("SELECT id, name, email, role, status, account FROM user WHERE email='%s' LIMIT 1", username)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var id int
		var name, email, role, status, account string

		if err := rows.Scan(&id, &name, &email, &role, &status, &account); err != nil {
			jsonResp := &ResponseRequset{
				Message: err.Error(),
				Data:    utils.Tokens{},
			}
			return jsonResp
		}

		users = append(users, User{
			Id:      id,
			Name:    name,
			Email:   email,
			Role:    role,
			Status:  status,
			Account: account,
		})
	}

	if err := rows.Err(); err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	tokens, err := utils.GenerateTokens(users[0].Id, username, users[0].Name, users[0].Role)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Success
	jsonResp := &ResponseRequset{
		Message: "",
		Data:    tokens,
	}
	return jsonResp
}

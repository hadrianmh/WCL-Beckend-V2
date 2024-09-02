package auth

import (
	"backend/adapters"
	"backend/utils"
	"errors"
	"fmt"
	"reflect"
)

type UserData struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Status  string `json:"status"`
	Account string `json:"account"`
	Picture string `json:"picture"`
	Hidden  string `json:"hidden"`
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
func Init(action string, email string, password string, refresh_token string) (*ResponseRequset, error) {

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

	return nil, errors.New("bad Request")
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

	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT id, name, email, role, status, account, picture FROM user WHERE email='%s' AND password=md5('%s') LIMIT 1", email, password)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	defer rows.Close()

	usersdata := []UserData{}
	for rows.Next() {
		var id, name, email, role, status, account, picture string

		if err := rows.Scan(&id, &name, &email, &role, &status, &account, &picture); err != nil {
			jsonResp := &ResponseRequset{
				Message: err.Error(),
				Data:    utils.Tokens{},
			}
			return jsonResp
		}

		usersdata = append(usersdata, UserData{
			Id:      id,
			Name:    name,
			Email:   email,
			Role:    role,
			Status:  status,
			Account: account,
			Picture: picture,
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
	if len(usersdata) < 1 {
		jsonResp := &ResponseRequset{
			Message: "Wrong email or password",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Fail // Email user not verified
	if usersdata[0].Status != "1" {
		jsonResp := &ResponseRequset{
			Message: "Please verify email to continue or contact the administrator",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	// Fail // User suspended
	if usersdata[0].Account != "1" {
		jsonResp := &ResponseRequset{
			Message: "Account suspended. contact the administrator",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	tokens, err := utils.GenerateTokens(usersdata[0].Id, email, usersdata[0].Name, usersdata[0].Role, usersdata[0].Status, usersdata[0].Account, usersdata[0].Picture)
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
			Message: "Bad Request",
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	userid, err := utils.ValidateToken(refresh_token, "refresh_token")
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

	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT id, name, email, role, status, account, picture FROM user WHERE id='%s' LIMIT 1", userid)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	defer rows.Close()

	fmt.Println(userid)
	fmt.Println(query)

	usersdata := []UserData{}
	for rows.Next() {
		var id, name, email, role, status, account, picture string

		if err := rows.Scan(&id, &name, &email, &role, &status, &account, &picture); err != nil {
			jsonResp := &ResponseRequset{
				Message: err.Error(),
				Data:    utils.Tokens{},
			}
			return jsonResp
		}

		usersdata = append(usersdata, UserData{
			Id:      id,
			Name:    name,
			Email:   email,
			Role:    role,
			Status:  status,
			Account: account,
			Picture: picture,
		})
	}

	if err := rows.Err(); err != nil {
		jsonResp := &ResponseRequset{
			Message: err.Error(),
			Data:    utils.Tokens{},
		}
		return jsonResp
	}

	tokens, err := utils.GenerateTokens(userid, usersdata[0].Email, usersdata[0].Name, usersdata[0].Role, usersdata[0].Status, usersdata[0].Account, usersdata[0].Picture)
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

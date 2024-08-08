package company

import (
	"backend/adapters"
	"backend/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data            []Company `json:"data"`
	RecordsTotal    string    `json:"recordsTotal,omitempty"`
	RecordsFiltered string    `json:"recordsFiltered,omitempty"`
}

type Company struct {
	Id          int    `json:"id,omitempty"`
	CompanyName string `json:"companyname,omitempty"`
	Address     string `json:"address,omitempty"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
	Logo        string `json:"logo"`
	InputBy     int    `json:"inputby,omitempty"`
	Hidden      int    `json:"hidden,omitempty"`
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
	if err != nil || limit < -1 {
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
			search = fmt.Sprintf(`WHERE (company LIKE '%%%s%%' OR address LIKE '%%%s%%' OR email LIKE '%%%s%%' OR phone LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue)
		}

		query_datatables = fmt.Sprintf(`SELECT COUNT(id) as totalrows FROM company %s ORDER BY id DESC`, search)
		if err = sql.Connection.QueryRow(query_datatables).Scan(&totalrows); err != nil {
			if err.Error() == `sql: no rows in result set` {
				totalrows = 0
			} else {
				return Response{}, err
			}
		}

		// If request limit -1 (pagination datatables) is show all
		if limit == -1 {
			limit = totalrows
		}

		query = fmt.Sprintf(`SELECT id, company, address, email, phone, logo, input_by, hidden FROM company %s ORDER BY id DESC LIMIT %d OFFSET %d`, search, limit, offset)

	} else {
		query = fmt.Sprintf(`SELECT id, company, address, email, phone, logo, input_by, hidden FROM company WHERE id = %d LIMIT 1`, id)
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	companys := []Company{}
	for rows.Next() {
		var id, inputby, hidden int
		var companyname, address, email, logo, phone string

		if err := rows.Scan(&id, &companyname, &address, &email, &phone, &logo, &inputby, &hidden); err != nil {
			return Response{}, err
		}

		companys = append(companys, Company{
			Id:          id,
			CompanyName: companyname,
			Address:     address,
			Email:       email,
			Phone:       phone,
			Logo:        logo,
			InputBy:     inputby,
			Hidden:      hidden,
		})
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}

	response.Data = companys

	return response, nil
}

func Create(Sessionid string, CompanyName string, Address string, Email string, Phone string, Logo string, InputBy int) ([]Company, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id FROM company WHERE company = '%s' LIMIT 1", CompanyName)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Unique email validation
	if rows.Next() {
		return nil, errors.New("company registered")
	}

	// Validate image format
	allowedFormats := map[string]bool{"jpg": true, "jpeg": true, "png": true}
	getFormat := utils.GetImageFormat(Logo)

	if getFormat != "" {
		format := strings.ToLower(getFormat)
		if !allowedFormats[format] {
			return nil, errors.New("only JPG, JPEG, and PNG formats are allowed")
		}
	}

	queryCreate := fmt.Sprintf("INSERT INTO company (company, address, email, phone, logo, input_by, hidden) VALUES ('%s', '%s', '%s', '%s', '%s', '%d','0')", CompanyName, Address, Email, Phone, Logo, InputBy)
	create, err := sql.Connection.Query(queryCreate)
	if err != nil {
		return nil, err
	}

	defer create.Close()

	// Log capture
	utils.Capture(
		`Company Created`,
		fmt.Sprintf(`Company: %s - Address: %s - Phone: %s - Email: %s`, CompanyName, Address, Phone, Email),
		Sessionid,
	)

	return []Company{}, nil
}

func Update(Sessionid string, Id int, CompanyName string, Address string, Email string, Phone string, Logo string, InputBy int, Hidden int) ([]Company, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM company WHERE id = '%d' LIMIT 1", Id)
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
	query_id = fmt.Sprintf(`SELECT id FROM company WHERE company = '%s' LIMIT 1`, CompanyName)
	if err = sql.Connection.QueryRow(query_id).Scan(&id); err != nil {
		if err.Error() == `sql: no rows in result set` {
			id = Id
		} else {
			return nil, err
		}
	}

	if Id != id {
		return nil, errors.New("company registered")
	}

	queryUpdate := fmt.Sprintf("UPDATE company SET company ='%s', address ='%s', email ='%s', phone ='%s', logo ='%s', input_by ='%d', hidden ='%d' WHERE id = %d", CompanyName, Address, Email, Phone, Logo, InputBy, Hidden, Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	// Log capture
	utils.Capture(
		`Company Edited`,
		fmt.Sprintf(`Companyid: %d - Company: %s - Address: %s - Phone: %s - Email: %s`, Id, CompanyName, Address, Phone, Email),
		Sessionid,
	)

	return []Company{}, nil
}

func Delete(Sessionid string, Id int) ([]Company, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	var company, address, email, phone string
	query_id := fmt.Sprintf(`SELECT company, address, email, phone FROM company WHERE id = %d LIMIT 1`, Id)
	if err = sql.Connection.QueryRow(query_id).Scan(&company, &address, &email, &phone); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	query := fmt.Sprintf("DELETE FROM company WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Log capture
	utils.Capture(
		`Company Deleted`,
		fmt.Sprintf(`Companyid: %d - Company: %s - Address: %s - Phone: %s - Email: %s`, Id, company, address, phone, email),
		Sessionid,
	)

	return []Company{}, nil
}

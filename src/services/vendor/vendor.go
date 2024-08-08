package vendor

import (
	"backend/adapters"
	"backend/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data            []Vendor `json:"data"`
	RecordsTotal    string   `json:"recordsTotal,omitempty"`
	RecordsFiltered string   `json:"recordsFiltered,omitempty"`
}

type Vendor struct {
	Id         int    `json:"id,omitempty"`
	VendorName string `json:"vendorname,omitempty"`
	Address    string `json:"address,omitempty"`
	Phone      string `json:"phone"`
	InputBy    int    `json:"inputby,omitempty"`
	Hidden     int    `json:"hidden,omitempty"`
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
			search = fmt.Sprintf(`WHERE (vendor LIKE '%%%s%%' OR address LIKE '%%%s%%' OR phone LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue)
		}

		query_datatables = fmt.Sprintf(`SELECT COUNT(id) as totalrows FROM vendor %s ORDER BY id DESC`, search)
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

		query = fmt.Sprintf(`SELECT id, vendor, address, CASE WHEN phone = '' THEN '' ELSE phone END AS phone, input_by, hidden FROM vendor %s ORDER BY id DESC LIMIT %d OFFSET %d`, search, limit, offset)

	} else {
		query = fmt.Sprintf(`SELECT id, vendor, address, phone, input_by, hidden FROM vendor WHERE id = %d LIMIT 1`, id)
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	vendors := []Vendor{}
	for rows.Next() {
		var id, inputby, hidden int
		var vendorname, address, phone string

		if err := rows.Scan(&id, &vendorname, &address, &phone, &inputby, &hidden); err != nil {
			return Response{}, err
		}

		vendors = append(vendors, Vendor{
			Id:         id,
			VendorName: vendorname,
			Address:    address,
			Phone:      phone,
			InputBy:    inputby,
			Hidden:     hidden,
		})
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}

	response.Data = vendors

	return response, nil
}

func Create(Sessionid string, VendorName string, Address string, Phone int, InputBy int) ([]Vendor, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id FROM vendor WHERE vendor = '%s' LIMIT 1", VendorName)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Unique email validation
	if rows.Next() {
		return nil, errors.New("vendor registered")
	}

	queryCreate := fmt.Sprintf("INSERT INTO vendor (vendor, address, phone, input_by, hidden) VALUES ('%s', '%s', '%d', '%d', '0')", VendorName, Address, Phone, InputBy)
	create, err := sql.Connection.Query(queryCreate)
	if err != nil {
		return nil, err
	}

	defer create.Close()

	// Log capture
	utils.Capture(
		`Vendor Created`,
		fmt.Sprintf(`Vendor: %s - Address: %s - Phone: %d`, VendorName, Address, Phone),
		Sessionid,
	)

	return []Vendor{}, nil
}

func Update(Sessionid string, Id int, VendorName string, Address string, Phone int, InputBy int, Hidden int) ([]Vendor, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM vendor WHERE id = '%d' LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	// User ID validation
	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	query_email := fmt.Sprintf("SELECT id FROM vendor WHERE vendor = '%s' LIMIT 1", VendorName)
	rows_email, err := sql.Connection.Query(query_email)
	if err != nil {
		return nil, err
	}

	defer rows_email.Close()

	// Unique email validation
	if rows_email.Next() {
		return nil, errors.New("vendor registered")
	}

	queryUpdate := fmt.Sprintf("UPDATE vendor SET vendor ='%s', address ='%s', phone ='%d', input_by ='%d', hidden ='%d' WHERE id ='%d'", VendorName, Address, Phone, InputBy, Hidden, Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	// Log capture
	utils.Capture(
		`Vendor Updated`,
		fmt.Sprintf(`Vendorid: %d - Vendor: %s - Address: %s - Phone: %d`, Id, VendorName, Address, Phone),
		Sessionid,
	)

	return []Vendor{}, nil
}

func Delete(Sessionid string, Id int) ([]Vendor, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}
	var vendor, address, phone string
	query_id := fmt.Sprintf(`SELECT vendor, address, phone FROM vendor WHERE id = %d LIMIT 1`, Id)
	if err = sql.Connection.QueryRow(query_id).Scan(&vendor, &address, &phone); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	query := fmt.Sprintf("DELETE FROM vendor WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Log capture
	utils.Capture(
		`Vendor Deleted`,
		fmt.Sprintf(`Vendor: %s - Address: %s - Phone: %s`, vendor, address, phone),
		Sessionid,
	)

	return []Vendor{}, nil
}

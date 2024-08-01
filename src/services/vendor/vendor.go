package vendor

import (
	"backend/adapters"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Vendor struct {
	Id         int    `json:"id,omitempty"`
	VendorName string `json:"vendorname,omitempty"`
	Address    string `json:"address,omitempty"`
	Phone      int    `json:"phone,omitempty"`
	InputBy    int    `json:"inputby,omitempty"`
	Hidden     int    `json:"hidden,omitempty"`
}

func Get(ctx *gin.Context) ([]Vendor, error) {
	var query string
	idParam := ctx.DefaultQuery("id", "0")
	LimitParam := ctx.DefaultQuery("limit", "10")
	OffsetParam := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(LimitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(OffsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		query = fmt.Sprintf(`SELECT id, vendor, address, phone, input_by, hidden FROM vendor ORDER BY id DESC LIMIT %d OFFSET %d`, limit, offset)
	} else {
		query = fmt.Sprintf(`SELECT id, vendor, address, phone, input_by, hidden FROM vendor WHERE id = %d LIMIT 1`, id)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	vendors := []Vendor{}
	for rows.Next() {
		var id, inputby, hidden, phone int
		var vendorname, address string

		if err := rows.Scan(&id, &vendorname, &address, &phone, &inputby, &hidden); err != nil {
			return nil, err
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
		return nil, err
	}

	return vendors, nil
}

func Create(VendorName string, Address string, Phone int, InputBy int) ([]Vendor, error) {
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

	return []Vendor{}, nil
}

func Update(Id int, VendorName string, Address string, Phone int, InputBy int, Hidden int) ([]Vendor, error) {
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

	return []Vendor{}, nil
}

func Delete(Id int) ([]Vendor, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM vendor WHERE id = %d LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	query := fmt.Sprintf("DELETE FROM vendor WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return []Vendor{}, nil
}

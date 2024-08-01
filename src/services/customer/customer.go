package customer

import (
	"backend/adapters"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	Id           int    `json:"id,omitempty"`
	CustomerName string `json:"customername,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postalcode,omitempty"`
	Phone        string `json:"phone,omitempty"`
	SName        string `json:"sname,omitempty"`
	SAddress     string `json:"saddress,omitempty"`
	SCity        string `json:"scity,omitempty"`
	SCountry     string `json:"scountry,omitempty"`
	SProvince    string `json:"sprovince,omitempty"`
	SPostalCode  string `json:"spostalcode,omitempty"`
	SPhone       string `json:"sphone,omitempty"`
	InputBy      int    `json:"inputby,omitempty"`
	Hidden       int    `json:"hidden,omitempty"`
}

func Get(ctx *gin.Context) ([]Customer, error) {
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
		query = fmt.Sprintf(`SELECT id, nama, alamat, kota, negara, provinsi, kodepos, telp, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, s_telp, input_by, hidden FROM customer ORDER BY id DESC LIMIT %d OFFSET %d`, limit, offset)
	} else {
		query = fmt.Sprintf(`SELECT id, nama, alamat, kota, negara, provinsi, kodepos, telp, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, s_telp, input_by, hidden FROM customer WHERE id = %d LIMIT 1`, id)
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

	customers := []Customer{}
	for rows.Next() {
		var id, inputby, hidden int
		var customername, address, city, postalcode, country, province, sname, saddress, scity, scountry, sprovince, spostalcode, phone, sphone string

		if err := rows.Scan(&id, &customername, &address, &city, &country, &province, &postalcode, &phone, &sname, &saddress, &scity, &scountry, &sprovince, &spostalcode, &sphone, &inputby, &hidden); err != nil {
			return nil, err
		}

		customers = append(customers, Customer{
			Id:           id,
			CustomerName: customername,
			Address:      address,
			City:         city,
			Country:      country,
			Province:     province,
			PostalCode:   string(postalcode),
			Phone:        string(phone),
			SName:        sname,
			SAddress:     saddress,
			SCity:        scity,
			SCountry:     scountry,
			SProvince:    sprovince,
			SPostalCode:  string(spostalcode),
			SPhone:       string(sphone),
			InputBy:      inputby,
			Hidden:       hidden,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func Create(CustomerName string, Address string, City string, Country string, Province string, PostalCode string, Phone string, SName string, SAddress string, SCity string, SCountry string, SProvince string, SPostalCode string, SPhone string, InputBy int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id FROM customer WHERE nama = '%s' LIMIT 1", CustomerName)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Unique email validation
	if rows.Next() {
		return nil, errors.New("customer registered")
	}

	queryCreate := fmt.Sprintf("INSERT INTO customer (nama, alamat, kota, negara, provinsi, kodepos, telp, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, s_telp, input_by, hidden) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '0')", CustomerName, Address, City, Country, Province, PostalCode, Phone, SName, SAddress, SCity, SCountry, SProvince, SPostalCode, SPhone, InputBy)
	create, err := sql.Connection.Query(queryCreate)
	if err != nil {
		return nil, err
	}

	defer create.Close()

	return []Customer{}, nil
}

func Update(Id int, CustomerName string, Address string, City string, Country string, Province string, PostalCode string, Phone string, SName string, SAddress string, SCity string, SCountry string, SProvince string, SPostalCode string, SPhone string, InputBy int, Hidden int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM customer WHERE id = '%d' LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	// User ID validation
	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	query_email := fmt.Sprintf("SELECT id FROM customer WHERE nama = '%s' LIMIT 1", CustomerName)
	rows_email, err := sql.Connection.Query(query_email)
	if err != nil {
		return nil, err
	}

	defer rows_email.Close()

	// Unique email validation
	if rows_email.Next() {
		return nil, errors.New("customer registered")
	}

	queryUpdate := fmt.Sprintf("UPDATE customer SET nama ='%s', alamat ='%s', kota ='%s', negara ='%s', provinsi ='%s', kodepos ='%s', telp ='%s', s_nama ='%s', s_alamat ='%s', s_kota ='%s', s_negara ='%s', s_provinsi ='%s', s_kodepos ='%s', s_telp ='%s', input_by ='%d', hidden ='%d' WHERE id ='%d'", CustomerName, Address, City, Country, Province, PostalCode, Phone, SName, SAddress, SCity, SCountry, SProvince, SPostalCode, SPhone, InputBy, Hidden, Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	return []Customer{}, nil
}

func Delete(Id int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM customer WHERE id = %d LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	defer rows_id.Close()

	if !rows_id.Next() {
		return nil, errors.New("invalid ID")
	}

	query := fmt.Sprintf("DELETE FROM customer WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return []Customer{}, nil
}

package customer

import (
	"backend/adapters"
	"backend/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data            []Customer `json:"data"`
	RecordsTotal    string     `json:"recordsTotal,omitempty"`
	RecordsFiltered string     `json:"recordsFiltered,omitempty"`
}
type Customer struct {
	Id           int    `json:"id,omitempty"`
	CustomerName string `json:"customername,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postalcode,omitempty"`
	Phone        string `json:"phone,omitempty"`
	SName        string `json:"sname"`
	SAddress     string `json:"saddress,omitempty"`
	SCity        string `json:"scity,omitempty"`
	SCountry     string `json:"scountry,omitempty"`
	SProvince    string `json:"sprovince,omitempty"`
	SPostalCode  string `json:"spostalcode,omitempty"`
	SPhone       string `json:"sphone,omitempty"`
	B_alamat     string `json:"b_alamat"`
	S_alamat     string `json:"s_alamat"`
	InputBy      int    `json:"inputby,omitempty"`
	Hidden       int    `json:"hidden,omitempty"`
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

	defer sql.Connection.Close()

	Id, err := strconv.Atoi(idParam)
	if err != nil || Id < 1 {
		// datatables total rows and filtered handling
		if SearchValue != "" {
			search = fmt.Sprintf(`WHERE (nama LIKE '%%%s%%' OR alamat LIKE '%%%s%%' OR kota LIKE '%%%s%%' OR negara LIKE '%%%s%%' OR provinsi LIKE '%%%s%%' OR kodepos LIKE '%%%s%%' OR telp LIKE '%%%s%%' OR s_nama LIKE '%%%s%%' OR s_alamat LIKE '%%%s%%' OR s_kota LIKE '%%%s%%' OR s_negara LIKE '%%%s%%' OR s_provinsi LIKE '%%%s%%' OR s_kodepos LIKE '%%%s%%' OR s_telp LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue)
		}

		query_datatables = fmt.Sprintf(`SELECT COUNT(id) as totalrows FROM customer %s ORDER BY id DESC`, search)
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

		query = fmt.Sprintf(`SELECT id, nama, alamat, kota, negara, provinsi, kodepos, telp, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, s_telp, input_by, hidden FROM customer %s ORDER BY id DESC LIMIT %d OFFSET %d`, search, limit, offset)

	} else {
		query = fmt.Sprintf(`SELECT id, nama, alamat, kota, negara, provinsi, kodepos, telp, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, s_telp, input_by, hidden FROM customer WHERE id = %d LIMIT 1`, Id)
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var id, inputby, hidden int
		var customername, address, city, postalcode, country, province, sname, saddress, scity, scountry, sprovince, spostalcode, phone, sphone string

		if err := rows.Scan(&id, &customername, &address, &city, &country, &province, &postalcode, &phone, &sname, &saddress, &scity, &scountry, &sprovince, &spostalcode, &sphone, &inputby, &hidden); err != nil {
			return Response{}, err
		}

		// Validation address
		if Id < 1 {
			if sname != "" {
				saddress = fmt.Sprintf(`%s. `, saddress)
				scity = fmt.Sprintf(`%s - `, scity)
				sprovince = fmt.Sprintf(`%s, `, sprovince)
				scountry = fmt.Sprintf(`%s. `, scountry)
				spostalcode = fmt.Sprintf(`%s. `, spostalcode)
			}

			if address != "" {
				address = fmt.Sprintf(`%s. `, address)
			}

			if city != "" {
				city = fmt.Sprintf(`%s - `, city)
			}

			if province != "" {
				province = fmt.Sprintf(`%s, `, province)
			}

			if country != "" {
				country = fmt.Sprintf(`%s. `, country)
			}

			if postalcode != "" {
				postalcode = fmt.Sprintf(`%s. `, postalcode)
			}
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
			B_alamat:     fmt.Sprintf(`%s%s%s%s%s`, address, city, province, country, postalcode),
			S_alamat:     fmt.Sprintf(`%s%s%s%s%s`, saddress, scity, sprovince, scountry, spostalcode),
			InputBy:      inputby,
			Hidden:       hidden,
		})
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}

	response.Data = customers

	return response, nil
}

func Create(Sessionid string, CustomerName string, Address string, City string, Country string, Province string, PostalCode string, Phone string, SName string, SAddress string, SCity string, SCountry string, SProvince string, SPostalCode string, SPhone string, InputBy int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	defer sql.Connection.Close()

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

	// Log capture
	utils.Capture(
		`Customer Created`,
		fmt.Sprintf(`Customer: %s - Address: %s - Phone: %s`, CustomerName, Address, Phone),
		Sessionid,
	)

	return []Customer{}, nil
}

func Update(Sessionid string, Id int, CustomerName string, Address string, City string, Country string, Province string, PostalCode string, Phone string, SName string, SAddress string, SCity string, SCountry string, SProvince string, SPostalCode string, SPhone string, InputBy int, Hidden int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	defer sql.Connection.Close()

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

	// Unique email validation
	var id int
	query_id = fmt.Sprintf(`SELECT id FROM customer WHERE nama = '%s' LIMIT 1`, CustomerName)
	if err = sql.Connection.QueryRow(query_id).Scan(&id); err != nil {
		if err.Error() == `sql: no rows in result set` {
			id = Id
		} else {
			return nil, err
		}
	}

	if Id != id {
		return nil, errors.New("customer registered")
	}

	queryUpdate := fmt.Sprintf("UPDATE customer SET nama ='%s', alamat ='%s', kota ='%s', negara ='%s', provinsi ='%s', kodepos ='%s', telp ='%s', s_nama ='%s', s_alamat ='%s', s_kota ='%s', s_negara ='%s', s_provinsi ='%s', s_kodepos ='%s', s_telp ='%s', input_by ='%d', hidden ='%d' WHERE id ='%d'", CustomerName, Address, City, Country, Province, PostalCode, Phone, SName, SAddress, SCity, SCountry, SProvince, SPostalCode, SPhone, InputBy, Hidden, Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	// Log capture
	utils.Capture(
		`Customer Updated`,
		fmt.Sprintf(`Customerid: %d - Customer: %s - Address: %s - Phone: %s`, Id, CustomerName, Address, Phone),
		Sessionid,
	)

	return []Customer{}, nil
}

func Delete(Sessionid string, Id int) ([]Customer, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	defer sql.Connection.Close()

	var customer, address, phone string
	query_id := fmt.Sprintf(`SELECT nama, alamat, telp FROM customer WHERE id = %d LIMIT 1`, Id)
	if err = sql.Connection.QueryRow(query_id).Scan(&customer, &address, &phone); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	query := fmt.Sprintf("DELETE FROM customer WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Log capture
	utils.Capture(
		`Customer Deleted`,
		fmt.Sprintf(`Customerid: %d - Customer: %s - Address: %s - Phone: %s`, Id, customer, address, phone),
		Sessionid,
	)

	return []Customer{}, nil
}

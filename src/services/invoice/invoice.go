package invoice

import (
	"backend/adapters"
	"backend/config"
	"backend/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type InvoiceItems struct {
	NoSj     string  `json:"no_sj"`
	NoSo     string  `json:"no_so"`
	Item     string  `json:"item"`
	Unit     string  `json:"unit"`
	SendQty  float64 `json:"send_qty"`
	Price    float64 `json:"price"`
	Ppn      float64 `json:"ppn,omitempty"`
	Cost     float64 `json:"cost,omitempty"`
	SubTotal float64 `json:"subtotal,omitempty"`
	Total    float64 `json:"total,omitempty"`
}

type Datatables struct {
	Id           string         `json:"id,omitempty"`
	Id_fk        string         `json:"id_fk,omitempty"`
	Id_sj        string         `json:"id_sj,omitempty"`
	InvoiceDate  string         `json:"invoice_date,omitempty"`
	Duration     string         `json:"duration,omitempty"`
	SjDate       string         `json:"sj_date,omitempty"`
	CustomerName string         `json:"customer,omitempty"`
	CompanyName  string         `json:"company,omitempty"`
	Shipto       string         `json:"shipto,omitempty"`
	NoPoCustomer string         `json:"po_customer,omitempty"`
	NoSj         string         `json:"no_sj,omitempty"`
	NoSo         string         `json:"no_so,omitempty"`
	NoInvoice    string         `json:"no_invoice,omitempty"`
	SendQty      string         `json:"send_qty,omitempty"`
	Item         string         `json:"item,omitempty"`
	Unit         string         `json:"unit,omitempty"`
	Ekspedisi    string         `json:"ekspedisi,omitempty"`
	Uom          string         `json:"uom,omitempty"`
	Jml          string         `json:"jml,omitempty"`
	Price        string         `json:"price,omitempty"`
	Bill         string         `json:"bill,omitempty"`
	Ppn          string         `json:"ppn,omitempty"`
	SubTotal     string         `json:"subtotal,omitempty"`
	Total        string         `json:"total,omitempty"`
	Cost         string         `json:"cost,omitempty"`
	PrintBy      string         `json:"print_by,omitempty"`
	InputBy      string         `json:"input_by,omitempty"`
	Address      string         `json:"address,omitempty"`
	SPhone       string         `json:"sphone,omitempty"`
	SName        string         `json:"sname,omitempty"`
	SAddress     string         `json:"saddress,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Billto       string         `json:"billto,omitempty"`
	Items        []InvoiceItems `json:"items,omitempty"`
	BankDetail   string         `json:"pilihBANK,omitempty"`
	BankAccount  string         `json:"an,omitempty"`
	BankNumber   string         `json:"rek,omitempty"`
	BankName     string         `json:"bank,omitempty"`
	Ttd          string         `json:"ttd,omitempty"`
	Logo         string         `json:"logo,omitempty"`
	Email        string         `json:"email,omitempty"`
	Note         string         `json:"note,omitempty"`
}

func Get(ctx *gin.Context) ([]Datatables, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	var querystr string
	LimitParam := ctx.DefaultQuery("limit", "10")
	OffsetParam := ctx.DefaultQuery("offset", "0")
	MonthlyReportParam := ctx.DefaultQuery("monthly_report", "")
	DataParam := ctx.DefaultQuery("data", "")
	DateNow := time.Now()
	DateNowFormat := DateNow.Format(config.App.DateFormat_Global)

	limit, err := strconv.Atoi(LimitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(OffsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	monthlyreport := MonthlyReportParam
	if MonthlyReportParam == "" {
		monthlyreport = DateNow.Format(config.App.DateFormat_MonthlyReport)
	}

	// Parse data request waiting, proccess, duedate to query
	if DataParam == "proccess" {
		querystr = fmt.Sprintf(`AND b.status = 0 AND b.duration >= '%s' AND b.id IS NOT NULL ORDER BY a.id DESC`, DateNowFormat)
	} else if DataParam == "duedate" {
		querystr = fmt.Sprintf(`AND b.status = 0 AND b.duration < '%s' AND b.id IS NOT NULL ORDER BY a.id DESC`, DateNowFormat)
	} else if DataParam == "done" {
		querystr = fmt.Sprintf(`AND b.status = 1 AND b.invoice_date LIKE '%s%%' ORDER BY a.id DESC`, monthlyreport)
	} else {
		querystr = `AND b.id IS NULL ORDER BY a.id DESC`
	}

	query := fmt.Sprintf(`SELECT a.id_fk, a.id_sj, a.no_delivery, a.send_qty, CASE WHEN b.id > 0 THEN b.id ELSE 0 END AS id, CASE WHEN b.print IS NOT NULL THEN b.print ELSE '' END AS print, CASE WHEN b.no_invoice IS NOT NULL THEN b.no_invoice ELSE '' END AS no_invoice, CASE WHEN b.invoice_date IS NOT NULL THEN b.invoice_date ELSE '' END AS invoice_date, CASE WHEN b.duration IS NOT NULL THEN b.duration ELSE '' END AS duration, CASE WHEN b.input_by > 0 THEN b.input_by ELSE 0 END AS input_by, CASE WHEN b.note IS NOT NULL THEN b.note ELSE '' END AS note, c.price, c.unit, d.no_so, e.customer, e.po_customer, f.sj_date, f.cost, f.ekspedisi, f.uom, f.jml, g.ppn, CASE WHEN i.name IS NOT NULL THEN i.name ELSE '' END AS name FROM delivery_orders_item AS a LEFT JOIN invoice AS b ON b.id_fk = a.id_fk AND b.id_sj = a.id_sj LEFT JOIN preorder_item AS c ON c.id_fk = a.id_fk AND c.item_to = a.item_to LEFT JOIN workorder_item AS d ON d.id_fk = a.id_fk AND d.item_to = a.item_to LEFT JOIN preorder_customer AS e ON e.id_fk = a.id_fk LEFT JOIN delivery_orders_customer AS f ON f.id_fk = a.id_fk AND f.id_sj = a.id_sj LEFT JOIN preorder_price AS g ON g.id_fk = a.id_fk LEFT JOIN status AS h ON h.id_fk = a.id_fk AND h.item_to = a.item_to LEFT JOIN user AS i ON i.id = b.input_by WHERE h.order_status BETWEEN 1 AND 2 %s LIMIT %d OFFSET %d`, querystr, limit, offset)

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	datatables := []Datatables{}
	arrInvoice := []string{}
	arrCost := []string{}

	for rows.Next() {
		var tagihanFormatted, ppnFormatted, totalFormatted, price, send_qty float64
		var id, id_fk, id_sj, input_by, cost, ppn int64
		var no_delivery, print, no_invoice, invoice_date, duration, unit, no_so, customername, nopocustomer, sj_date, ekspedisi, uom, jml, username, note string

		if err := rows.Scan(&id_fk, &id_sj, &no_delivery, &send_qty, &id, &print, &no_invoice, &invoice_date, &duration, &input_by, &note, &price, &unit, &no_so, &customername, &nopocustomer, &sj_date, &cost, &ekspedisi, &uom, &jml, &ppn, &username); err != nil {
			return nil, err
		}

		if send_qty > 0 {
			exNoSo := strings.Split(no_so, `/`)
			tagihanFormatted = price * send_qty
			if ppn > 0 {
				ppnFormatted = float64(tagihanFormatted * 11 / 10)
				totalFormatted = float64(tagihanFormatted + (tagihanFormatted * 11 / 10))
			} else {
				ppnFormatted = float64(0)
				totalFormatted = tagihanFormatted
			}

			// Filtering output print and username to except dupe
			if utils.InArray(arrInvoice, no_invoice) {
				print = ``
				username = ``
			} else {
				if print == `1` {
					print = `SUDAH`
				} else {
					print = `BELUM`
				}
			}

			// Filtering output cost to except dupe
			if utils.InArray(arrCost, fmt.Sprintf(`%d-%d-%d`, id_fk, id_sj, cost)) {
				cost = 0
			}

			datatables = append(datatables, Datatables{
				Id:           fmt.Sprintf(`%d-%d`, id_fk, id_sj),
				SjDate:       sj_date,
				InvoiceDate:  invoice_date,
				Duration:     duration,
				CustomerName: customername,
				NoPoCustomer: nopocustomer,
				NoSo:         fmt.Sprintf(`%s/%s%s`, exNoSo[0], exNoSo[1], exNoSo[2]),
				NoSj:         no_delivery,
				NoInvoice:    no_invoice,
				SendQty:      fmt.Sprintf(`%.2f`, send_qty),
				Unit:         unit,
				Price:        fmt.Sprintf(`%.2f`, price),
				Ekspedisi:    ekspedisi,
				Uom:          uom,
				Note:         note,
				Jml:          jml,
				Bill:         fmt.Sprintf(`%.2f`, tagihanFormatted),
				Ppn:          fmt.Sprintf(`%.2f`, ppnFormatted),
				Total:        fmt.Sprintf("%.2f", totalFormatted),
				Cost:         fmt.Sprintf("%d", cost),
				PrintBy:      print,
				InputBy:      username,
			})

			arrInvoice = append(arrInvoice, no_invoice)
			arrCost = append(arrCost, fmt.Sprintf(`%d-%d-%d`, id_fk, id_sj, cost))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return datatables, nil
}

func Create(Id string, Date string, InputBy int) ([]Datatables, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, fmt.Errorf("[err3] %s", err)
	}

	var invoice_date, no_invoice, invoice, customername, no_delivery string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// ambil podate dan no invoice pada tabel workorder_customer
	query := `SELECT invoice_date, no_invoice FROM invoice ORDER BY id DESC LIMIT 1`
	if err = sql.Connection.QueryRow(query).Scan(&invoice_date, &no_invoice); err != nil {
		if err.Error() == `sql: no rows in result set` {
		} else {
			return nil, err
		}
	}

	// Parse invoice date
	exInvoiceDate := strings.Split(invoice_date, `-`)
	InvoiceDateConv, err := strconv.Atoi(exInvoiceDate[0][2:len(exInvoiceDate[0])]) // output: 2024 -> 24
	if err != nil {
		return nil, fmt.Errorf("[err1] %s", err)
	}

	// check tahun saat ini
	timenow := time.Now()
	yearnow := timenow.Format(config.App.DateFormat_Year)
	yearnowconv, err := strconv.Atoi(yearnow) // output: 24
	if err != nil {
		return nil, fmt.Errorf("[err2] %s", err)
	}

	// check apakah nilai no invoice, invoice date, atau waktu ssat ini lebih besar dari invoice date
	if no_invoice == "" || invoice_date == "" || yearnowconv > InvoiceDateConv {
		invoice = fmt.Sprintf(`%d000001`, yearnowconv) // output: 24000001
	} else {
		NoInvoiceConv, err := strconv.Atoi(no_invoice)
		if err != nil {
			return nil, fmt.Errorf("[err3] %s", err)
		}

		invoice = fmt.Sprintf(`%06d`, (NoInvoiceConv + 1)) // output: 24000035 or 24000036 or 24000037 and more ....
	}

	// Generate estimasi +30 hari invoice
	DateParse, err := time.Parse(config.App.DateFormat_Global, Date)
	if err != nil {
		return nil, err
	}
	estDate := DateParse.AddDate(0, 0, 30)
	Duration := estDate.Format(config.App.DateFormat_Global)

	// Inisiasi
	arrayId := []string{}
	exId := strings.Split(Id, `,`)
	for _, value := range exId {
		parseValue := strings.Split(value, `-`)
		arrayId = append(arrayId, parseValue[0])
	}

	if len(utils.ArrayUnique(arrayId)) > 1 {
		return nil, fmt.Errorf("[err3] different Customer not allowed")
	} else {
		for _, value := range exId {
			exId = strings.Split(value, `-`)

			query = fmt.Sprintf(`SELECT a.customer, b.no_delivery FROM preorder_customer AS a LEFT JOIN delivery_orders_item AS b ON a.id_fk = b.id_fk WHERE b.id_fk = '%s' AND b.id_sj = '%s' GROUP BY b.id_sj`, exId[0], exId[1])

			if err = sql.Connection.QueryRow(query).Scan(&customername, &no_delivery); err != nil {
				if err.Error() == `sql: no rows in result set` {
				} else {
					return nil, err
				}
			}

			query = fmt.Sprintf(`INSERT INTO invoice (id_fk, id_sj, no_invoice, invoice_date, duration, input_by, print, print_date, status, complete_date, note) VALUES ('%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %d, '%s', '%s')`, exId[0], exId[1], invoice, Date, Duration, InputBy, 0, "0000-00-00", 0, "", "")

			create, err := sql.Connection.Query(query)
			if err != nil {
				return nil, err
			}
			defer create.Close()
		}
	}

	return []Datatables{}, nil
}

func Printview(Id int) ([]Datatables, error) {
	var no_invoice string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT no_invoice FROM invoice WHERE id = %d`, Id)
	if err = sql.Connection.QueryRow(query).Scan(&no_invoice); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("[err1] invalid ID")
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`SELECT a.id_fk, a.no_invoice, a.input_by, a.invoice_date, b.id_sj, b.item_to, b.no_delivery, b.send_qty, c.item, c.unit, c.price, d.customer, d.po_customer, e.ppn, f.shipto, f.cost, g.no_so, g.ingredient, g.size, g.volume, h.alamat, h.kota, h.negara, h.provinsi, h.kodepos, h.telp AS s_phone, h.s_nama, h.s_alamat, h.s_kota, h.s_negara, h.s_provinsi, h.s_kodepos, i.company, i.address, i.phone FROM invoice AS a LEFT JOIN delivery_orders_item AS b ON b.id_fk = a.id_fk AND b.id_sj = a.id_sj LEFT JOIN preorder_item AS c ON c.id_fk = a.id_fk AND c.item_to = b.item_to LEFT JOIN preorder_customer AS d ON d.id_fk = a.id_fk LEFT JOIN preorder_price AS e ON e.id_fk = a.id_fk LEFT JOIN delivery_orders_customer AS f ON f.id_fk = b.id_fk AND f.id_sj = b.id_sj LEFT JOIN workorder_item AS g ON g.id_fk = a.id_fk AND g.item_to = b.item_to LEFT JOIN customer AS h ON h.id = d.id_customer LEFT JOIN company AS i ON i.id = d.id_company WHERE a.no_invoice = '%s' ORDER BY b.no_delivery ASC`, no_invoice)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datatables := []Datatables{}
	arrCost := []string{}
	for rows.Next() {

		var tagihan, total, price, send_qty, ppn float64
		var id_fk, id_sj, input_by, cost, sequence_item int64
		var no_invoice, invoice_date, no_delivery, item, unit, customername, nopocustomer, shipto, no_so, material, size, volume, alamat, kota, negara, provinsi, kodepos, s_phone, s_nama, s_alamat, s_kota, s_negara, s_provinsi, s_kodepos, companyname, address, phone string

		if err := rows.Scan(&id_fk, &no_invoice, &input_by, &invoice_date, &id_sj, &sequence_item, &no_delivery, &send_qty, &item, &unit, &price, &customername, &nopocustomer, &ppn, &shipto, &cost, &no_so, &material, &size, &volume, &alamat, &kota, &negara, &provinsi, &kodepos, &s_phone, &s_nama, &s_alamat, &s_kota, &s_negara, &s_provinsi, &s_kodepos, &companyname, &address, &phone); err != nil {
			return nil, err
		}

		if send_qty > 0 {

			// Validation address
			if alamat != "" {
				alamat = fmt.Sprintf(`%s. `, alamat)
			}

			if kota != "" {
				kota = fmt.Sprintf(`%s - `, kota)
			}

			if provinsi != "" {
				provinsi = fmt.Sprintf(`%s, `, provinsi)
			}

			if negara != "" {
				negara = fmt.Sprintf(`%s. `, negara)
			}

			if s_kodepos != "" {
				kodepos = fmt.Sprintf(`%s. `, negara)
			}

			billto := fmt.Sprintf(`%s%s%s%s%s`, s_alamat, s_kota, s_provinsi, s_negara, s_kodepos)

			// Validation tax
			tagihan = send_qty * price
			if ppn > 0 {
				ppn = tagihan * 11 / 100
				total = tagihan + ppn
			} else {
				ppn = 0
				total = tagihan
			}

			// Filtering output cost to except dupe
			if utils.InArray(arrCost, fmt.Sprintf(`%d-%d-%d`, id_fk, id_sj, cost)) {
				cost = 0
			}

			// Parse nomor so
			exNoSo := strings.Split(no_so, `/`)

			datatables = append(datatables, Datatables{
				Id_fk:        fmt.Sprintf("%d", id_fk),
				CompanyName:  companyname,
				Address:      address,
				Phone:        phone,
				SPhone:       s_phone,
				Billto:       billto,
				SName:        s_nama,
				Shipto:       fmt.Sprintf(`%s%s`, shipto, s_phone),
				CustomerName: customername,
				NoPoCustomer: nopocustomer,
				NoSj:         no_delivery,
				NoInvoice:    no_invoice,
				Item:         strings.ToUpper(item),
				SendQty:      fmt.Sprintf("%.2f", send_qty),
				Total:        fmt.Sprintf("%.2f", total),
				Cost:         fmt.Sprintf("%d", cost),
				Price:        fmt.Sprintf("%.2f", price),
				Ppn:          fmt.Sprintf("%.2f", ppn),
				NoSo:         fmt.Sprintf(`%s/%s%s`, exNoSo[0], exNoSo[1], exNoSo[2]),
				InvoiceDate:  invoice_date,
				Unit:         unit,
			})

			arrCost = append(arrCost, fmt.Sprintf(`%d-%d-%d`, id_fk, id_sj, cost))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return datatables, nil
}

func Printnow(BodyReq []byte) ([]Datatables, error) {
	var data Datatables
	if err := json.Unmarshal([]byte(BodyReq), &data); err != nil {
		return nil, err
	}

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, fmt.Errorf("[err3] %s", err)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// ambil id_fk pada tabel workorder_customer
	var duration, logo, email string
	query := fmt.Sprintf(`SELECT a.duration, c.logo, c.email FROM invoice AS a LEFT JOIN preorder_customer AS b ON b.id_fk = a.id_fk LEFT JOIN company AS c ON c.id = b.id_company WHERE a.id = '%s' LIMIT 1`, data.Id)
	if err = sql.Connection.QueryRow(query).Scan(&duration, &logo, &email); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")
		} else {
			return nil, err
		}
	}

	if email != "" {
		email = fmt.Sprintf(`Email: %s`, email)
	}

	if data.BankDetail == "" {
		data.BankDetail = `BCA Cabang Pasar Baru-0021598690-Iskandar Zulkarnain`
	}
	exBank := strings.Split(data.BankDetail, `-`)

	// Generate estimasi +30 hari invoice
	durationParse, err := time.Parse(config.App.DateFormat_Global, duration)
	Duration := durationParse.Format(config.App.DateFormat_Print)
	if err != nil {
		return nil, err
	}

	datatables := Datatables{
		CompanyName:  data.CompanyName,
		Address:      data.Address,
		Phone:        data.Phone,
		CustomerName: data.CustomerName,
		NoInvoice:    data.NoInvoice,
		InvoiceDate:  data.InvoiceDate,
		SPhone:       data.SPhone,
		BankAccount:  exBank[2],
		BankNumber:   exBank[1],
		BankName:     exBank[0],
		NoSo:         data.NoSo,
		Billto:       data.Billto,
		Shipto:       data.Shipto,
		SName:        data.SName,
		Duration:     Duration,
		Logo:         logo,
		Email:        email,
	}

	var totalPpn, totalCost, subTotal float64
	for _, item := range data.Items {
		invoiceitems := InvoiceItems{
			NoSj:    item.NoSj,
			NoSo:    item.NoSo,
			Item:    item.Item,
			Unit:    item.Unit,
			SendQty: item.SendQty,
			Price:   item.Price,
		}
		subTotal += item.SendQty * item.Price
		totalPpn += item.Ppn
		totalCost += item.Cost

		datatables.Items = append(datatables.Items, invoiceitems)
	}

	datatables.SubTotal = fmt.Sprintf("%.2f", subTotal)
	datatables.Ppn = fmt.Sprintf("%.2f", totalPpn)
	datatables.Cost = fmt.Sprintf("%.2f", totalCost)
	datatables.Total = fmt.Sprintf("%.2f", totalPpn+subTotal)

	queryUpdate := fmt.Sprintf("UPDATE invoice SET print = 1, print_date ='%s' WHERE id = %s", data.InvoiceDate, data.Id)
	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}
	defer update.Close()

	// Log capture
	utils.Capture(
		`Print Invoice`,
		fmt.Sprintf(`Customer: %s - Invoice: %s - Total: Rp %.2f - Bank: %s`, data.CustomerName, data.NoInvoice, totalPpn+subTotal, data.BankDetail),
		data.Ttd,
	)

	return []Datatables{datatables}, nil
}

func Paid(id string, date string, note string) ([]Datatables, error) {
	var invoice_date, no_invoice, no_delivery, customername string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT a.invoice_date, a.no_invoice, b.customer, GROUP_CONCAT(DISTINCT c.no_delivery SEPARATOR ', ') AS no_delivery FROM invoice AS a LEFT JOIN preorder_customer AS b ON b.id_fk = a.id_fk LEFT JOIN delivery_orders_item AS c ON c.id_fk = a.id_fk AND c.id_sj = a.id_sj WHERE a.id = '%s' GROUP BY a.no_invoice`, id)
	if err := sql.Connection.QueryRow(query).Scan(&invoice_date, &no_invoice, &customername, &no_delivery); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`UPDATE invoice SET status = 1, complete_date = '%s', note = '%s' WHERE id = '%s'`, invoice_date, note, id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return []Datatables{}, nil

}

func UnPaid(id string) ([]Datatables, error) {
	var invoice_date, no_invoice, no_delivery, customername string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT a.invoice_date, a.no_invoice, b.customer, GROUP_CONCAT(DISTINCT c.no_delivery SEPARATOR ', ') AS no_delivery FROM invoice AS a LEFT JOIN preorder_customer AS b ON b.id_fk = a.id_fk LEFT JOIN delivery_orders_item AS c ON c.id_fk = a.id_fk AND c.id_sj = a.id_sj WHERE a.id = '%s' GROUP BY a.no_invoice`, id)
	if err := sql.Connection.QueryRow(query).Scan(&invoice_date, &no_invoice, &customername, &no_delivery); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`UPDATE invoice SET status = 0, complete_date = '', note = '' WHERE id = '%s'`, id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return []Datatables{}, nil

}

func Delete(id int) ([]Datatables, error) {
	var no_invoice string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT no_invoice FROM invoice WHERE id = %d`, id)
	if err := sql.Connection.QueryRow(query).Scan(&no_invoice); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`DELETE FROM invoice WHERE no_invoice = '%s'`, no_invoice)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return []Datatables{}, nil
}

package deliveryorder

import (
	"backend/adapters"
	"backend/config"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DeliveryOrders struct {
	Id           int64                `json:"id,omitempty"`
	CustomerName string               `json:"customer,omitempty"`
	CompanyName  string               `json:"company,omitempty"`
	Item         string               `json:"item,omitempty"`
	SequenceItem int64                `json:"item_to,omitempty"`
	NoSo         string               `json:"no_so,omitempty"`
	NoPoCustomer string               `json:"po_customer,omitempty"`
	ReqQty       int64                `json:"req_qty,omitempty"`
	Shipto       string               `json:"shipto,omitempty"`
	SpkDate      string               `json:"spk_date,omitempty"`
	Unit         string               `json:"unit,omitempty"`
	NoSj         string               `json:"no_sj,omitempty"`
	Resi         string               `json:"resi,omitempty"`
	SjDate       string               `json:"sj_date,omitempty"`
	Courier      string               `json:"courier,omitempty"`
	SendQty      int64                `json:"send_qty,omitempty"`
	Qty          int64                `json:"qty,omitempty"`
	NoSpk        string               `json:"no_spk,omitempty"`
	Cost         string               `json:"cost,omitempty"`
	InputBy      int64                `json:"input_by,omitempty"`
	Ttd          string               `json:"ttd,omitempty"`
	Address      string               `json:"address,omitempty"`
	Phone        string               `json:"phone,omitempty"`
	Logo         string               `json:"logo,omitempty"`
	Items        []DeliveryOrdersItem `json:"do_item,omitempty"`
}

type DeliveryOrdersItem struct {
	SequenceItem int   `json:"item_to,omitempty"`
	Qty          int64 `json:"qty,omitempty"`
}

type DataTables struct {
	Id           int    `json:"id"`
	SpkDate      string `json:"spk_date"`
	CustomerName string `json:"customer"`
	NoPoCustomer string `json:"po_customer"`
	Duration     string `json:"duration"`
	NoSo         string `json:"no_so"`
}

func Get(ctx *gin.Context) ([]DeliveryOrders, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	LimitParam := ctx.DefaultQuery("limit", "10")
	OffsetParam := ctx.DefaultQuery("offset", "0")
	MonthlyReportParam := ctx.DefaultQuery("monthly_report", "")

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
		DateNow := time.Now()
		monthlyreport = DateNow.Format(config.App.DateFormat_MonthlyReport)
	}

	query := fmt.Sprintf(`SELECT
	a.id,
	a.id_fk,
	a.id_sj,
	a.item_to,
	a.no_delivery,
	a.send_qty,
	b.order_status,
	c.sj_date,
	c.shipto,
	c.courier,
	c.no_tracking,
	c.cost,
	c.input_by,
	d.customer,
	d.po_customer,
	e.no_so,
	e.item,
	e.unit,
	f.name
	FROM
		delivery_orders_item AS a
	LEFT JOIN
		status AS b ON a.id_fk = b.id_fk AND a.item_to = b.item_to
	LEFT JOIN
		delivery_orders_customer AS c ON a.id_fk = c.id_fk AND a.id_sj = c.id_sj
	LEFT JOIN
		workorder_customer AS d ON a.id_fk = d.id_fk LEFT JOIN workorder_item AS e ON a.id_fk = e.id_fk AND a.item_to = e.item_to
	LEFT JOIN
		user AS f ON f.id = c.input_by
	WHERE
		b.order_status BETWEEN 1 AND 2 AND c.sj_date LIKE '%s%%'
	GROUP
		BY a.id
	ORDER BY
		a.id DESC LIMIT %d OFFSET %d`, monthlyreport, limit, offset)

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryorder := []DeliveryOrders{}
	for rows.Next() {
		var id, id_fk, id_sj, sequence_item, send_qty, order_status, input_by int64
		var customername, nopocustomer, no_sj, sj_date, shipto, item, unit, courier, resi, cost, no_so, username string

		if err := rows.Scan(&id, &id_fk, &id_sj, &sequence_item, &no_sj, &send_qty, &order_status, &sj_date, &shipto, &courier, &resi, &cost, &input_by, &customername, &nopocustomer, &no_so, &item, &unit, &username); err != nil {
			return nil, err
		}

		// Filter sj_date format
		sj_dateConv, err := time.Parse(config.App.DateFormat_Global, sj_date)
		if err != nil {
			return nil, err
		}
		sj_date = sj_dateConv.Format(config.App.DateFormat_Frontend)

		// Parse nomor SPK
		exNoSpk := strings.Split(no_so, "/")

		if send_qty > 0 {
			deliveryorder = append(deliveryorder, DeliveryOrders{
				Id:           id,
				CustomerName: customername,
				NoPoCustomer: nopocustomer,
				NoSpk:        fmt.Sprintf(`%s/%s%s`, exNoSpk[0], exNoSpk[1], exNoSpk[2]),
				NoSj:         no_sj,
				SjDate:       sj_date,
				Shipto:       shipto,
				Item:         item,
				SendQty:      send_qty,
				Unit:         unit,
				Courier:      courier,
				Resi:         resi,
				Cost:         cost,
				InputBy:      input_by,
			})
		}

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveryorder, nil
}

func Get_Waiting(ctx *gin.Context) ([]DataTables, error) {
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

	query := fmt.Sprintf(`SELECT
	a.id,
	a.spk_date,
	a.customer,
	a.po_customer,
	a.duration,
	GROUP_CONCAT(CONCAT(SUBSTRING_INDEX(c.no_so,'/',2),SUBSTRING_INDEX(c.no_so,'/',-1)) SEPARATOR ', ') AS no_so
	FROM
		workorder_customer AS a
	LEFT JOIN
		status AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		workorder_item AS c ON c.id_fk = a.id_fk AND c.item_to = b.item_to
	WHERE
		b.order_status BETWEEN 2 AND 3
	GROUP BY
		a.id_fk
	ORDER BY
		a.id ASC LIMIT %d OFFSET %d`, limit, offset)

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datatables := []DataTables{}
	for rows.Next() {
		var id int
		var spk_date, customername, nopocustomer, duration, no_so string

		if err := rows.Scan(&id, &spk_date, &customername, &nopocustomer, &duration, &no_so); err != nil {
			return nil, err
		}

		datatables = append(datatables, DataTables{
			Id:           id,
			SpkDate:      spk_date,
			CustomerName: customername,
			NoPoCustomer: nopocustomer,
			Duration:     duration,
			NoSo:         no_so,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return datatables, nil
}

func GetItem(Id int) ([]DeliveryOrders, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT
	a.spk_date,
	a.po_customer,
	a.customer,
	b.no_so,
	b.item,
	b.unit,
	b.qty,
	c.item_to,
	CASE
		WHEN d.total_send_qty > 0 THEN d.total_send_qty
		ELSE 0
		END
	AS total_send_qty,
	CASE
		WHEN e.shipto > 0 THEN e.shipto
		ELSE ''
		END
	AS shipto,
	g.alamat,
	g.kota,
	g.provinsi,
	g.negara,
	g.kodepos,
	g.s_alamat,
	g.s_kota,
	g.s_provinsi,
	g.s_negara,
	g.s_kodepos
	FROM
		workorder_customer AS a
	LEFT JOIN
		workorder_item AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		status AS c ON a.id_fk = c.id_fk AND b.item_to = c.item_to
	LEFT JOIN (
		SELECT
			y.id_fk,
			y.item_to,
			sum(y.send_qty) AS total_send_qty
		FROM
			workorder_customer AS x
		LEFT JOIN
			delivery_orders_item AS y ON x.id_fk = y.id_fk
		WHERE
			x.id = %d
		GROUP BY
			y.item_to
	) AS d ON a.id_fk = d.id_fk AND b.item_to = d.item_to
	LEFT JOIN (
		SELECT
			id_fk,
			shipto
		FROM
			delivery_orders_customer
		ORDER BY
			id
		DESC LIMIT 1
	) AS e ON a.id_fk = e.id_fk
	LEFT JOIN
		preorder_customer AS f ON f.id_fk = a.id_fk
	LEFT JOIN
		customer AS g ON g.id = f.id_customer
	WHERE
		a.id = %d AND c.order_status
	BETWEEN 2 AND 3`, Id, Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryorders := []DeliveryOrders{}
	for rows.Next() {
		var total_send_qty, sequence_item, qty int64
		var spk_date, nopocustomer, customer, so_no, shipto, item, unit, alamat, s_alamat, kota, s_kota, provinsi, s_provinsi, negara, s_negara, kodepos, s_kodepos string

		if err := rows.Scan(&spk_date, &nopocustomer, &customer, &so_no, &item, &unit, &qty, &sequence_item, &total_send_qty, &shipto, &alamat, &kota, &provinsi, &negara, &kodepos, &s_alamat, &s_kota, &s_provinsi, &s_negara, &s_kodepos); err != nil {
			return nil, err
		}

		// Validation address
		if shipto == "" {
			if s_alamat == "" {
				s_alamat = fmt.Sprintf(`%s. `, alamat)
			} else {
				s_alamat = fmt.Sprintf(`%s. `, s_alamat)
			}

			if s_kota == "" {
				s_kota = fmt.Sprintf(`%s - `, kota)
			} else {
				s_kota = fmt.Sprintf(`%s - `, s_kota)
			}

			if s_provinsi == "" {
				s_provinsi = fmt.Sprintf(`%s, `, provinsi)
			} else {
				s_kota = fmt.Sprintf(`%s, `, s_provinsi)
			}

			if s_negara == "" {
				s_negara = fmt.Sprintf(`%s. `, negara)
			} else {
				s_negara = fmt.Sprintf(`%s. `, s_negara)
			}

			if s_kodepos == "" {
				s_kodepos = kodepos
			}

			shipto = fmt.Sprintf(`%s%s%s%s%s`, s_alamat, s_kota, s_provinsi, s_negara, s_kodepos)
		}

		// Validate total send qty
		if total_send_qty > 0 {
			if total_send_qty > qty {
				total_send_qty = 0
			} else {
				total_send_qty = qty - total_send_qty
			}
		} else {
			total_send_qty = qty
		}

		// Parse nomor SO
		exNoSo := strings.Split(so_no, "/")

		deliveryorders = append(deliveryorders, DeliveryOrders{
			SpkDate:      spk_date,
			CustomerName: customer,
			NoPoCustomer: nopocustomer,
			NoSo:         fmt.Sprintf(`%s/%s%s`, exNoSo[0], exNoSo[1], exNoSo[2]),
			Item:         item,
			Unit:         unit,
			ReqQty:       total_send_qty,
			SequenceItem: sequence_item,
			Shipto:       shipto,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(deliveryorders) < 1 {
		return nil, errors.New("invalid ID")
	}

	return deliveryorders, nil
}

func GetNumber() ([]DeliveryOrders, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, fmt.Errorf("[err1] %s", err)
	}

	// ambil tanggal surat jalan terakhir pada tabel delivery_orders_customer
	var id_sj int
	var sj_date string

	query := `SELECT id_sj, sj_date FROM delivery_orders_customer ORDER BY id DESC LIMIT 1`
	err = sql.Connection.QueryRow(query).Scan(&id_sj, &sj_date)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			id_sj = 1
			sj_date = `1970-01-01`

		} else {
			return nil, fmt.Errorf("[err2] %s", err)
		}
	}

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, fmt.Errorf("[err3] %s", err)
	}

	// check tahun saat ini
	timenow := time.Now()
	yearnow := timenow.Format(config.App.DateFormat_Years)
	yearnowconv, err := strconv.Atoi(yearnow)
	if err != nil {
		return nil, fmt.Errorf("[err4] %s", err)
	}

	// parse tanggal surat jalan
	exSj_date := strings.Split(sj_date, `-`)
	sj_dateconv, err := strconv.Atoi(exSj_date[0])
	if err != nil {
		return nil, fmt.Errorf("[err5] %s", err)
	}

	// generate
	deliveryorder := []DeliveryOrders{}
	if id_sj < 1 || yearnowconv > sj_dateconv {
		yearupdate := timenow.Format(config.App.DateFormat_Year)
		deliveryorder = append(deliveryorder, DeliveryOrders{
			NoSj: fmt.Sprintf(`%s000001`, yearupdate),
		})

	} else {
		var no_delivery string

		yearupdate := timenow.Format(config.App.DateFormat_Year)

		query := `SELECT no_delivery FROM delivery_orders_item ORDER BY id DESC LIMIT 1`
		if err := sql.Connection.QueryRow(query).Scan(&no_delivery); err != nil {
			return nil, fmt.Errorf("[err6] %s", err)
		}

		exSj := strings.Split(no_delivery, `/`)
		substr := exSj[0][2:len(exSj[0])]
		noSj, err := strconv.Atoi(substr)
		if err != nil {
			return nil, fmt.Errorf("[err5] %s", err)
		}

		deliveryorder = append(deliveryorder, DeliveryOrders{
			NoSj: fmt.Sprintf(`%s%06d`, yearupdate, (noSj + 1)),
		})

	}

	return deliveryorder, nil
}

func Create(BodyReq []byte) ([]DeliveryOrders, error) {
	var id_fk, id_sj int

	var deliveryorder DeliveryOrders
	if err := json.Unmarshal([]byte(BodyReq), &deliveryorder); err != nil {
		return nil, err
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// ambil id_fk pada tabel workorder_customer
	query := fmt.Sprintf("SELECT id_fk FROM workorder_customer WHERE id = '%d' LIMIT 1", deliveryorder.Id)
	if err = sql.Connection.QueryRow(query).Scan(&id_fk); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")

		} else {
			return nil, err
		}
	}

	// ID validation
	if id_fk < 1 {
		return nil, fmt.Errorf("invalid id")
	}

	// check id_sj berdasarkan id_fk
	query = fmt.Sprintf("SELECT id_sj FROM delivery_orders_customer WHERE id_fk = %d ORDER BY id DESC LIMIT 1", deliveryorder.Id)
	err = sql.Connection.QueryRow(query).Scan(&id_sj)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			id_sj = 0
		} else {
			return nil, err
		}
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	// membuat kondisi jika id_sj berdasarkan id_fk ditemukan maka value id_sj + 1, sedangkan 0 maka akan diberi value 1
	if id_sj > 0 {
		id_sj = id_sj + 1
	} else {
		id_sj = 1
	}

	// input data ke tabel delivery_orders_customer
	query = fmt.Sprintf(`INSERT INTO delivery_orders_customer (id_fk, id_sj, sj_date, shipto, courier, no_tracking, cost, ekspedisi, uom, jml, input_by) VALUES (%d, %d, '%s', '%s', '%s', '%s', %d, '%s', '%s', '%s', %d)`, id_fk, id_sj, deliveryorder.SjDate, deliveryorder.Shipto, deliveryorder.Courier, deliveryorder.Resi, 0, "", "", "", deliveryorder.InputBy)
	if _, err = tx.Exec(query); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err1] %s", err)
	}

	// input data ke tabel delivery_orders_item
	for _, do := range deliveryorder.Items {
		var req_qty, send_qty, total, order_status int64

		query = fmt.Sprintf(`INSERT INTO delivery_orders_item (id_fk, id_sj, item_to, no_delivery, send_qty) VALUES (%d, %d, %d, '%s', %d)`, id_fk, id_sj, do.SequenceItem, deliveryorder.NoSj, do.Qty)
		if _, err = tx.Exec(query); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err)
		}

		// Menjumlah total tiap item yg dikirim berdasarkan id_fk
		query = fmt.Sprintf(`SELECT CASE WHEN a.qty > 0 THEN a.qty ELSE 0 END AS req_qty, CASE WHEN sum(b.send_qty) > 0 THEN sum(b.send_qty) ELSE 0 END AS send_qty FROM preorder_item AS a LEFT JOIN delivery_orders_item AS b ON a.id_fk = b.id_fk AND a.item_to = b.item_to WHERE a.id_fk = %d AND a.item_to = %d`, id_fk, do.SequenceItem)
		if err = sql.Connection.QueryRow(query).Scan(&req_qty, &send_qty); err != nil {
			return nil, err
		}

		// Ubah jika send qty sudah dikirim seadanya, mencukupi atau lebih dari req qty
		total = send_qty + do.Qty
		if total >= req_qty {
			order_status = 1
		} else {
			order_status = 2
		}
		query = fmt.Sprintf(`UPDATE status SET order_status = %d WHERE id_fk = %d AND item_to = %d`, order_status, id_fk, do.SequenceItem)
		if _, err = tx.Exec(query); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err3] %s", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return []DeliveryOrders{}, nil
}

func Delete(Id int) ([]DeliveryOrders, error) {
	var id_fk, id_sj, sequence_item, total_rows int64

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Mencari id_fk, id_sj, item_to
	query := fmt.Sprintf(`SELECT id_fk, id_sj, item_to FROM delivery_orders_item WHERE id = %d`, Id)
	if err = sql.Connection.QueryRow(query).Scan(&id_fk, &id_sj, &sequence_item); err != nil {
		return nil, err
	}

	// Menghitung total rows atas id_fk dan item_to
	query = fmt.Sprintf(`SELECT COUNT(id_fk) AS total_rows FROM delivery_orders_item WHERE id_fk = %d AND id_sj = %d GROUP BY id_fk`, id_fk, id_sj)
	if err = sql.Connection.QueryRow(query).Scan(&total_rows); err != nil {
		return nil, err
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	// Menghapus row
	query = fmt.Sprintf(`DELETE FROM delivery_orders_item WHERE id = %d`, Id)
	if _, err = tx.Exec(query); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err1] %s", err)
	}

	if total_rows < 2 {
		query = fmt.Sprintf(`DELETE FROM delivery_orders_customer WHERE id_fk = %d AND id_sj = %d`, id_fk, id_sj)
		if _, err = tx.Exec(query); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err)
		}

		query = fmt.Sprintf(`DELETE FROM invoice WHERE id_fk = %d AND id_sj = %d`, id_fk, id_sj)
		if _, err = tx.Exec(query); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err)
		}
	}

	// Mengubah status
	query = fmt.Sprintf(`UPDATE status SET order_status = 2 WHERE id_fk = %d AND item_to = %d`, id_fk, id_sj)
	if _, err = tx.Exec(query); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err3] %s", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return []DeliveryOrders{}, nil
}

func Printview(Id int) ([]DeliveryOrders, error) {
	var id_fk, id_sj, sequence_item int64

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Mencari id_fk, id_sj, item_to
	query := fmt.Sprintf(`SELECT id_fk, id_sj, item_to FROM delivery_orders_item WHERE id = %d`, Id)
	if err = sql.Connection.QueryRow(query).Scan(&id_fk, &id_sj, &sequence_item); err != nil {
		return nil, err
	}

	query = fmt.Sprintf(`SELECT a.no_delivery, a.send_qty, b.shipto, b.sj_date, c.customer, c.po_customer, e.item, e.unit, e.ingredient, e.size, e.volume FROM delivery_orders_item AS a LEFT JOIN delivery_orders_customer AS b ON a.id_fk = b.id_fk AND a.id_sj = b.id_sj LEFT JOIN workorder_customer AS c ON a.id_fk = c.id_fk LEFT JOIN workorder_item AS e ON a.id_fk = e.id_fk AND a.item_to = e.item_to WHERE a.id_fk = %d AND a.id_sj = %d`, id_fk, id_sj)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryorder := []DeliveryOrders{}
	for rows.Next() {
		var send_qty int64
		var no_sj, shipto, sj_date, customername, nopocustomer, item, unit, material, size, volume string

		if err := rows.Scan(&no_sj, &send_qty, &shipto, &sj_date, &customername, &nopocustomer, &item, &unit, &material, &size, &volume); err != nil {
			return nil, err
		}

		ItemUpperStr := strings.ToUpper(item)

		deliveryorder = append(deliveryorder, DeliveryOrders{
			CustomerName: customername,
			SjDate:       sj_date,
			Shipto:       shipto,
			NoSj:         no_sj,
			NoPoCustomer: nopocustomer,
			Item:         ItemUpperStr,
			Qty:          send_qty,
			Unit:         unit,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveryorder, nil
}

func Printnow(Id int, ttd string) ([]DeliveryOrders, error) {
	var id_fk, id_sj, sequence_item int64

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Mencari id_fk, id_sj, item_to
	query := fmt.Sprintf(`SELECT id_fk, id_sj, item_to FROM delivery_orders_item WHERE id = %d`, Id)
	if err = sql.Connection.QueryRow(query).Scan(&id_fk, &id_sj, &sequence_item); err != nil {
		return nil, err
	}

	query = fmt.Sprintf(`SELECT a.no_so, a.item, a.unit, a.ingredient, a.size, a.volume, b.send_qty, d.company, d.address, d.logo, d.phone, e.shipto, e.sj_date, c.po_customer, b.no_delivery, c.customer, f.name FROM workorder_item AS a LEFT JOIN delivery_orders_item AS b ON b.id_fk = a.id_fk AND b.item_to = a.item_to LEFT JOIN preorder_customer AS c ON c.id_fk = %d LEFT JOIN company AS d ON d.id = c.id_company LEFT JOIN delivery_orders_customer AS e ON a.id_fk = e.id_fk AND b.id_sj = e.id_sj LEFT JOIN user AS f ON c.input_by = f.id WHERE a.id_fk = %d AND b.id_sj = %d GROUP BY a.id;`, id_fk, id_fk, id_sj)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryorder := []DeliveryOrders{}
	for rows.Next() {
		var send_qty int64
		var no_so, item, unit, material, size, volume, company, address, logo, phone, shipto, sj_date, nopocustomer, no_sj, customername, ttd string

		if err := rows.Scan(&no_so, &item, &unit, &material, &size, &volume, &send_qty, &company, &address, &logo, &phone, &shipto, &sj_date, &nopocustomer, &no_sj, &customername, &ttd); err != nil {
			return nil, err
		}

		ItemUpperStr := strings.ToUpper(item)

		deliveryorder = append(deliveryorder, DeliveryOrders{
			Item:         ItemUpperStr,
			Qty:          send_qty,
			Unit:         unit,
			SjDate:       sj_date,
			NoPoCustomer: nopocustomer,
			NoSj:         no_sj,
			CustomerName: customername,
			Shipto:       shipto,
			Ttd:          ttd,
			NoSo:         no_so,
			CompanyName:  company,
			Address:      address,
			Phone:        phone,
			Logo:         logo,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveryorder, nil
}

package metrics

import (
	"backend/adapters"
	"backend/config"
	"backend/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Datatables struct {
	Item         string `json:"item"`
	Price        string `json:"price"`
	Qty          int    `json:"qty"`
	FkId         int    `json:"fkid"`
	OrderGrade   string `json:"order_grade"`
	NoSpk        string `json:"no_spk"`
	CustomerName string `json:"customer"`
	NoPoCustomer string `json:"po_customer"`
	PoDate       string `json:"po_date"`
	Ppn          string `json:"tax"`
	SoNumber     string `json:"so_no"`
	Size         string `json:"size"`
	Unit         string `json:"unit"`
	Qore         string `json:"qore"`
	Lin          string `json:"lin"`
	Roll         string `json:"roll"`
	Material     string `json:"ingredient"`
	Volume       string `json:"volume"`
	Porporasi    string `json:"porporasi"`
	Note         string `json:"annotation"`
	UkBahanBaku  string `json:"uk_bahan_baku"`
	QtyBahanBaku string `json:"qty_bahan_baku"`
	Sources      string `json:"sources"`
	Merk         string `json:"merk"`
	WoType       string `json:"type"`
	SpkDate      string `json:"spk_date"`
	NoSj         string `json:"sj_no"`
	SendQty      string `json:"send_qty"`
	SjId         string `json:"id_sj"`
	SjDate       string `json:"sj_date"`
	Courier      string `json:"courier"`
	OrderStatus  string `json:"order_status"`
	CompanyName  string `json:"company"`
	Isi          string `json:"isi"`
	Etd          string `json:"etd"`
	Total        string `json:"total"`
	SoDate       string `json:"so_date"`
	PriceBefore  string `json:"price_before"`
	Resi         string `json:"resi"`
	Cost         string `json:"cost"`
}

func Notification(Sessionid string) (interface{}, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	var role, account, do_waiting_counter, do_duedate_counter, wo_waiting_counter, inv_waiting_counter, inv_duedate_counter int

	query := fmt.Sprintf(`SELECT role, account FROM user WHERE id = %s`, Sessionid)
	if err = sql.Connection.QueryRow(query).Scan(&role, &account); err != nil {
		if err.Error() == `sql: no rows in result set` {
			role = 0
			account = 0
		} else {
			return nil, err
		}
	}

	data := []map[string]int{}
	Now := time.Now()
	Date := Now.Format(config.App.DateFormat_Global)

	query1 := fmt.Sprintf(`SELECT count(CASE WHEN a.duration >= '%s' THEN 1 END) AS do_waiting_counter, count(CASE WHEN a.duration < '%s' THEN 1 END) AS do_duedate_counter FROM workorder_customer AS a JOIN status AS b ON a.id_fk = b.id_fk WHERE b.order_status BETWEEN 3 AND 15;`, Date, Date)
	if err = sql.Connection.QueryRow(query1).Scan(&do_waiting_counter, &do_duedate_counter); err != nil {
		if err.Error() == `sql: no rows in result set` {
			do_waiting_counter = 0
			do_duedate_counter = 0
		} else {
			return nil, err
		}
	}

	query2 := `SELECT count(id_fk) as wo_waiting_counter FROM status WHERE order_status = 0`
	if err = sql.Connection.QueryRow(query2).Scan(&wo_waiting_counter); err != nil {
		if err.Error() == `sql: no rows in result set` {
			wo_waiting_counter = 0
		} else {
			return nil, err
		}
	}

	query3 := `SELECT COUNT(sub.count) AS inv_waiting_counter FROM (SELECT CASE WHEN a.id_fk > 0 THEN 1 ELSE 0 END as count FROM delivery_orders_customer AS a JOIN status AS b ON a.id_fk = b.id_fk JOIN invoice AS c ON a.id_fk = c.id_fk WHERE b.order_status BETWEEN 1 AND 2 AND c.id_fk IS NOT NULL GROUP BY a.id_fk) AS sub`
	if err = sql.Connection.QueryRow(query3).Scan(&inv_waiting_counter); err != nil {
		if err.Error() == `sql: no rows in result set` {
			inv_waiting_counter = 0
		} else {
			return nil, err
		}
	}

	query4 := fmt.Sprintf(`SELECT COUNT(sub.count) AS inv_duedate_counter FROM (SELECT CASE WHEN a.id_fk > 0 THEN 1 ELSE 0 END AS count FROM delivery_orders_customer AS a JOIN status AS b ON a.id_fk = b.id_fk JOIN invoice AS c ON a.id_fk = c.id_fk WHERE b.order_status BETWEEN 1 AND 2 AND c.id_fk IS NOT NULL AND c.duration < '%s' GROUP BY a.id_fk) AS sub`, Date)
	if err = sql.Connection.QueryRow(query4).Scan(&inv_duedate_counter); err != nil {
		if err.Error() == `sql: no rows in result set` {
			inv_duedate_counter = 0
		} else {
			return nil, err
		}
	}

	data = append(data, map[string]int{
		"role":                role,
		"account":             account,
		"total_counter":       do_waiting_counter + do_duedate_counter + wo_waiting_counter + inv_waiting_counter + inv_duedate_counter,
		"do_waiting_counter":  do_waiting_counter,
		"do_duedate_counter":  do_duedate_counter,
		"wo_waiting_counter":  wo_waiting_counter,
		"inv_waiting_counter": inv_waiting_counter,
		"inv_duedate_counter": inv_duedate_counter,
	})

	return data, nil
}

func SoTracking(ctx *gin.Context) ([]Datatables, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	var totalrows int
	var query, search, query_datatables, Report string
	// idParam := ctx.DefaultQuery("id", "0")
	LimitParam := ctx.DefaultQuery("limit", "10")
	OffsetParam := ctx.DefaultQuery("offset", "0")
	ReportParam := ctx.DefaultQuery("report", "")
	StartDateParam := ctx.DefaultQuery("startdate", "")
	EndDateParam := ctx.DefaultQuery("enddate", "")

	// datatables handling
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
	if err != nil || offset < -1 {
		offset = 0
	}

	// filter data based on monthly, yearly or periode
	DateNow := time.Now()
	MonthDefault := DateNow.Format(config.App.DateFormat_MonthlyReport)

	if ReportParam == "month" && StartDateParam != "" && utils.ValidateReportFormatDate(strings.ReplaceAll(StartDateParam, "-", "/"), config.App.DateFormat_MonthlyReport) {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01' AND '%s-31'`,
			strings.ReplaceAll(StartDateParam, "/", "-"),
			strings.ReplaceAll(StartDateParam, "/", "-")) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else if ReportParam == "year" && StartDateParam != "" && utils.ValidateReportFormatDate(StartDateParam, config.App.DateFormat_Years) {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01-01' AND '%s-12-31'`,
			StartDateParam,
			StartDateParam) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else if ReportParam == "periode" && StartDateParam != "" && EndDateParam != "" && utils.ValidateReportFormatDate(StartDateParam, config.App.DateFormat_Global) && utils.ValidateReportFormatDate(EndDateParam, config.App.DateFormat_Global) {
		Report = fmt.Sprintf(
			`BETWEEN '%s' AND '%s'`,
			StartDateParam,
			EndDateParam) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01' AND '%s-31'`,
			strings.ReplaceAll(MonthDefault, "/", "-"),
			strings.ReplaceAll(MonthDefault, "/", "-")) // BETWEEN '2024-01-01' AND '2024-01-31'
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_datatables = fmt.Sprintf(`SELECT COUNT(b.id) as totalrows FROM preorder_item AS a LEFT JOIN preorder_customer AS b ON b.id_fk = a.id_fk LEFT JOIN preorder_price AS c ON c.id_fk = b.id_fk LEFT JOIN workorder_item AS d ON d.id_fk = a.id_fk AND d.item_to = a.item_to LEFT JOIN workorder_customer AS e ON e.id_fk = b.id_fk LEFT JOIN delivery_orders_item AS f ON f.id_fk = a.id_fk AND f.item_to = a.item_to LEFT JOIN delivery_orders_customer AS g ON g.id_fk = b.id_fk AND g.id_sj = f.id_sj LEFT JOIN status AS h ON h.id_fk = a.id_fk AND h.item_to = a.item_to LEFT JOIN company AS i ON i.id = b.id_company LEFT JOIN setting AS j ON j.id = d.detail WHERE b.po_date %s %s ORDER BY b.id, f.no_delivery ASC`, Report, search)
	if err = sql.Connection.QueryRow(query_datatables).Scan(&totalrows); err != nil {
		if err.Error() == `sql: no rows in result set` {
			totalrows = 0
		} else {
			return nil, err
		}
	}

	// If request limit -1 (pagination datatables) is show all
	// Based on monthlyreport
	if limit == -1 {
		limit = totalrows
	}

	query = fmt.Sprintf(`SELECT
	a.item,
	a.price,
	a.qty AS req_qty,
	a.unit,
	b.id_fk,
	b.order_grade,
	b.customer,
	b.po_customer,
	b.po_date,
	c.ppn,
	d.no_so,
	d.size,
	d.qore,
	d.lin,
	d.roll,
	d.ingredient,
	d.volume,
	d.porporasi,
	d.annotation,
	d.uk_bahan_baku,
	d.qty_bahan_baku,
	d.sources,
	d.merk,
	d.type,
	e.spk_date,
	CASE WHEN f.no_delivery IS NOT NULL THEN f.no_delivery ELSE '' END AS no_delivery,
	CASE WHEN f.send_qty IS NOT NULL THEN f.send_qty ELSE '' END AS send_qty,
	CASE WHEN g.id_sj IS NOT NULL THEN g.id_sj ELSE '' END AS id_sj,
	CASE WHEN g.sj_date IS NOT NULL THEN g.sj_date ELSE '0000-00-00' END AS sj_date,
	CASE WHEN g.courier IS NOT NULL THEN g.courier ELSE '' END AS courier,
	CASE WHEN g.no_tracking IS NOT NULL THEN g.no_tracking ELSE '' END AS no_tracking,
	CASE WHEN g.cost IS NOT NULL THEN g.cost ELSE 0 END AS cost,
	h.order_status,
	i.company,
	CASE WHEN j.isi IS NOT NULL THEN j.isi ELSE '' END AS isi,
	(a.price * a.qty) AS etd
	FROM
		preorder_item AS a
	LEFT JOIN
		preorder_customer AS b ON b.id_fk = a.id_fk
	LEFT JOIN
		preorder_price AS c ON c.id_fk = b.id_fk
	LEFT JOIN
		workorder_item AS d ON d.id_fk = a.id_fk AND d.item_to = a.item_to
	LEFT JOIN
		workorder_customer AS e ON e.id_fk = b.id_fk
	LEFT JOIN
		delivery_orders_item AS f ON f.id_fk = a.id_fk AND f.item_to = a.item_to
	LEFT JOIN
		delivery_orders_customer AS g ON g.id_fk = b.id_fk AND g.id_sj = f.id_sj
	LEFT JOIN
		status AS h ON h.id_fk = a.id_fk AND h.item_to = a.item_to
	LEFT JOIN
		company AS i ON i.id = b.id_company
	LEFT JOIN
		setting AS j ON j.id = d.detail
	WHERE b.po_date %s %s ORDER BY b.id, f.no_delivery ASC LIMIT %d OFFSET %d`, Report, search, limit, offset)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	datatables := []Datatables{}
	arrCost := []string{}
	for rows.Next() {
		var etd, ppn, total, cost float64
		var fkid, qty, order_grade int
		var item, price, order_gradeStr, customername, nopo, po_date, so_no, size, unit, qore, lin, roll, material, volume, porporasi, note, uk_bahan_baku, qty_bahan_baku, sources, merk, itemtype, spk_date, sj_no, send_qty, id_sj, sj_date, courier, resi, orderStatus, companyname, isi string

		if err := rows.Scan(&item, &price, &qty, &unit, &fkid, &order_grade, &customername, &nopo, &po_date, &ppn, &so_no, &size, &qore, &lin, &roll, &material, &volume, &porporasi, &note, &uk_bahan_baku, &qty_bahan_baku, &sources, &merk, &itemtype, &spk_date, &sj_no, &send_qty, &id_sj, &sj_date, &courier, &resi, &cost, &orderStatus, &companyname, &isi, &etd); err != nil {
			return nil, err
		}

		// Parsing nomor so untuk generate nomor spk
		noSpk := strings.Split(so_no, "/")

		if order_grade > 0 {
			order_gradeStr = `Spesial`
		} else {
			order_gradeStr = `Reguler`
		}

		// Filter Order Status format
		if orderStatus == "0" {
			orderStatus = `PO baru dibuat`
		} else if orderStatus == "1" || orderStatus == "2" {
			orderStatus = `Delivery`
		} else if orderStatus == "3" {
			orderStatus = `Packing`
		} else if orderStatus == "4" {
			orderStatus = `Cetak SPK`

		} else if orderStatus == "5" {
			orderStatus = `Pembuatan Pisau`

		} else if orderStatus == "6" {
			orderStatus = `Antri Sliting`

		} else if orderStatus == "7" {
			orderStatus = `Antri Cetak`

		} else if orderStatus == "8" {
			orderStatus = `Proses Cetak`

		} else if orderStatus == "9" {
			orderStatus = `Proses Bahan Baku`

		} else if orderStatus == "10" {
			orderStatus = `Proses Film`

		} else if orderStatus == "11" {
			orderStatus = `Proses Toyobo`

		} else if orderStatus == "12" {
			orderStatus = `Proses ACC`

		} else if orderStatus == "13" {
			orderStatus = `Proses Sliting`

		} else if orderStatus == "14" {
			orderStatus = `Reture`

		} else if orderStatus == "15" {
			orderStatus = `Proses Sample`

		} else if orderStatus == "16" {
			orderStatus = `Input PO`

		} else {
			orderStatus = `PO baru dibuat`
		}

		// Filter Porporasi format
		if porporasi == "1" {
			porporasi = "Ya"
		} else {
			porporasi = "Tidak"
		}

		// Filter SPK date format
		if spk_date == "1970-01-01" || spk_date == "0000-00-00" {
			spk_date = ``
		} else {
			SpkDateParse, err := time.Parse(config.App.DateFormat_Global, spk_date)
			if err != nil {
				return nil, err
			}

			spk_date = SpkDateParse.Format(config.App.DateFormat_Frontend)
		}

		// Filter Po date format
		if po_date == "1970-01-01" || po_date == "0000-00-00" {
			po_date = ``
		} else {
			po_dateParse, err := time.Parse(config.App.DateFormat_Global, po_date)
			if err != nil {
				return nil, err
			}

			po_date = po_dateParse.Format(config.App.DateFormat_Frontend)
		}

		// Filter Sj date format
		if sj_date == "1970-01-01" || sj_date == "0000-00-00" {
			sj_date = ``
		} else {
			sj_dateParse, err := time.Parse(config.App.DateFormat_Global, sj_date)
			if err != nil {
				return nil, err
			}

			sj_date = sj_dateParse.Format(config.App.DateFormat_Frontend)
		}

		// Parsing dan filter sources
		ParseSources := strings.Split(sources, "|")
		if ParseSources[0] == "2" {
			sourcesDateParse, err := time.Parse(config.App.DateFormat_Global, ParseSources[2])
			if err != nil {
				return nil, err
			}
			sourcesStr := sourcesDateParse.Format(config.App.DateFormat_Frontend)
			sources = fmt.Sprintf(`SUBCONT (%s, %s)`, ParseSources[1], sourcesStr)

		} else if ParseSources[0] == "3" {
			sources = fmt.Sprintf(`IN STOCK (%s %s)`, ParseSources[1], unit)

		} else {
			sources = `Internal`
		}

		if ppn > 0 {
			ppn = etd * 11 / 100
			total = etd + ppn
		} else {
			ppn = 0
			total = etd
		}

		// Filtering output cost to except dupe
		if utils.InArray(arrCost, fmt.Sprintf(`%d-%s-%.2f`, fkid, id_sj, cost)) {
			cost = 0
		}

		datatables = append(datatables, Datatables{
			Item:         item,
			Price:        price,
			Qty:          qty,
			FkId:         fkid,
			OrderGrade:   order_gradeStr,
			CustomerName: customername,
			NoPoCustomer: nopo,
			PoDate:       po_date,
			Size:         size,
			Unit:         unit,
			Qore:         qore,
			Lin:          lin,
			Roll:         roll,
			Material:     material,
			Volume:       volume,
			Porporasi:    porporasi,
			Note:         note,
			UkBahanBaku:  uk_bahan_baku,
			QtyBahanBaku: qty_bahan_baku,
			Sources:      sources,
			Merk:         merk,
			WoType:       itemtype,
			SpkDate:      spk_date,
			NoSj:         sj_no,
			SendQty:      send_qty,
			SjId:         id_sj,
			SjDate:       sj_date,
			Courier:      courier,
			OrderStatus:  orderStatus,
			CompanyName:  companyname,
			Isi:          isi,
			NoSpk:        fmt.Sprintf(`%s/%s%s`, noSpk[0], noSpk[1], noSpk[2]),
			Ppn:          fmt.Sprintf("%.2f", ppn),
			Etd:          fmt.Sprintf("%.2f", etd),
			Total:        fmt.Sprintf("%.2f", total),
			Resi:         resi,
			Cost:         fmt.Sprintf("%.2f", cost),
			SoDate:       po_date,
			PriceBefore:  fmt.Sprintf("%.2f", etd),
		})

		arrCost = append(arrCost, fmt.Sprintf(`%d-%s-%.2f`, fkid, id_sj, cost))
	}

	return datatables, nil
}

func Static(ctx *gin.Context) (interface{}, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	ReportParam := ctx.DefaultQuery("report", "")
	StartDateParam := ctx.DefaultQuery("startdate", "")
	EndDateParam := ctx.DefaultQuery("enddate", "")

	var Report string
	var po_total, wo_total, do_total, inv_total int
	data := []map[string]int{}

	// filter data based on monthly, yearly or periode
	DateNow := time.Now()
	MonthDefault := DateNow.Format(config.App.DateFormat_MonthlyReport)

	if ReportParam == "month" && StartDateParam != "" && utils.ValidateReportFormatDate(strings.ReplaceAll(StartDateParam, "-", "/"), config.App.DateFormat_MonthlyReport) {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01' AND '%s-31'`,
			strings.ReplaceAll(StartDateParam, "/", "-"),
			strings.ReplaceAll(StartDateParam, "/", "-")) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else if ReportParam == "year" && StartDateParam != "" && utils.ValidateReportFormatDate(StartDateParam, config.App.DateFormat_Years) {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01-01' AND '%s-12-31'`,
			StartDateParam,
			StartDateParam) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else if ReportParam == "periode" && StartDateParam != "" && EndDateParam != "" && utils.ValidateReportFormatDate(StartDateParam, config.App.DateFormat_Global) && utils.ValidateReportFormatDate(EndDateParam, config.App.DateFormat_Global) {
		Report = fmt.Sprintf(
			`BETWEEN '%s' AND '%s'`,
			StartDateParam,
			EndDateParam) // BETWEEN '2024-01-01' AND '2024-01-31'

	} else {
		Report = fmt.Sprintf(
			`BETWEEN '%s-01' AND '%s-31'`,
			strings.ReplaceAll(MonthDefault, "/", "-"),
			strings.ReplaceAll(MonthDefault, "/", "-")) // BETWEEN '2024-01-01' AND '2024-01-31'
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query1 := fmt.Sprintf(`SELECT count(id) AS po_total FROM preorder_customer WHERE po_date %s`, Report)
	if err = sql.Connection.QueryRow(query1).Scan(&po_total); err != nil {
		if err.Error() == `sql: no rows in result set` {
			po_total = 0
		} else {
			return nil, err
		}
	}

	query2 := fmt.Sprintf(`SELECT COUNT(sub.jml_wo) AS wo_total FROM (SELECT CASE WHEN a.id > 0 THEN 1 END AS jml_wo FROM workorder_customer AS a LEFT JOIN status AS b ON a.id_fk = b.id_fk WHERE b.order_status BETWEEN 4 AND 13 AND a.spk_date %s GROUP BY a.id_fk) AS sub`, Report)
	if err = sql.Connection.QueryRow(query2).Scan(&wo_total); err != nil {
		if err.Error() == `sql: no rows in result set` {
			wo_total = 0
		} else {
			return nil, err
		}
	}

	query3 := fmt.Sprintf(`SELECT COUNT(sub.jml_do) AS do_total FROM (SELECT CASE WHEN a.id > 0 THEN 1 END AS jml_do FROM delivery_orders_customer AS a LEFT JOIN status AS b ON a.id_fk = b.id_fk WHERE b.order_status BETWEEN 1 AND 2 AND a.sj_date %s GROUP BY a.id_fk) AS sub`, Report)
	if err = sql.Connection.QueryRow(query3).Scan(&do_total); err != nil {
		if err.Error() == `sql: no rows in result set` {
			inv_total = 0
		} else {
			return nil, err
		}
	}

	query4 := fmt.Sprintf(`SELECT count(DISTINCT no_invoice) as inv_total FROM invoice WHERE status = 1 AND invoice_date %s`, Report)
	if err = sql.Connection.QueryRow(query4).Scan(&inv_total); err != nil {
		if err.Error() == `sql: no rows in result set` {
			inv_total = 0
		} else {
			return nil, err
		}
	}

	data = append(data, map[string]int{
		"po_total":  po_total,
		"wo_total":  wo_total,
		"do_total":  do_total,
		"inv_total": inv_total,
	})

	return data, nil
}

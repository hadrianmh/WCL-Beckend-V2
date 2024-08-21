package salesorder

import (
	"backend/adapters"
	"backend/config"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data            []DataTables `json:"data"`
	RecordsTotal    string       `json:"recordsTotal,omitempty"`
	RecordsFiltered string       `json:"recordsFiltered,omitempty"`
}

type Preorder struct {
	Id           string         `json:"poid"`
	Fkid         string         `json:"fkid"`
	CompanyId    int            `json:"companyid"`
	CompanyName  string         `json:"company,omitempty"`
	CustomerId   int            `json:"customerid,omitempty"`
	CustomerName string         `json:"customer"`
	PoDate       string         `json:"po_date"`
	NoPoCustomer string         `json:"po_customer,omitempty"`
	OrderGrade   int            `json:"order_grade"`
	InputBy      int            `json:"input_by,omitempty"`
	Ppn          int            `json:"tax"`
	Items        []PreorderItem `json:"items,omitempty"`
}

type PreorderItem struct {
	Poitemid     int    `json:"poitemid,omitempty"`
	Woitemid     int    `json:"woitemid,omitempty"`
	SequenceItem int    `json:"sequence_item,omitempty"`
	Detail       int    `json:"detail"`
	Item         string `json:"item"`
	Size         string `json:"size"`
	Price        string `json:"price"`
	Qty          int64  `json:"qty"`
	SoNumber     string `json:"so_no,omitempty"`
	Unit         string `json:"unit"`
	Qore         string `json:"qore"`
	Lin          string `json:"lin"`
	Roll         string `json:"roll"`
	Material     string `json:"ingredient"`
	Volume       int64  `json:"volume"`
	Total        int    `json:"total,omitempty"`
	Note         string `json:"annotation"`
	Porporasi    int    `json:"porporasi"`
	UkBahanBaku  string `json:"uk_bahan_baku"`
	QtyBahanBaku string `json:"qty_bahan_baku"`
	Sources      string `json:"sources"`
	Merk         string `json:"merk"`
	WoType       string `json:"type"`
	Etc1         string `json:"etc1,omitempty"`
	Etc2         string `json:"etc2,omitempty"`
	InputBy      int    `json:"input_by,omitempty"`
}

type PreorderShippingCost struct {
	Id        int    `json:"id"`
	Detail    string `json:"detail"`
	Cost      string `json:"cost"`
	Ekspedisi string `json:"ekspedisi"`
	Uom       string `json:"uom"`
	Jml       string `json:"jml"`
}

type DataTables struct {
	Itemid       int    `json:"itemid"`
	SequenceItem int    `json:"sequence_item"`
	Detail       string `json:"detail"`
	Item         string `json:"item"`
	Size         string `json:"size"`
	Price        string `json:"price"`
	Qty          int    `json:"qty"`
	CompanyId    int    `json:"companyid"`
	CustomerId   int    `json:"customerid"`
	FkId         int    `json:"fkid"`
	CustomerName string `json:"customer"`
	Estimated    string `json:"estimasi"`
	PoDate       string `json:"po_date"`
	NoPoCustomer string `json:"po_customer"`
	OrderGrade   string `json:"order_grade"`
	InputBy      int    `json:"input_by"`
	SoNumber     string `json:"so_no"`
	Unit         string `json:"unit"`
	Qore         string `json:"qore"`
	Lin          string `json:"lin"`
	Roll         string `json:"roll"`
	Material     string `json:"ingredient"`
	Volume       string `json:"volume"`
	Total        string `json:"total"`
	Note         string `json:"annotation"`
	Porporasi    string `json:"porporasi"`
	UkBahanBaku  string `json:"uk_bahan_baku"`
	QtyBahanBaku string `json:"qty_bahan_baku"`
	Sources      string `json:"sources"`
	Merk         string `json:"merk"`
	WoType       string `json:"type"`
	Ppn          string `json:"ppn"`
	OrderStatus  int    `json:"order_status"`
	Etd          string `json:"etd"`
	Isi          string `json:"isi"`
	Ongkir       string `json:"ongkir"`
	SjId         string `json:"id_sj"`
	CompanyName  string `json:"company"`
	InputName    string `json:"username"`
}

type SuggestionsType struct {
	Id   string `json:"id,omitempty"`
	Item string `json:"item,omitempty"`
}

type SuggestionsCustomer struct {
	Poid         string `json:"po_id,omitempty"`
	Customerid   string `json:"customerid"`
	CustomerName string `json:"customername,omitempty"`
	NoPoCustomer string `json:"nopocustomer,omitempty"`
	Fkid         string `json:"fkid"`
	Item         string `json:"item,omitempty"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Category     string `json:"category"`
}

type SuggestionsItem struct {
	PoId       int            `json:"po_id,omitempty"`
	Customerid int            `json:"customerid,omitempty"`
	Items      []PreorderItem `json:"items,omitempty"`
}

func Get(ctx *gin.Context) (Response, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return Response{}, err
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
	if err != nil || offset < 1 {
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

	// datatables total rows and filtered handling
	if SearchValue != "" {
		search = fmt.Sprintf(`AND (a.detail LIKE '%%%s%%' OR a.item LIKE '%%%s%%' OR a.size LIKE '%%%s%%' OR a.price LIKE '%%%s%%' OR a.qty LIKE '%%%s%%' OR a.unit LIKE '%%%s%%' OR b.customer LIKE '%%%s%%' OR b.po_date LIKE '%%%s%%' OR b.po_customer LIKE '%%%s%%' OR c.no_so LIKE '%%%s%%' OR c.qore LIKE '%%%s%%' OR c.lin LIKE '%%%s%%' OR c.roll LIKE '%%%s%%' OR c.ingredient LIKE '%%%s%%' OR c.volume LIKE '%%%s%%' OR c.annotation LIKE '%%%s%%' OR c.uk_bahan_baku LIKE '%%%s%%' OR c.qty_bahan_baku LIKE '%%%s%%' OR c.merk LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return Response{}, err
	}

	query_datatables = fmt.Sprintf(`SELECT
	COUNT(a.id) as totalrows FROM preorder_item AS a
	LEFT JOIN
		preorder_customer AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		workorder_item AS c ON c.id_fk = a.id_fk AND c.item_to = a.item_to
	LEFT JOIN
		preorder_price AS d ON a.id_fk = d.id_fk
	LEFT JOIN
		status AS e ON a.id_fk = e.id_fk AND a.item_to = e.item_to
	LEFT JOIN
		setting AS f ON a.detail = f.id
	LEFT JOIN
		(SELECT id_fk, id_sj, SUM(cost) AS total_ongkir FROM delivery_orders_customer GROUP BY id_fk) AS g ON g.id_fk = b.id_fk
	LEFT JOIN
		company AS h ON h.id = b.id_company WHERE b.po_date %s %s ORDER BY a.id DESC`, Report, search)
	if err = sql.Connection.QueryRow(query_datatables).Scan(&totalrows); err != nil {
		if err.Error() == `sql: no rows in result set` {
			totalrows = 0
		} else {
			return Response{}, err
		}
	}

	// If request limit -1 (pagination datatables) is show all
	// Based on monthlyreport
	if limit == -1 {
		limit = totalrows
	}

	query = fmt.Sprintf(`SELECT
		a.id AS itemid,
		a.item_to,
		a.detail,
		a.item,
		a.size,
		a.price,
		a.qty,
		a.unit,
		b.id_company AS companyid,
		b.id AS customerid,
		b.id_fk,
		b.customer,
		b.po_date,
		b.order_grade,
		b.po_customer,
		b.input_by,
		c.no_so,
		c.qore,
		c.lin,
		c.roll,
		c.ingredient,
		c.volume,
		c.annotation,
		c.porporasi,
		c.uk_bahan_baku,
		c.qty_bahan_baku,
		c.sources,
		c.merk,
		c.type,
		d.ppn,
		e.order_status,
		(a.price * a.qty) AS etd,
		CASE WHEN f.isi != '' THEN f.isi ELSE '' END AS isi,
		coalesce(g.total_ongkir, "0") AS ongkir,
		coalesce(g.id, "0") AS doid,
		coalesce(g.id_sj, "0") AS id_sj,
		h.company
	FROM (SELECT * FROM preorder_item WHERE hidden = 0) AS a
	LEFT JOIN
		preorder_customer AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		workorder_item AS c ON c.id_fk = a.id_fk AND c.item_to = a.item_to
	LEFT JOIN
		preorder_price AS d ON a.id_fk = d.id_fk
	LEFT JOIN
		status AS e ON a.id_fk = e.id_fk AND a.item_to = e.item_to
	LEFT JOIN
		setting AS f ON a.detail = f.id
	LEFT JOIN
		(SELECT id, id_fk, id_sj, SUM(cost) AS total_ongkir FROM delivery_orders_customer GROUP BY id_fk) AS g ON g.id_fk = b.id_fk
	LEFT JOIN
		company AS h ON h.id = b.id_company
	WHERE b.po_date %s %s ORDER BY a.id DESC LIMIT %d OFFSET %d`, Report, search, limit, offset)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	datatables := []DataTables{}
	arrCost := []string{}
	for rows.Next() {
		var etd, ppn, total float64
		var itemId, sequenceItem, detail, qty, companyid, customerid, fkid, orderGrade, inputBy, orderStatus, doid int
		var poDate, itemName, size, price, unit, customerName, estDate, customerPoNumber, soNumber, qore, lin, roll, material, volume, note, ukBahanBaku, qtyBahanBaku, sources, merk, woType, isi, companyName, ongkir, sjId, userName, orderGradeStr, porporasi string

		if err := rows.Scan(&itemId, &sequenceItem, &detail, &itemName, &size, &price, &qty, &unit, &companyid, &customerid, &fkid, &customerName, &poDate, &orderGrade, &customerPoNumber, &inputBy, &soNumber, &qore, &lin, &roll, &material, &volume, &note, &porporasi, &ukBahanBaku, &qtyBahanBaku, &sources, &merk, &woType, &ppn, &orderStatus, &etd, &isi, &ongkir, &doid, &sjId, &companyName); err != nil {
			return Response{}, err
		}

		poDateParse, err := time.Parse(config.App.DateFormat_Global, poDate)
		if err != nil {
			return Response{}, err
		}

		poDate = string(poDateParse.Format(config.App.DateFormat_Frontend))
		poDateFuture := poDateParse.AddDate(0, 0, 16)
		estDate = string(poDateFuture.Format(config.App.DateFormat_Frontend))

		queryUserName := fmt.Sprintf(`SELECT name FROM user WHERE id = %d`, inputBy)
		if err = sql.Connection.QueryRow(queryUserName).Scan(&userName); err != nil {
			return Response{}, err
		}

		if orderGrade > 0 {
			orderGradeStr = `Spesial`
		} else {
			orderGradeStr = `Reguler`
		}

		if porporasi == "1" {
			porporasi = `YA`
		} else {
			porporasi = `TIDAK`
		}

		// Parsing dan filter sources
		ParseSources := strings.Split(sources, "|")
		if ParseSources[0] == "2" {
			sourcesDateParse, err := time.Parse(config.App.DateFormat_Global, ParseSources[2])
			if err != nil {
				return Response{}, err
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

		// Parse nomor So
		So := strings.Split(soNumber, "/")

		// Filtering output cost to except dupe
		if utils.InArray(arrCost, fmt.Sprintf(`%d`, doid)) {
			ongkir = `0`
		}

		datatables = append(datatables, DataTables{
			Itemid:       itemId,
			SequenceItem: sequenceItem,
			Detail:       isi,
			Item:         itemName,
			Size:         size,
			Price:        price,
			Qty:          qty,
			Unit:         unit,
			CompanyId:    companyid,
			CustomerId:   customerid,
			FkId:         fkid,
			CustomerName: customerName,
			PoDate:       poDate,
			Estimated:    estDate,
			OrderGrade:   orderGradeStr,
			NoPoCustomer: customerPoNumber,
			InputBy:      inputBy,
			SoNumber:     fmt.Sprintf(`%s/%s%s`, So[0], So[1], So[2]),
			Qore:         qore,
			Lin:          lin,
			Roll:         roll,
			Material:     material,
			Volume:       volume,
			Note:         note,
			Porporasi:    porporasi,
			UkBahanBaku:  ukBahanBaku,
			QtyBahanBaku: qtyBahanBaku,
			Sources:      sources,
			Merk:         merk,
			WoType:       woType,
			OrderStatus:  orderStatus,
			Isi:          isi,
			Ongkir:       ongkir,
			SjId:         sjId,
			CompanyName:  companyName,
			InputName:    userName,
			Etd:          fmt.Sprintf("%.2f", etd),
			Ppn:          fmt.Sprintf(`%.2f`, ppn),
			Total:        fmt.Sprintf(`%.2f`, total),
		})

		arrCost = append(arrCost, fmt.Sprintf(`%d`, doid))
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}
	response.Data = datatables

	return response, nil
}

func Create(Sessionid string, BodyReq []byte) ([]Preorder, error) {
	var preorder Preorder
	var fkId sql.NullInt64
	var soNumber, noWso string
	var id_fk, total int64

	err := json.Unmarshal([]byte(BodyReq), &preorder)
	if err != nil {
		return nil, err
	}

	if preorder.CompanyId < 1 || preorder.CustomerId < 1 {
		return nil, errors.New("pelanggan ID tidak valid")
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	queryFkid := `SELECT id_fk FROM preorder_customer ORDER BY id DESC LIMIT 1`
	if err = sql.Connection.QueryRow(queryFkid).Scan(&fkId); err != nil {
		id_fk = 0 //set default data null
	}

	querySoNumber := `SELECT no_so FROM workorder_item ORDER BY id DESC LIMIT 1`
	if err = sql.Connection.QueryRow(querySoNumber).Scan(&soNumber); err != nil {
		soNumber = `WSO/1807/001` //set default data null
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	soNumberQueue := strings.Split(soNumber, "/")
	now := time.Now()
	c_time := now.Format("0601")

	id_fk = fkId.Int64 + 1

	queryPoCustomer := fmt.Sprintf(`INSERT INTO preorder_customer (id_fk, id_company, id_customer, customer, order_grade, po_date, po_customer, input_by) VALUES (%d, %d, %d, '%s', %d, '%s', '%s', %s)`, id_fk, preorder.CompanyId, preorder.CustomerId, preorder.CustomerName, preorder.OrderGrade, preorder.PoDate, preorder.NoPoCustomer, Sessionid)
	_, errs := tx.Exec(queryPoCustomer)
	if errs != nil {
		tx.Rollback()
		return nil, errs
	}

	queryPoItem := `INSERT INTO preorder_item (id_fk, item_to, detail, item, size, price, qty, unit, input_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmtPoItem, err := tx.Prepare(queryPoItem)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtPoItem.Close()

	queryWoItem := `INSERT INTO workorder_item (id_fk, item_to, detail, no_so, item, size, unit, qore, lin, roll, ingredient, qty, volume, total, annotation, porporasi, uk_bahan_baku, qty_bahan_baku, sources, merk, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmtWoItem, err := tx.Prepare(queryWoItem)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtWoItem.Close()

	queryStatus := `INSERT INTO status (id_fk, item_to, order_status, hidden) VALUES (?, ?, ?, ?)`
	stmtStatus, err := tx.Prepare(queryStatus)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtStatus.Close()

	for index, item := range preorder.Items {
		// validate price is float
		_, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("price format not allowed")
		}

		////////// REDEFIND VARIABLE ////////////////
		SequenceItem := index + 1

		if fkId.Int64 == 0 || c_time > soNumberQueue[1] {
			noWso = fmt.Sprintf("WSO/%s/00%d", c_time, SequenceItem)
		} else {
			lastSoNumber, err := strconv.Atoi(soNumberQueue[2])
			if err != nil {
				return nil, err
			}
			noWso = fmt.Sprintf("WSO/%s/%03d", c_time, lastSoNumber+SequenceItem)
		}

		if item.Unit == "PCS" {
			total = item.Qty / item.Volume
		} else {
			total = item.Qty * item.Volume
		}

		if item.Sources == "2" {
			item.Sources = fmt.Sprintf("%s|%s|%s", item.Sources, strings.ReplaceAll(item.Etc1, "|", "-"), item.Etc2)
		} else if item.Sources == "3" {
			item.Sources = fmt.Sprintf("%s|%s", item.Sources, item.Etc1)
		}

		////////////////////////////////////

		_, err1 := stmtPoItem.Exec(id_fk, SequenceItem, item.Detail, item.Item, item.Size, item.Price, item.Qty, item.Unit, item.InputBy)
		if err1 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err1] %s", err1)
		}

		_, err2 := stmtWoItem.Exec(id_fk, SequenceItem, item.Detail, noWso, item.Item, item.Size, item.Unit, item.Qore, item.Lin, item.Roll, item.Material, item.Qty, item.Volume, total, item.Note, item.Porporasi, item.UkBahanBaku, item.QtyBahanBaku, item.Sources, item.Merk, item.WoType)
		if err2 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err2)
		}

		_, err3 := stmtStatus.Exec(id_fk, SequenceItem, 0, 0)
		if err3 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err3] %s", err3)
		}
	}

	queryPoPrice := fmt.Sprintf(`INSERT INTO preorder_price (id_fk, ppn) VALUES (%d, %d)`, id_fk, preorder.Ppn)
	_, err4 := tx.Exec(queryPoPrice)
	if err4 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err4] %s", err4)
	}

	queryWoCustomer := fmt.Sprintf(`INSERT INTO workorder_customer (id_fk, po_date, spk_date, duration, po_customer, customer, input_by) VALUES (%d, '%s', '%s', '%s', '%s', '%s', %s)`, id_fk, preorder.PoDate, "0000-00-00", "0000-00-00", preorder.NoPoCustomer, preorder.CustomerName, Sessionid)
	_, err5 := tx.Exec(queryWoCustomer)
	if err5 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err5] %s", err5)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`SO Created`,
		fmt.Sprintf(`Customer: %s - PO No: %s - Date: %s`, preorder.CustomerName, preorder.NoPoCustomer, preorder.PoDate),
		Sessionid,
	)

	return []Preorder{}, nil
}

func AddItem(Sessionid string, BodyReq []byte) ([]Preorder, error) {
	var preorder Preorder
	var soNumber, noWso string
	var total, last_sequence_item int64

	err := json.Unmarshal([]byte(BodyReq), &preorder)
	if err != nil {
		return nil, err
	}

	Fkid, _ := strconv.Atoi(preorder.Fkid)
	Id, _ := strconv.Atoi(preorder.Id)
	if Fkid < 1 || Id < 1 {
		return nil, errors.New("PO ID tidak valid")
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	querySoNumber := `SELECT no_so FROM workorder_item ORDER BY id DESC LIMIT 1`
	if err = sql.Connection.QueryRow(querySoNumber).Scan(&soNumber); err != nil {
		soNumber = `WSO/1807/001` //set default data null
	}

	querySequenceItem := fmt.Sprintf(`SELECT COUNT(id) FROM preorder_item WHERE id_fk = %s`, preorder.Fkid)
	if err = sql.Connection.QueryRow(querySequenceItem).Scan(&last_sequence_item); err != nil {
		soNumber = `WSO/1807/001` //set default data null
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	c_time := now.Format("0601")

	queryPoItem := `INSERT INTO preorder_item (id_fk, item_to, detail, item, size, price, qty, unit, input_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmtPoItem, err := tx.Prepare(queryPoItem)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtPoItem.Close()

	queryWoItem := `INSERT INTO workorder_item (id_fk, item_to, detail, no_so, item, size, unit, qore, lin, roll, ingredient, qty, volume, total, annotation, porporasi, uk_bahan_baku, qty_bahan_baku, sources, merk, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmtWoItem, err := tx.Prepare(queryWoItem)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtWoItem.Close()

	queryStatus := `INSERT INTO status (id_fk, item_to, order_status, hidden) VALUES (?, ?, ?, ?)`
	stmtStatus, err := tx.Prepare(queryStatus)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmtStatus.Close()

	fmt.Println(len(preorder.Items))

	log := []map[string]string{}
	for index, item := range preorder.Items {
		// validate price is float
		_, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("price format not allowed")
		}

		////////// REDEFIND VARIABLE ////////////////
		SequenceItem := int(last_sequence_item) + (index + 1)

		noWso = fmt.Sprintf("WSO/%s/%03d", c_time, SequenceItem)

		if item.Unit == "PCS" {
			total = item.Qty / item.Volume
		} else {
			total = item.Qty * item.Volume
		}

		if item.Sources == "2" {
			item.Sources = fmt.Sprintf("%s|%s|%s", item.Sources, strings.ReplaceAll(item.Etc1, "|", "-"), item.Etc2)
		} else if item.Sources == "3" {
			item.Sources = fmt.Sprintf("%s|%s", item.Sources, item.Etc1)
		}

		////////////////////////////////////

		_, err1 := stmtPoItem.Exec(preorder.Fkid, SequenceItem, item.Detail, item.Item, item.Size, item.Price, item.Qty, item.Unit, item.InputBy)
		if err1 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err1] %s", err1)
		}

		_, err2 := stmtWoItem.Exec(preorder.Fkid, SequenceItem, item.Detail, noWso, item.Item, item.Size, item.Unit, item.Qore, item.Lin, item.Roll, item.Material, item.Qty, item.Volume, total, item.Note, item.Porporasi, item.UkBahanBaku, item.QtyBahanBaku, item.Sources, item.Merk, item.WoType)
		if err2 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err2)
		}

		_, err3 := stmtStatus.Exec(preorder.Fkid, SequenceItem, 0, 0)
		if err3 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err3] %s", err3)
		}

		log = append(log, map[string]string{
			"sequence_item": fmt.Sprintf(`%d`, SequenceItem),
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	datalog, _ := json.Marshal(log)
	utils.Capture(
		`SO Created [Add Item]`,
		fmt.Sprintf(`FkId: %s - PoId: %s - data: %s`, preorder.Fkid, preorder.Id, datalog),
		Sessionid,
	)

	return []Preorder{}, nil
}

func GetCustomer(Id int) ([]Preorder, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT a.id AS so_cusid, a.id_company as companyid, a.id_customer AS customerid, a.id_fk, a.customer, a.po_date, a.po_customer, a.order_grade, a.input_by, b.id AS so_priceid, b.ppn, c.company FROM preorder_customer AS a LEFT JOIN preorder_price AS b ON a.id_fk = b.id_fk LEFT JOIN company AS c ON a.id_company = c.id WHERE a.id = %d LIMIT 1", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	preorder := []Preorder{}
	for rows.Next() {
		var so_cusid, so_priceid, companyid, customerid, fkid, ordergrade, inputby, ppn int
		var customername, ponumber, podate, companyname string

		if err := rows.Scan(&so_cusid, &companyid, &customerid, &fkid, &customername, &podate, &ponumber, &ordergrade, &inputby, &so_priceid, &ppn, &companyname); err != nil {
			return nil, err
		}

		preorder = append(preorder, Preorder{
			CompanyName:  companyname,
			CompanyId:    companyid,
			CustomerName: customername,
			CustomerId:   customerid,
			OrderGrade:   ordergrade,
			PoDate:       podate,
			NoPoCustomer: ponumber,
			Ppn:          ppn,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(preorder) < 1 {
		return nil, errors.New("invalid ID")
	}

	return preorder, nil
}

func GetItem(Id int) ([]PreorderItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT a.id as po_itemid, a.id_fk, a.item_to, a.detail, a.item, a.size, a.price, a.qty, a.unit, a.input_by, b.id AS wo_itemid, b.qore, b.lin, b.roll, b.ingredient, b.volume, b.annotation, b.porporasi, b.uk_bahan_baku, b.qty_bahan_baku, b.sources, b.detail, b.type, b.merk FROM preorder_item AS a LEFT JOIN workorder_item AS b ON a.id_fk = b.id_fk AND a.item_to = b.item_to WHERE a.id = %d LIMIT 1", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	preorderitem := []PreorderItem{}
	for rows.Next() {
		var poid, woid, fkid, sequenceitem, porporasi, detail, inputby int
		var volume, qty int64
		var item, size, ukbahanbaku, qore, lin, qtybahanbaku, roll, material, unit, note, price, sources, merk, itemtype string

		if err := rows.Scan(&poid, &fkid, &sequenceitem, &detail, &item, &size, &price, &qty, &unit, &inputby, &woid, &qore, &lin, &roll, &material, &volume, &note, &porporasi, &ukbahanbaku, &qtybahanbaku, &sources, &detail, &itemtype, &merk); err != nil {
			return nil, err
		}

		preorderitem = append(preorderitem, PreorderItem{
			Poitemid:     poid,
			Woitemid:     woid,
			Item:         item,
			Size:         size,
			UkBahanBaku:  ukbahanbaku,
			Qore:         qore,
			Lin:          lin,
			QtyBahanBaku: qtybahanbaku,
			Roll:         roll,
			Material:     material,
			Unit:         unit,
			Volume:       volume,
			Note:         note,
			Price:        utils.StrReplaceAll(price, ".", ","),
			Qty:          qty,
			Sources:      sources,
			Porporasi:    porporasi,
			Detail:       detail,
			Merk:         merk,
			WoType:       itemtype,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(preorderitem) < 1 {
		return nil, errors.New("invalid ID")
	}

	return preorderitem, nil
}

func UpdateCustomer(Sessionid string, Id int, companyid int, customerid int, customername string, ordergrade int, podate string, ponumcustomer string, ppn int) ([]Preorder, error) {
	var id_fk int
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id_fk FROM preorder_customer WHERE id = '%d' LIMIT 1", Id)
	if err = sql.Connection.QueryRow(query_id).Scan(&id_fk); err != nil {
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

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	queryUpdate_PoCustomer := fmt.Sprintf(`UPDATE preorder_customer SET customer = '%s', id_company = %d, id_customer = %d, order_grade = %d, po_date = '%s', po_customer = '%s' WHERE id = %d`, customername, companyid, customerid, ordergrade, podate, ponumcustomer, Id)

	queryUpdate_WoCustomer := fmt.Sprintf(`UPDATE workorder_customer SET po_date = '%s', po_customer = '%s', customer = '%s' WHERE id_fk = %d`, podate, ponumcustomer, customername, id_fk)

	queryUpdate_PoPrice := fmt.Sprintf(`UPDATE preorder_price SET ppn = %d WHERE id_fk = %d`, ppn, id_fk)

	_, err1 := tx.Exec(queryUpdate_PoCustomer)
	if err1 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err1] %s", err1)
	}

	_, err2 := tx.Exec(queryUpdate_WoCustomer)
	if err2 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err2] %s", err2)
	}

	_, err3 := tx.Exec(queryUpdate_PoPrice)
	if err3 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err3] %s", err3)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`SO Updated [customer]`,
		fmt.Sprintf(`PO Id: %d - PO No: %s - CustomerId: %d - CompanyId: %d - PO Date: %s - PPN: %d`, Id, ponumcustomer, customerid, companyid, podate, ppn),
		Sessionid,
	)

	return []Preorder{}, nil
}

func UpdateItem(Sessionid string, poitemid int, woitemid int, item string, size string, ukbahanbaku string, qore string, lin string, qtybahanbaku string, roll string, material string, unit string, volume int64, note string, price string, qty int64, sources string, porporasi int, detail int, merk string, wotype string, etc1 string, etc2 string) ([]PreorderItem, error) {
	var id_fk, item_to int
	var total, total_send_qty int64

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id_fk FROM preorder_item WHERE id = %d LIMIT 1", poitemid)
	if err = sql.Connection.QueryRow(query_id).Scan(&id_fk); err != nil {
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

	if unit == "PCS" {
		total = qty / volume
	} else {
		total = qty * volume
	}

	if sources == "2" {
		sources = fmt.Sprintf("%s|%s|%s", sources, strings.ReplaceAll(etc1, "|", "-"), etc2)
	} else if sources == "3" {
		sources = fmt.Sprintf("%s|%s", sources, etc1)
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	queryUpdate_PoItem := fmt.Sprintf(`UPDATE preorder_item SET item = '%s', size = '%s', qty = %d, unit = '%s', price = '%s' WHERE id = %d`, item, size, qty, unit, price, poitemid)
	_, err1 := tx.Exec(queryUpdate_PoItem)
	if err1 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err1] %s", err1)
	}

	// validasi sources
	if sources == "2" {
		sources = fmt.Sprintf("%s|%s|%s", sources, strings.ReplaceAll(etc1, "|", "-"), etc2)
	} else if sources == "3" {
		sources = fmt.Sprintf("%s|%s", sources, etc1)
	}

	queryUpdate_WoItem := fmt.Sprintf(`UPDATE workorder_item SET item = '%s', detail = %d, merk = '%s', type = '%s', size = '%s', unit = '%s', qty = %d, volume = %d, total = %d, uk_bahan_baku = '%s', qty_bahan_baku = '%s', qore = '%s', lin = '%s', roll = '%s', ingredient = '%s', annotation = '%s', porporasi = %d, sources = '%s' WHERE id = %d`, item, detail, merk, wotype, size, unit, qty, volume, total, ukbahanbaku, qtybahanbaku, qore, lin, roll, material, note, porporasi, sources, woitemid)
	_, err2 := tx.Exec(queryUpdate_WoItem)
	if err2 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err2] %s", err2)
	}

	// validasi $qty jika lebih dari jumlah pengiriman status berubah 6, sedangkan kurang
	// atau cukup akan menjadi 7
	checkTotalQty := fmt.Sprintf(`SELECT a.id_fk, a.item_to, CASE WHEN sum(b.send_qty) > 0 THEN sum(b.send_qty) ELSE 0 END AS total_send_qty FROM preorder_item AS a LEFT JOIN delivery_orders_item AS b ON b.id_fk = a.id_fk AND b.item_to = a.item_to WHERE a.id = '%d'`, poitemid)

	if err3 := sql.Connection.QueryRow(checkTotalQty).Scan(&id_fk, &item_to, &total_send_qty); err3 != nil {
		return nil, fmt.Errorf("[err3] %s", err3)
	}

	queryStatus := `UPDATE status SET order_status =? where id_fk =? AND item_to =?`
	stmtStatus, err := tx.Prepare(queryStatus)
	if err != nil {
		return nil, err
	}
	defer stmtStatus.Close()

	if total_send_qty < 1 {
		if _, err4 := stmtStatus.Exec(0, id_fk, item_to); err4 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err4] %s", err4)
		}
	} else if qty > total_send_qty {
		if _, err5 := stmtStatus.Exec(2, id_fk, item_to); err5 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err5] %s", err5)
		}

	} else {
		if _, err6 := stmtStatus.Exec(1, id_fk, item_to); err6 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err6] %s", err6)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`SO Updated [item]`,
		fmt.Sprintf(`PoitemId: %d - WoitemId: %d - Item: %s - Size: %s - Qty: %d - Unit: %s - Price: %s`, poitemid, woitemid, item, size, qty, unit, price),
		Sessionid,
	)

	return []PreorderItem{}, nil
}

func GetShippingCost(Id int) ([]PreorderShippingCost, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT a.id, a.courier, a.no_tracking, a.cost, a.ekspedisi, a.uom, a.jml, b.no_delivery, b.send_qty FROM delivery_orders_customer AS a LEFT JOIN (SELECT id_fk, no_delivery, send_qty FROM delivery_orders_item WHERE id_fk = %d GROUP BY no_delivery) AS b ON a.id_fk = b.id_fk WHERE a.id_fk = %d GROUP BY a.id ORDER BY b.no_delivery ASC;`, Id, Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	preordershipcost := []PreorderShippingCost{}
	for rows.Next() {
		var id, send_qty int
		var cost, courier, no_tracking, ekspedisi, uom, jml, no_delivery, costStr string

		if err := rows.Scan(&id, &courier, &no_tracking, &cost, &ekspedisi, &uom, &jml, &no_delivery, &send_qty); err != nil {
			return nil, err
		}

		if no_delivery != "" {
			no_delivery = fmt.Sprintf(`SJ: %s`, no_delivery)
		}

		if courier != "" {
			courier = fmt.Sprintf(` - Kurir: %s`, courier)
		}

		if no_tracking != "" {
			no_tracking = fmt.Sprintf(`- No Tracking: %s`, no_tracking)
		}

		if cost != "" {
			costStr = fmt.Sprintf(` - Ongkir: %s`, cost)
		} else {
			costStr = fmt.Sprintf(` - Ongkir: %s`, `0`)
		}

		if send_qty > 0 {
			preordershipcost = append(preordershipcost, PreorderShippingCost{
				Detail:    fmt.Sprintf(`%s%s%s%s`, no_delivery, courier, no_tracking, costStr),
				Cost:      cost,
				Id:        id,
				Ekspedisi: ekspedisi,
				Uom:       uom,
				Jml:       jml,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return preordershipcost, nil
}

func UpdateShipCost(Sessionid string, Id int, cost string, ekspedisi string, uom string, jml string) ([]PreorderShippingCost, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query2 := fmt.Sprintf(`UPDATE delivery_orders_customer SET cost = '%s', ekspedisi = '%s', uom = '%s', jml = '%s' WHERE id = %d`, strings.ReplaceAll(cost, `.`, ``), ekspedisi, uom, jml, Id)
	if _, err := sql.Connection.Query(query2); err != nil {
		return nil, fmt.Errorf("[err] %s", err)
	}

	// Log capture
	utils.Capture(
		`PO Updated [ongkir]`,
		fmt.Sprintf(`DoId: %d - Cost: %s - Eksedisi: %s - Uom: %s - Jumlah: %s`, Id, cost, ekspedisi, uom, jml),
		Sessionid,
	)

	return []PreorderShippingCost{}, nil
}

func Delete(Sessionid string, Id int) ([]PreorderItem, error) {
	var id_fk, item_to, jml_item_dlm_wo int
	var customername, ponumcustomer, podate string

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ambil id_fk dan item_to dari tabel PO_item
	query_id := fmt.Sprintf("SELECT id_fk, item_to FROM preorder_item WHERE id = %d LIMIT 1", Id)
	if err = sql.Connection.QueryRow(query_id).Scan(&id_fk, &item_to); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, fmt.Errorf("invalid ID")

		} else {
			return nil, err
		}
	}

	// Ambil data customer dari tabel PO_customer berdasarkan id_fk diatas (untuk logging)
	query_cus := fmt.Sprintf("SELECT customer, po_customer, po_date FROM preorder_customer WHERE id_fk = %d LIMIT 1", id_fk)
	if err = sql.Connection.QueryRow(query_cus).Scan(&customername, &ponumcustomer, &podate); err != nil {
		return nil, err
	}

	// Mencari total salesorder dalam workorder_item
	query_totalwoitem := fmt.Sprintf("SELECT count(id) AS jml_item_dlm_wo FROM workorder_item WHERE id_fk = %d GROUP BY id_fk", id_fk)
	if err = sql.Connection.QueryRow(query_totalwoitem).Scan(&jml_item_dlm_wo); err != nil {
		return nil, err
	}

	// mencari item dalam delivery orders item
	query_totaldoitem := fmt.Sprintf(`SELECT id_fk, item_to, count(no_delivery) AS jml_item_dlm_sj FROM delivery_orders_item WHERE id_fk = %d GROUP BY no_delivery`, id_fk)
	rows, err := sql.Connection.Query(query_totaldoitem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPointers := make([]interface{}, len(columns))
		for i := range values {
			valuesPointers[i] = &values[i]
		}

		if err := rows.Scan(valuesPointers...); err != nil {
			return nil, err
		}

		results = append(results, values)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	// Menghapus data salesorder, workorder, dan status
	/*if jml_item_dlm_wo < 2 {
		del1 := fmt.Sprintf(`DELETE FROM preorder_customer WHERE id_fk = %d`, id_fk)
		del2 := fmt.Sprintf(`DELETE FROM preorder_item WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)
		del3 := fmt.Sprintf(`DELETE FROM preorder_price WHERE id_fk = %d`, id_fk)
		del4 := fmt.Sprintf(`DELETE FROM workorder_customer WHERE id_fk = %d`, id_fk)
		del5 := fmt.Sprintf(`DELETE FROM workorder_item WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)
		del6 := fmt.Sprintf(`DELETE FROM status WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)

		if _, err1 := tx.Exec(del1); err1 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err1] %s", err1)
		}

		if _, err2 := tx.Exec(del2); err2 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err2)
		}

		if _, err3 := tx.Exec(del3); err3 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err3] %s", err3)
		}

		if _, err4 := tx.Exec(del4); err4 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err4] %s", err4)
		}

		if _, err5 := tx.Exec(del5); err5 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err5] %s", err5)
		}

		if _, err6 := tx.Exec(del6); err6 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err6] %s", err6)
		}
	} else {
		del2 := fmt.Sprintf(`DELETE FROM preorder_item WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)
		del5 := fmt.Sprintf(`DELETE FROM workorder_item WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)
		del6 := fmt.Sprintf(`DELETE FROM status WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)

		if _, err2 := tx.Exec(del2); err2 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err2] %s", err2)
		}

		if _, err5 := tx.Exec(del5); err5 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err5] %s", err5)
		}

		if _, err6 := tx.Exec(del6); err6 != nil {
			tx.Rollback()
			return nil, fmt.Errorf("[err6] %s", err6)
		}
	}*/

	// Menghapus data DO customer, DO item dan invoice berdasarkan kondisi
	for _, row := range results {
		var id_fk, item_to, jml_item_dlm_sj int

		if row[0] != nil {
			id_fk = int(row[0].(int64))
		}

		if row[1] != nil {
			item_to = int(row[1].(int64))
		}

		if row[2] != nil {
			jml_item_dlm_sj = int(row[2].(int64))
		}

		//menghapus DO customer, DO item dan invoice berdasarkan kondisi
		if jml_item_dlm_sj < 2 {
			del7 := fmt.Sprintf(`DELETE FROM delivery_orders_item WHERE id_fk = %d AND item_to = %d`, id_fk, item_to)
			del8 := fmt.Sprintf(`DELETE FROM delivery_orders_customer WHERE id_fk = %d`, id_fk)
			del9 := fmt.Sprintf(`DELETE FROM invoice WHERE id_fk = %d`, id_fk)

			if _, err7 := tx.Exec(del7); err7 != nil {
				tx.Rollback()
				return nil, fmt.Errorf("[err7] %s", err7)
			}

			if _, err8 := tx.Exec(del8); err8 != nil {
				tx.Rollback()
				return nil, fmt.Errorf("[err8] %s", err8)
			}

			if _, err9 := tx.Exec(del9); err9 != nil {
				tx.Rollback()
				return nil, fmt.Errorf("[err9] %s", err9)
			}
		} else {
			// ???????????????????????????????????????????????????
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`SO Delete`,
		fmt.Sprintf(`SoitemId: %d - Fkid: %d - Sequence Item: %d`, Id, id_fk, item_to),
		Sessionid,
	)

	return []PreorderItem{}, nil
}

func SuggestType(ctx *gin.Context) ([]SuggestionsType, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, isi FROM setting WHERE ket ='SO_ITEM'`
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	suggestType := []SuggestionsType{}
	for rows.Next() {
		var id, isi string

		if err := rows.Scan(&id, &isi); err != nil {
			return nil, err
		}

		suggestType = append(suggestType, SuggestionsType{
			Id:   id,
			Item: isi,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suggestType, nil
}

func SuggestCustomer(ctx *gin.Context) ([]SuggestionsCustomer, error) {
	var prevValue, customer, nopo string
	Keyword := ctx.DefaultQuery("keyword", "")

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT a.id AS id_customer, a.nama, CASE WHEN b.po_customer IS NOT NULL THEN b.po_customer ELSE '' END AS po_customer, CASE WHEN b.id IS NOT NULL THEN b.id ELSE '' END AS id_po, CASE WHEN b.id_fk > 0 THEN b.id_fk ELSE 0 END AS id_fk, CASE WHEN c.item IS NOT NULL THEN GROUP_CONCAT(c.item SEPARATOR ' - ') ELSE '' END AS item FROM customer AS a LEFT JOIN preorder_customer AS b ON a.id = b.id_customer LEFT JOIN preorder_item AS c ON c.id_fk = b.id_fk WHERE a.nama LIKE '%%%s%%' GROUP BY b.id_fk ORDER BY b.id DESC`, Keyword)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	datas := []SuggestionsCustomer{}
	for rows.Next() {
		var customerid, customername, poid, nopocustomer, fkid, item string

		if err := rows.Scan(&customerid, &customername, &nopocustomer, &poid, &fkid, &item); err != nil {
			return nil, err
		}

		datas = append(datas, SuggestionsCustomer{
			Customerid:   customerid,
			CustomerName: customername,
			NoPoCustomer: nopocustomer,
			Poid:         poid,
			Item:         item,
			Fkid:         fkid,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	suggestions := []SuggestionsCustomer{}
	if len(datas) > 0 {
		prevValue = ``
		for _, data := range datas {
			customer = strings.ReplaceAll(data.CustomerName, ` `, `_`)
			nopo = strings.ReplaceAll(data.NoPoCustomer, ` `, `_`)

			if nopo == "" {
				suggestions = append(suggestions, SuggestionsCustomer{
					Customerid: data.Customerid,
					Poid:       ``,
					Label:      `Buat sales baru`,
					Value:      data.CustomerName,
					Category:   data.CustomerName,
				})

			} else {

				if prevValue != customer {
					suggestions = append(suggestions, SuggestionsCustomer{
						Customerid: data.Customerid,
						Poid:       ``,
						Label:      `Buat sales baru`,
						Value:      data.CustomerName,
						Category:   data.CustomerName,
					})

				} else {
					suggestions = append(suggestions, SuggestionsCustomer{
						Customerid: data.Customerid,
						Poid:       data.Poid,
						Fkid:       data.Fkid,
						Label:      fmt.Sprintf(`No PO: %s - %s`, data.NoPoCustomer, data.Item),
						Value:      data.CustomerName,
						Category:   data.CustomerName,
					})
				}
			}

			prevValue = customer
		}

	} else {
		suggestions = append(suggestions, SuggestionsCustomer{
			Customerid: ``,
			Poid:       ``,
			Label:      `Tidak terdaftar, silakan daftar sebagai customer baru.`,
			Value:      Keyword,
			Category:   ``,
		})
	}

	return suggestions, nil
}

func SuggestItem(ctx *gin.Context) ([]SuggestionsItem, error) {
	Customerid := ctx.DefaultQuery("customerid", "")
	Poid := ctx.DefaultQuery("poid", "")

	customerid, err := strconv.Atoi(Customerid)
	if err != nil {
		return nil, err
	}

	poid, err := strconv.Atoi(Poid)
	if err != nil {
		return nil, err
	}

	suggestionsitem := SuggestionsItem{
		Customerid: customerid,
	}

	if poid > 0 {
		if customerid == 0 {
			return nil, errors.New("invalid ID")
		}

		sql, err := adapters.NewSql()
		if err != nil {
			return nil, err
		}

		query := fmt.Sprintf(`SELECT c.item_to, CASE WHEN c.price > 0 THEN c.price ELSE 0 END AS price, d.item, d.size, d.unit, d.qore, d.lin, d.roll, d.ingredient, d.qty, d.volume, d.annotation, d.porporasi, d.uk_bahan_baku, d.qty_bahan_baku, d.detail, d.merk, d.type FROM customer AS a LEFT JOIN preorder_customer AS b ON a.id = b.id_customer LEFT JOIN preorder_item AS c ON c.id_fk = b.id_fk LEFT JOIN workorder_item AS d ON d.id_fk = c.id_fk AND d.item_to = c.item_to WHERE a.id = %d AND b.id = %d`, customerid, poid)

		rows, err := sql.Connection.Query(query)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		suggestionsitem.PoId = poid

		soitem := []PreorderItem{}
		for rows.Next() {
			var price float64
			var sequence_item, porporasi, detail int
			var qty, volume int64
			var item, size, unit, qore, lin, roll, material, note, uk_bahan_baku, qty_bahan_baku, merk, itemtype string

			if err := rows.Scan(&sequence_item, &price, &item, &size, &unit, &qore, &lin, &roll, &material, &qty, &volume, &note, &porporasi, &uk_bahan_baku, &qty_bahan_baku, &detail, &merk, &itemtype); err != nil {
				return nil, err
			}

			soitem = append(soitem, PreorderItem{
				SequenceItem: sequence_item,
				Price:        fmt.Sprintf(`%.2f`, price),
				Item:         item,
				Size:         size,
				Unit:         unit,
				Qore:         qore,
				Lin:          lin,
				Roll:         roll,
				Material:     material,
				Qty:          qty,
				Volume:       volume,
				Note:         note,
				Porporasi:    porporasi,
				UkBahanBaku:  uk_bahan_baku,
				QtyBahanBaku: qty_bahan_baku,
				Detail:       detail,
				Merk:         merk,
				WoType:       itemtype,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		suggestionsitem.Items = soitem

	} else {
		suggestionsitem.PoId = 0
	}

	return []SuggestionsItem{suggestionsitem}, nil
}

func SuggestAttr(ctx *gin.Context) (interface{}, error) {
	Id := ctx.DefaultQuery("id", "")

	IdInt, err := strconv.Atoi(Id)
	if err != nil {
		return nil, err
	}

	if IdInt < 1 {
		return nil, errors.New("invalid ID")
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	var isi, value string
	query := fmt.Sprintf(`SELECT isi, value FROM setting WHERE id = %d`, IdInt)
	err = sql.Connection.QueryRow(query).Scan(&isi, &value)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	// Unmarshal the JSON string into a map without struct
	var result map[string]string
	if err = json.Unmarshal([]byte(value), &result); err != nil {
		return nil, err
	}

	output := map[string]string{
		"field": result["input"],
	}

	return output, nil
}

func SuggestSO(ctx *gin.Context) (interface{}, error) {
	Id := ctx.DefaultQuery("id", "")

	IdInt, err := strconv.Atoi(Id)
	if err != nil {
		return nil, err
	}

	if IdInt < 1 {
		return nil, errors.New("invalid ID")
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	var value string
	query := fmt.Sprintf(`SELECT value FROM setting WHERE id = %d`, IdInt)
	err = sql.Connection.QueryRow(query).Scan(&value)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	// Unmarshal the JSON string into a map without struct
	var result map[string]string
	if err = json.Unmarshal([]byte(value), &result); err != nil {
		return nil, err
	}

	output := map[string]string{
		"field": result["input"],
	}

	return output, nil
}

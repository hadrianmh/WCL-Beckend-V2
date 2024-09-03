package workorder

import (
	"backend/adapters"
	"backend/config"
	"backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
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

type WorkOrder struct {
	Id           int    `json:"id"`
	SequenceItem int    `json:"item_to"`
	SpkDate      string `json:"spk_date"`
	OrderStatus  int    `json:"order_status"`
	InputBy      int    `json:"input_by"`
}

type DataTables struct {
	WoCusid      int     `json:"wocusid,omitempty"`
	FkId         int     `json:"fkid,omitempty"`
	SpkDate      string  `json:"spk_date"`
	Duration     string  `json:"duration"`
	NoPoCustomer string  `json:"po_customer,omitempty"`
	CustomerId   int     `json:"customerid,omitempty"`
	CustomerName string  `json:"customer"`
	InputBy      int     `json:"input_by"`
	WoItemid     int     `json:"woitemid,omitempty"`
	SequenceItem int     `json:"sequence_item,omitempty"`
	Detail       int     `json:"detail"`
	Item         string  `json:"item"`
	Size         string  `json:"size"`
	Unit         string  `json:"unit"`
	Qore         string  `json:"qore"`
	Lin          string  `json:"lin"`
	Roll         string  `json:"roll"`
	Material     string  `json:"ingredient"`
	Qty          int     `json:"qty,omitempty"`
	Volume       string  `json:"volume,omitempty"`
	Total        int     `json:"total,omitempty"`
	Ttl          float64 `json:"ttl,omitempty"`
	Note         string  `json:"annotation"`
	Porporasi    string  `json:"porporasi"`
	UkBahanBaku  string  `json:"uk_bahan_baku,"`
	QtyBahanBaku string  `json:"qty_bahan_baku,"`
	Sources      string  `json:"sources,omitempty"`
	Merk         string  `json:"merk,omitempty"`
	WoType       string  `json:"type,omitempty"`
	OrderStatus  string  `json:"orderstatus,omitempty"`
	QtyProd      int     `json:"qty_produksi,omitempty"`
	Isi          int     `json:"isi,omitempty"`
	NoSpk        string  `json:"no_spk,omitempty"`
	SatuanUnit   string  `json:"satuanunit,omitempty"`
	DateNow      string  `json:"tgl,omitempty"`
	Ttd          string  `json:"ttd,omitempty"`
	Username     string  `json:"username"`
}

type WorkorderItem struct {
	PoDate       string `json:"po_date"`
	SpkDate      string `json:"spk_date"`
	NoSpk        string `json:"no_spk"`
	NoPoCustomer string `json:"po_customer"`
	CustomerName string `json:"customer"`
	OrderStatus  string `json:"orderstatus"`
}

func Get(ctx *gin.Context) (Response, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return Response{}, err
	}

	var totalrows int
	var query, search, query_datatables, Report string
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

	if SearchValue != "" {
		wso := SearchValue // wso search spesific
		if strings.Contains(strings.ToLower(SearchValue), config.App.PrefixSearch_wso) {
			trimPrefix := strings.TrimPrefix(strings.ToLower(SearchValue), config.App.PrefixSearch_wso)
			if len(trimPrefix) < 4 {
				wso = config.App.PrefixSearch_wso + trimPrefix
			} else {
				wso = config.App.PrefixSearch_wso + trimPrefix[:4] + "/" + trimPrefix[4:]
			}
		}

		search = fmt.Sprintf(`AND (a.po_date LIKE '%%%s%%' OR a.spk_date LIKE '%%%s%%' OR a.duration LIKE '%%%s%%' OR a.po_customer LIKE '%%%s%%' OR a.customer LIKE '%%%s%%' OR b.no_so LIKE '%%%s%%' OR b.item LIKE '%%%s%%' OR b.size LIKE '%%%s%%' OR b.unit LIKE '%%%s%%' OR b.qore LIKE '%%%s%%' OR b.lin LIKE '%%%s%%' OR b.roll LIKE '%%%s%%' OR b.ingredient LIKE '%%%s%%' OR b.qty LIKE '%%%s%%' OR b.volume LIKE '%%%s%%' OR b.total LIKE '%%%s%%' OR b.annotation LIKE '%%%s%%' OR b.uk_bahan_baku LIKE '%%%s%%' OR b.qty_bahan_baku LIKE '%%%s%%' OR b.sources LIKE '%%%s%%' OR b.merk LIKE '%%%s%%' OR b.type LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, wso, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return Response{}, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query_datatables = fmt.Sprintf(`SELECT
	COUNT(b.id)
	FROM workorder_customer AS a
	LEFT JOIN
		workorder_item AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		status AS c ON a.id_fk = c.id_fk AND b.item_to = c.item_to
	WHERE
		a.po_date %s %s ORDER BY a.id DESC`, Report, search)

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
		a.id AS woid,
		a.id_fk,
		a.po_date,
		a.spk_date,
		a.duration,
		a.po_customer,
		a.customer,
		a.input_by,
		b.id AS woitemid,
		b.item_to,
		b.detail,
		b.no_so,
		b.item,
		b.size,
		b.unit,
		b.qore,
		b.lin,
		b.roll,
		b.ingredient,
		b.qty,
		b.volume,
		b.total,
		b.annotation,
		b.porporasi,
		b.uk_bahan_baku,
		b.qty_bahan_baku,
		b.sources,
		b.merk,
		b.type,
		c.order_status,
		a.id AS id_customer,
		e.name
	FROM workorder_customer AS a
	LEFT JOIN
		workorder_item AS b ON a.id_fk = b.id_fk
	LEFT JOIN
		status AS c ON a.id_fk = c.id_fk AND b.item_to = c.item_to
	LEFT JOIN
		user AS e ON a.input_by = e.id
	WHERE
		a.po_date %s %s ORDER BY b.id DESC LIMIT %d OFFSET %d`, Report, search, limit, offset)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	datatables := []DataTables{}
	for rows.Next() {
		var volume float64
		var wocusid, fkid, customerid, inputby, woitemid, sequence_item, detail, qty, total int
		var podate, spk_date, duration, nopocustomer, customername, so_no, itemname, size, unit, qore, lin, roll, material, note, ukBahanBaku, qtyBahanBaku, sources, merk, woType, porporasi, orderStatus, username string

		if err := rows.Scan(&wocusid, &fkid, &podate, &spk_date, &duration, &nopocustomer, &customername, &inputby, &woitemid, &sequence_item, &detail, &so_no, &itemname, &size, &unit, &qore, &lin, &roll, &material, &qty, &volume, &total, &note, &porporasi, &ukBahanBaku, &qtyBahanBaku, &sources, &merk, &woType, &orderStatus, &customerid, &username); err != nil {
			return Response{}, err
		}

		// Parsing nomor so untuk generate nomor spk
		noSpk := strings.Split(so_no, "/")

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
				return Response{}, err
			}

			spk_date = SpkDateParse.Format(config.App.DateFormat_Frontend)
		}

		// Estimate SPK date format
		if duration == "1970-01-01" || duration == "0000-00-00" {
			duration = ``
		} else {
			durationDateParse, err := time.Parse(config.App.DateFormat_Global, duration)
			if err != nil {
				return Response{}, err
			}

			duration = durationDateParse.Format(config.App.DateFormat_Frontend)
		}

		datatables = append(datatables, DataTables{
			WoCusid:      wocusid,
			FkId:         fkid,
			SpkDate:      spk_date,
			Duration:     duration,
			NoPoCustomer: nopocustomer,
			CustomerName: customername,
			InputBy:      inputby,
			WoItemid:     woitemid,
			SequenceItem: sequence_item,
			Detail:       detail,
			Item:         itemname,
			Size:         size,
			Unit:         unit,
			Qore:         qore,
			Lin:          lin,
			Roll:         roll,
			Material:     material,
			Qty:          qty,
			Volume:       fmt.Sprintf(`%.2f`, volume),
			Total:        total,
			Note:         note,
			Porporasi:    porporasi,
			UkBahanBaku:  ukBahanBaku,
			QtyBahanBaku: qtyBahanBaku,
			Sources:      sources,
			Merk:         merk,
			WoType:       woType,
			OrderStatus:  orderStatus,
			CustomerId:   customerid,
			NoSpk:        fmt.Sprintf(`%s/%s%s`, noSpk[0], noSpk[1], noSpk[2]),
			Username:     username,
		})
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

func Create(Id int, sequence_item int, spkdate string, orderstatus int, Sessionid string) ([]WorkOrder, error) {
	var id_fk int
	var duration string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	// ambil id_fk pada tabel workorder_customer
	query_id := fmt.Sprintf("SELECT id_fk FROM workorder_customer WHERE id = '%d' LIMIT 1", Id)
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

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	// generate estimasi spk date + 16 days
	SpkDateParse, err := time.Parse(config.App.DateFormat_Global, spkdate)
	if err != nil {
		return nil, err
	}

	SpkDateFuture := SpkDateParse.AddDate(0, 0, 16)
	duration = string(SpkDateFuture.Format(config.App.DateFormat_Global))

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	// update spk date dan durasi pd tabel WO_customer
	queryUpdate_WoCustomer := fmt.Sprintf(`UPDATE workorder_customer SET spk_date = '%s', duration = '%s', input_by = '%s' WHERE id = %d`, spkdate, duration, Sessionid, Id)

	// update pd tabel status
	queryUpdate_Status := fmt.Sprintf(`UPDATE status SET order_status = '%d' WHERE id_fk = %d AND item_to = %d`, orderstatus, id_fk, sequence_item)

	_, err1 := tx.Exec(queryUpdate_WoCustomer)
	if err1 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err1] %s", err1)
	}

	_, err2 := tx.Exec(queryUpdate_Status)
	if err2 != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[err2] %s", err2)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`WO Process`,
		fmt.Sprintf(`Woid: %d - Sequence Item: %d - SpkDate: %s - Status: %d`, Id, sequence_item, spkdate, orderstatus),
		Sessionid,
	)

	return []WorkOrder{}, nil
}

func GetProcess(Id string, SequenceItem string) ([]WorkorderItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf(`SELECT a.po_date, a.spk_date, a.po_customer, a.customer, b.order_status, c.no_so FROM workorder_customer AS a LEFT JOIN status AS b ON a.id_fk = b.id_fk LEFT JOIN workorder_item AS c ON c.id_fk = a.id_fk WHERE a.id = '%s' AND b.item_to = '%s' AND c.item_to = '%s'`, Id, SequenceItem, SequenceItem)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	woitem := []WorkorderItem{}
	for rows.Next() {
		var podate, spkdate, nopo, customer, orderstatus, so_no string

		if err := rows.Scan(&podate, &spkdate, &nopo, &customer, &orderstatus, &so_no); err != nil {
			return nil, err
		}

		// Parsing nomor so untuk generate nomor spk
		noSpk := strings.Split(so_no, "/")

		woitem = append(woitem, WorkorderItem{
			PoDate:       podate,
			SpkDate:      spkdate,
			NoSpk:        fmt.Sprintf(`%s/%s%s`, noSpk[0], noSpk[1], noSpk[2]),
			NoPoCustomer: nopo,
			CustomerName: customer,
			OrderStatus:  orderstatus,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(woitem) < 1 {
		return nil, errors.New("invalid ID")
	}

	return woitem, nil
}

func Printview(Wocusid int, Sequenceitem int) ([]DataTables, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf(`SELECT a.id AS wocusid, a.id_fk, a.po_date, a.spk_date, a.duration, a.po_customer, a.customer, a.input_by, b.id AS woitemid, b.item_to, b.detail, b.no_so, b.item, b.size, b.unit, b.qore, b.lin, b.roll, b.ingredient, b.qty, b.volume, b.total, b.annotation, b.porporasi, b.uk_bahan_baku, b.qty_bahan_baku, CASE WHEN b.sources > 0 THEN b.sources ELSE 0 END AS sources, b.merk, b.type FROM workorder_customer AS a LEFT JOIN workorder_item AS b ON a.id_fk = b.id_fk LEFT JOIN status AS c ON a.id_fk = c.id_fk AND b.item_to = c.item_to WHERE a.id = %d AND b.item_to = %d`, Wocusid, Sequenceitem)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datatables := []DataTables{}
	for rows.Next() {
		var Total, total float64
		var wocusid, id_fk, input_by, woitemid, sequence_item, detail, qty, Qty, Isi, exSources1, exSources2, volume int
		var po_date, spk_date, duration, nopocustomer, customername, no_so, item, size, unit, qore, lin, roll, material, note, uk_bahan_baku, qty_bahan_baku, sources, merk, itemType, porporasi string

		if err := rows.Scan(&wocusid, &id_fk, &po_date, &spk_date, &duration, &nopocustomer, &customername, &input_by, &woitemid, &sequence_item, &detail, &no_so, &item, &size, &unit, &qore, &lin, &roll, &material, &qty, &volume, &total, &note, &porporasi, &uk_bahan_baku, &qty_bahan_baku, &sources, &merk, &itemType); err != nil {
			return nil, err
		}

		// Parse nomor so
		exNoSo := strings.Split(no_so, "/")

		// Inisiasi sources value
		exSources := strings.Split(sources, "|")
		if len(exSources) > 1 {
			conv1, err := strconv.Atoi(exSources[0])
			if err != nil {
				return nil, err
			}

			if exSources[1] == "" || reflect.TypeOf(exSources[1]).Kind() != reflect.Int {
				exSources1 = conv1
				exSources2 = 0
			} else {
				conv2, err := strconv.Atoi(exSources[1])
				if err != nil {
					return nil, err
				}
				exSources1 = conv1
				exSources2 = conv2
			}

		} else {
			conv1, err := strconv.Atoi(exSources[0])
			if err != nil {
				return nil, err
			}
			exSources1 = conv1
			exSources2 = 0
		}

		// Inisiasi Total, Qty, Isi
		if exSources1 == 3 {
			if unit == "PCS" {
				if exSources2 >= qty {
					Total = 0
					Qty = 0
					Isi = 0
				} else {
					Qty = qty - exSources2
					Isi = volume
					Total = math.Round(float64((Qty/Isi)*10) / 10)
				}

			} else {
				if exSources2 >= qty {
					Total = 0
					Qty = 0
					Isi = 0
				} else {
					Qty = qty - exSources2
					Isi = volume
					Total = float64(Qty * Isi)
				}
			}

		} else {
			Qty = qty
			Isi = volume
			Total = total
		}

		// Inisiasi porporasi value
		if porporasi == "1" {
			porporasi = "Ya"
		} else {
			porporasi = "Tidak"
		}

		datatables = append(datatables, DataTables{
			SpkDate:      spk_date,
			CustomerName: customername,
			NoSpk:        fmt.Sprintf(`%s/%s%s`, exNoSo[0], exNoSo[1], exNoSo[2]),
			Size:         size,
			Unit:         unit,
			Ttl:          Total,
			Qore:         qore,
			Lin:          lin,
			Roll:         roll,
			Material:     material,
			QtyProd:      Qty,
			Isi:          Isi,
			Note:         note,
			NoPoCustomer: nopocustomer,
			UkBahanBaku:  uk_bahan_baku,
			QtyBahanBaku: qty_bahan_baku,
			Porporasi:    porporasi,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return datatables, nil
}

func Printnow(BodyReq []byte) ([]DataTables, error) {
	var datatables DataTables

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, fmt.Errorf("[err3] %s", err)
	}

	err = json.Unmarshal([]byte(BodyReq), &datatables)
	if err != nil {
		return nil, err
	}

	// Formating Workorder date
	SpkDateParse, err := time.Parse(config.App.DateFormat_Global, datatables.SpkDate)
	SpkDate := SpkDateParse.Format(config.App.DateFormat_Print)
	if err != nil {
		return nil, err
	}

	// Generate
	workorder := []DataTables{}
	workorder = append(workorder, DataTables{
		DateNow:      datatables.DateNow,
		SpkDate:      SpkDate,
		CustomerName: datatables.CustomerName,
		NoSpk:        datatables.NoSpk,
		Note:         datatables.Note,
		Size:         datatables.Size,
		UkBahanBaku:  datatables.UkBahanBaku,
		Material:     datatables.Material,
		Roll:         datatables.Roll,
		Qore:         datatables.Qore,
		Lin:          datatables.Lin,
		Porporasi:    datatables.Porporasi,
		QtyBahanBaku: datatables.QtyBahanBaku,
		Ttl:          datatables.Ttl,
		Isi:          datatables.Isi,
		Ttd:          datatables.Ttd,
		NoPoCustomer: datatables.NoPoCustomer,
		QtyProd:      datatables.QtyProd,
	})

	return workorder, nil
}

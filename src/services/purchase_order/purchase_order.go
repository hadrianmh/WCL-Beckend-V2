package purchaseorder

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

type PurchaseOrder struct {
	PoId          int      `json:"po_id,omitempty"`
	Vendorid      int      `json:"vendorid,omitempty"`
	Companyid     int      `json:"companyid,omitempty"`
	PoDate        string   `json:"po_date,omitempty"`
	PoNumber      string   `json:"po_number,omitempty"`
	Note          string   `json:"note,omitempty"`
	Ppn           string   `json:"tax,omitempty"`
	PoType        string   `json:"po_type,omitempty"`
	MonthlyReport string   `json:"monthly_report,omitempty"`
	InputBy       int      `json:"inputby,omitempty"`
	VendorName    string   `json:"vendor,omitempty"`
	VendorAddress string   `json:"vendor_address,omitempty"`
	CompanyName   string   `json:"company,omitempty"`
	UserId        int      `json:"userid,omitempty"`
	UserName      string   `json:"user,omitempty"`
	Items         []PoItem `json:"items,omitempty"`
	Fkid          int      `json:"fkid,omitempty"`
}

type PoItem struct {
	Id           int    `json:"itemid,omitempty"`
	Fkid         int    `json:"fkid,omitempty"`
	SequenceItem int    `json:"sequence_item"`
	Detail       string `json:"detail"`
	Size         string `json:"size"`
	Price1       string `json:"price_1"`
	Price2       string `json:"price_2"`
	Subtotal     string `json:"subtotal,omitempty"`
	Qty          int    `json:"qty"`
	Unit         string `json:"unit"`
	Merk         string `json:"merk"`
	ItemType     string `json:"item_type"`
	Core         string `json:"core"`
	Roll         string `json:"roll"`
	Material     string `json:"material"`
	InputAttr    string `json:"inputattr,omitempty"`
	PrintAttr    string `json:"printattr,omitempty"`
}

type DataTables struct {
	PoId        int    `json:"po_id"`
	PoDate      string `json:"po_date"`
	CompanyName string `json:"company"`
	VendorName  string `json:"vendor"`
	PoNumber    string `json:"po_number"`
	PoType      string `json:"po_type"`
	Detail      string `json:"detail"`
	Size        string `json:"size"`
	Price1      string `json:"price_1"`
	Price2      string `json:"price_2"`
	Qty         int    `json:"qty"`
	Unit        string `json:"unit"`
	Merk        string `json:"merk"`
	ItemType    string `json:"item_type"`
	Core        string `json:"core"`
	Roll        string `json:"roll"`
	Material    string `json:"material"`
	Note        string `json:"note"`
	UserName    string `json:"user"`
	ItemId      int    `json:"itemid"`
	Subtotal    string `json:"subtotal"`
	Tax         string `json:"tax"`
	Total       string `json:"total"`
}

type Print struct {
	VendorName     string   `json:"vendor,omitempty"`
	VendorAddress  string   `json:"vendor_address,omitempty"`
	PoDate         string   `json:"po_date,omitempty"`
	PoNumber       string   `json:"po_number,omitempty"`
	PoType         string   `json:"po_type,omitempty"`
	Note           string   `json:"note,omitempty"`
	Ppn            string   `json:"tax,omitempty"`
	Ttd            string   `json:"ttd,omitempty"`
	PrintDate      string   `json:"print_date,omitempty"`
	InputAttr      string   `json:"inputattr,omitempty"`
	PrintAttr      string   `json:"printattr,omitempty"`
	CompanyName    string   `json:"companyname,omitempty"`
	CompanyAddress string   `json:"companyaddress,omitempty"`
	Companyemail   string   `json:"companyemail,omitempty"`
	CompanyPhone   string   `json:"companyphone,omitempty"`
	CompanyLogo    string   `json:"companylogo,omitempty"`
	Subtotal       string   `json:"subtotal,omitempty"`
	Taxtotal       string   `json:"taxtotal,omitempty"`
	Total          string   `json:"total,omitempty"`
	Items          []PoItem `json:"items,omitempty"`
}

type SuggestionsVendor struct {
	VendorName   string `json:"vendorname,omitempty"`
	NoPoCustomer string `json:"nopocustomer,omitempty"`
	ItemType     string `json:"itemtype,omitempty"`
	Detail       string `json:"detail,omitempty"`
	Isi          string `json:"isi,omitempty"`
	Vendorid     string `json:"vendorid"`
	Poid         string `json:"po_id,omitempty"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Category     string `json:"category"`
}

type SuggestionsItem struct {
	PoId     int      `json:"po_id,omitempty"`
	Vendorid int      `json:"vendorid,omitempty"`
	Type     string   `json:"type,omitempty"`
	Attr     string   `json:"attr,omitempty"`
	Items    []PoItem `json:"items,omitempty"`
}

type SuggestionsType struct {
	Id   string `json:"id,omitempty"`
	Item string `json:"item,omitempty"`
}

type SuggestionsPO struct {
	Fkid         string `json:"fkid"`
	SequenceItem string `json:"sequence_item"`
	NoPoCustomer string `json:"po_number"`
	ItemType     string `json:"item_type"`
	Detail       string `json:"detail"`
	Isi          string `json:"isi"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Category     string `json:"category"`
}

type Attr struct {
	Input string `json:"input,omitempty"`
	Print string `json:"print,omitempty"`
}
type AddItemPurchaseOrder struct {
	Fkid     int      `json:"fkid" binding:"required"`
	PoNumber string   `json:"po_number" binding:"required"`
	Items    []PoItem `json:"items,omitempty"`
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
		search = fmt.Sprintf(`AND (a.po_date LIKE '%%%s%%' OR a.nopo LIKE '%%%s%%' OR a.note LIKE '%%%s%%' OR a.ppn LIKE '%%%s%%' OR b.vendor LIKE '%%%s%%' OR c.detail LIKE '%%%s%%' OR c.size LIKE '%%%s%%' OR c.price_1 LIKE '%%%s%%' OR c.price_2 LIKE '%%%s%%' OR c.qty LIKE '%%%s%%' OR c.unit LIKE '%%%s%%' OR c.merk LIKE '%%%s%%' OR c.type LIKE '%%%s%%' OR c.core LIKE '%%%s%%' OR c.core LIKE '%%%s%%' OR c.gulungan LIKE '%%%s%%' OR c.bahan LIKE '%%%s%%' OR d.name LIKE '%%%s%%' OR e.company LIKE '%%%s%%' OR f.isi LIKE '%%%s%%')`, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue, SearchValue)
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return Response{}, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query_datatables = fmt.Sprintf(`SELECT COUNT(a.id) as totalrows FROM po_customer AS a JOIN vendor AS b ON a.id_vendor = b.id JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON a.id = c.id_fk JOIN user AS d ON a.input_by = d.id JOIN company AS e ON e.id = a.id_company JOIN setting AS f ON f.id = a.type WHERE a.po_date %s %s ORDER BY a.id DESC`, Report, search)
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

	// id, err := strconv.Atoi(idParam)
	// if err != nil || id < 1 {
	query = fmt.Sprintf(`SELECT a.id AS id_po, a.po_date, a.nopo, a.note, a.ppn, a.input_by, CASE WHEN b.vendor IS NOT NULL THEN b.vendor ELSE '' END AS vendor, c.id AS id_po_item, c.detail, c.size, c.price_1, CASE WHEN c.price_2 = '' THEN 0 ELSE price_2 END AS price_2, c.qty, c.unit, c.merk, c.type, c.core, c.gulungan, c.bahan, d.id AS userid, d.name, e.company, f.isi FROM po_customer AS a LEFT JOIN vendor AS b ON a.id_vendor = b.id LEFT JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON a.id = c.id_fk LEFT JOIN user AS d ON a.input_by = d.id LEFT JOIN company AS e ON e.id = a.id_company LEFT JOIN setting AS f ON f.id = a.type WHERE a.po_date %s %s ORDER BY a.id DESC LIMIT %d OFFSET %d`, Report, search, limit, offset)

	// } else {
	// 	query = fmt.Sprintf(`SELECT a.id AS id_po, a.po_date, a.nopo, a.note, a.ppn, a.input_by, b.vendor, c.id AS id_po_item, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.merk, c.type, c.core, c.gulungan, c.bahan, d.id AS userid, d.name, e.company, f.isi FROM po_customer AS a JOIN vendor AS b ON a.id_vendor = b.id JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON a.id = c.id_fk JOIN user AS d ON a.input_by = d.id JOIN company AS e ON e.id = a.id_company JOIN setting AS f ON f.id = a.type WHERE a.id = %d`, id)
	// }

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return Response{}, err
	}

	defer rows.Close()

	// 2. purchaseorders := []PurchaseOrder{}
	// 1. poMap := make(map[int]*PurchaseOrder)
	datatables := []DataTables{}
	for rows.Next() {
		var price1, price2, subtotal, tax, total, qty float64
		var poId, inputBy, itemId, userId, ppn int
		var poDate, poNumber, note, vendorName, detail, size, unit, merk, itemType, core, roll, material, userName, companyName, poType string

		if err := rows.Scan(
			&poId, &poDate, &poNumber, &note, &ppn, &inputBy, &vendorName,
			&itemId, &detail, &size, &price1, &price2, &qty, &unit, &merk,
			&itemType, &core, &roll, &material, &userId, &userName, &companyName, &poType,
		); err != nil {
			return Response{}, err
		}

		if price2 > 0 {
			subtotal = qty * price2
		} else {
			subtotal = qty * price1
		}

		if ppn > 0 {
			tax = subtotal * 11 / 100
		} else {
			tax = 0
		}

		total = subtotal + tax

		datatables = append(datatables, DataTables{
			PoId:        poId,
			PoDate:      poDate,
			CompanyName: companyName,
			VendorName:  vendorName,
			PoNumber:    poNumber,
			PoType:      poType,
			Detail:      detail,
			Size:        size,
			Price1:      fmt.Sprintf(`%.2f`, price1),
			Price2:      fmt.Sprintf(`%.2f`, price2),
			Qty:         int(qty),
			Unit:        unit,
			Merk:        merk,
			ItemType:    itemType,
			Core:        core,
			Roll:        roll,
			Material:    material,
			Note:        note,
			UserName:    userName,
			ItemId:      itemId,
			Subtotal:    fmt.Sprintf(`%.2f`, subtotal),
			Tax:         fmt.Sprintf(`%.2f`, tax),
			Total:       fmt.Sprintf(`%.2f`, total),
		})

		// 2.
		//  purchaseorders = append(purchaseorders, PurchaseOrder{
		// 	PoId:        poId,
		// 	PoDate:      poDate,
		// 	PoNumber:    poNumber,
		// 	Note:        note,
		// 	Ppn:         ppn,
		// 	InputBy:     inputBy,
		// 	VendorName:  vendorName,
		// 	CompanyName: companyName,
		// 	UserId:      userId,
		// 	UserName:    userName,
		// 	PoType:      poType,
		// 	Items: []PoItem{
		// 		{
		// 			Id:       itemId,
		// 			Detail:   detail,
		// 			Size:     size,
		// 			Price1:   price1,
		// 			Price2:   price2,
		// 			Qty:      qty,
		// 			Unit:     unit,
		// 			Merk:     merk,
		// 			ItemType: itemType,
		// 			Core:     core,
		// 			Roll:     roll,
		// 			Material: material,
		// 		},
		// 	},
		// })

		// 1.
		// if _, exists := poMap[poId]; !exists {
		// 	poMap[poId] = &PurchaseOrder{
		// 		PoId:        poId,
		// 		PoDate:      poDate,
		// 		PoNumber:    poNumber,
		// 		Note:        note,
		// 		Ppn:         ppn,
		// 		InputBy:     inputBy,
		// 		VendorName:  vendorName,
		// 		CompanyName: companyName,
		// 		UserId:      userId,
		// 		UserName:    userName,
		// 		PoType:      poType,
		// 		Items:       []PoItem{},
		// 	}
		// }

		// if itemId != 0 {
		// 	poMap[poId].Items = append(poMap[poId].Items, PoItem{
		// 		Id:       itemId,
		// 		Detail:   detail,
		// 		Size:     size,
		// 		Price1:   price1,
		// 		Price2:   price2,
		// 		Qty:      qty,
		// 		Unit:     unit,
		// 		Merk:     merk,
		// 		ItemType: itemType,
		// 		Core:     core,
		// 		Roll:     roll,
		// 		Material: material,
		// 	})
		// }
	}

	if err := rows.Err(); err != nil {
		return Response{}, err
	}

	// 1.
	// Convert PoMap to slice
	// purchaseorders := make([]PurchaseOrder, 0, len(poMap))
	// for _, po := range poMap {
	// 	purchaseorders = append(purchaseorders, *po)
	// }

	response := Response{
		RecordsTotal:    fmt.Sprintf(`%d`, totalrows),
		RecordsFiltered: fmt.Sprintf(`%d`, totalrows),
	}
	response.Data = datatables

	return response, nil
}

func SuggestVendor(ctx *gin.Context) ([]SuggestionsVendor, error) {
	var prevValue, vendor, nopo string
	Keyword := ctx.DefaultQuery("keyword", "")

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf(`SELECT a.id AS id_vendor, a.vendor, CASE WHEN b.id != '' THEN b.id ELSE '' END AS id_po, CASE WHEN b.nopo != '' THEN b.nopo ELSE '' END AS nopo, CASE WHEN b.type != '' THEN b.type ELSE '' END AS type, CASE WHEN c.detail != '' THEN GROUP_CONCAT(c.detail SEPARATOR ' - ') ELSE '' END AS detail, CASE WHEN d.isi != '' THEN d.isi ELSE '' END AS isi FROM vendor AS a LEFT JOIN po_customer AS b ON a.id = b.id_vendor LEFT JOIN po_item AS c ON c.id_fk = b.id LEFT JOIN setting AS d ON d.id = b.type WHERE a.vendor LIKE '%%%s%%' GROUP BY b.id ORDER BY a.id DESC`, Keyword)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	datas := []SuggestionsVendor{}
	for rows.Next() {
		var vendorid, vendorname, poid, nopocustomer, itemtype, detail, isi string

		if err := rows.Scan(&vendorid, &vendorname, &poid, &nopocustomer, &itemtype, &detail, &isi); err != nil {
			return nil, err
		}

		datas = append(datas, SuggestionsVendor{
			Vendorid:     vendorid,
			VendorName:   vendorname,
			ItemType:     itemtype,
			Detail:       detail,
			Isi:          isi,
			Poid:         poid,
			NoPoCustomer: nopocustomer,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	suggestions := []SuggestionsVendor{}
	if len(datas) > 0 {
		prevValue = ``
		for _, data := range datas {
			vendor = strings.ReplaceAll(data.VendorName, ` `, `_`)
			nopo = strings.ReplaceAll(data.NoPoCustomer, ` `, `_`)

			if nopo == "" {
				suggestions = append(suggestions, SuggestionsVendor{
					Vendorid: data.Vendorid,
					Poid:     ``,
					Label:    `Buat preorder baru`,
					Value:    data.VendorName,
					Category: data.VendorName,
				})

			} else {

				if prevValue == "" {
					suggestions = append(suggestions, SuggestionsVendor{
						Vendorid: data.Vendorid,
						Poid:     ``,
						Label:    `Buat preorder baru`,
						Value:    data.VendorName,
						Category: data.VendorName,
					})

					suggestions = append(suggestions, SuggestionsVendor{
						Vendorid: data.Vendorid,
						Poid:     data.Poid,
						Label:    fmt.Sprintf(`[%s] - %s`, data.Isi, data.Detail),
						Value:    data.VendorName,
						Category: data.VendorName,
					})

				} else if prevValue == vendor {
					suggestions = append(suggestions, SuggestionsVendor{
						Vendorid: data.Vendorid,
						Poid:     ``,
						Label:    `Buat preorder baru`,
						Value:    data.VendorName,
						Category: data.VendorName,
					})

				} else {
					suggestions = append(suggestions, SuggestionsVendor{
						Vendorid: data.Vendorid,
						Poid:     data.Poid,
						Label:    fmt.Sprintf(`[%s] - %s`, data.Isi, data.Detail),
						Value:    data.VendorName,
						Category: data.VendorName,
					})
				}
			}

			prevValue = vendor
		}

	} else {
		suggestions = append(suggestions, SuggestionsVendor{
			Vendorid: ``,
			Poid:     ``,
			Value:    Keyword,
			Label:    `Tidak terdaftar, silakan daftar sebagai vendor baru.`,
			Category: ``,
		})
	}

	return suggestions, nil
}

func SuggestItem(ctx *gin.Context) ([]SuggestionsItem, error) {
	Vendorid := ctx.DefaultQuery("vendorid", "0")
	Poid := ctx.DefaultQuery("poid", "0")

	if Vendorid == "0" {
		return nil, errors.New("invalid ID")
	}

	Vendor, err := strconv.Atoi(Vendorid)
	if err != nil {
		return nil, err
	}

	Po, err := strconv.Atoi(Poid)
	if err != nil {
		return nil, err
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	var item_type, value string
	query := fmt.Sprintf(`SELECT a.type, b.value FROM po_customer AS a LEFT JOIN setting AS b ON b.id = a.type WHERE a.id = %d`, Po)

	err = sql.Connection.QueryRow(query).Scan(&item_type, &value)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			Po = 0
		} else {
			return nil, err
		}
	}

	suggestionsitem := SuggestionsItem{
		PoId:     Po,
		Vendorid: Vendor,
	}

	if Po > 0 {
		var attr Attr
		if err := json.Unmarshal([]byte(value), &attr); err != nil {
			return nil, err
		}

		suggestionsitem.Type = item_type
		suggestionsitem.Attr = attr.Input

		query = fmt.Sprintf(`SELECT item_to, detail, size, price_1, price_2, qty, unit, merk, type, core, gulungan, bahan FROM po_item WHERE id_fk = %d ORDER BY id ASC`, Po)

		rows, err := sql.Connection.Query(query)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		poitem := []PoItem{}
		for rows.Next() {
			var sequence_item, qty int
			var detail, size, price1, price2, unit, merk, itemtype, core, roll, material string

			if err := rows.Scan(&sequence_item, &detail, &size, &price1, &price2, &qty, &unit, &merk, &itemtype, &core, &roll, &material); err != nil {
				return nil, err
			}

			poitem = append(poitem, PoItem{
				SequenceItem: sequence_item,
				Detail:       detail,
				Size:         size,
				Price1:       price1,
				Price2:       price2,
				Qty:          qty,
				Unit:         unit,
				Merk:         merk,
				ItemType:     itemtype,
				Core:         core,
				Roll:         roll,
				Material:     material,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		suggestionsitem.Items = poitem
	}

	return []SuggestionsItem{suggestionsitem}, nil
}

func SuggestType(ctx *gin.Context) ([]SuggestionsType, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := `SELECT id, isi FROM setting WHERE ket ='PO_ITEM'`
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

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

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

func SuggestPO(ctx *gin.Context) ([]SuggestionsPO, error) {
	Keyword := ctx.DefaultQuery("keyword", "")

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf(`SELECT CASE WHEN c.id_fk > 0 THEN c.id_fk ELSE '0' END AS id_fk, CASE WHEN c.item_to > 0 THEN c.item_to ELSE '0' END AS item_to, CASE WHEN b.nopo != '' THEN b.nopo ELSE '' END AS nopo, CASE WHEN b.type != '' THEN b.type ELSE '' END AS type, CASE WHEN c.detail != '' THEN GROUP_CONCAT(c.detail SEPARATOR ' - ') ELSE '' END AS detail, CASE WHEN d.isi != '' THEN d.isi ELSE '' END AS isi FROM po_customer AS b LEFT JOIN po_item AS c ON c.id_fk = b.id LEFT JOIN setting AS d ON d.id = b.type WHERE b.nopo LIKE '%%%s%%' GROUP BY c.id_fk ORDER BY b.id DESC;`, Keyword)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	suggestions := []SuggestionsPO{}
	for rows.Next() {
		var fkid, sequence_item, nopocustomer, itemtype, detail, isi string

		if err := rows.Scan(&fkid, &sequence_item, &nopocustomer, &itemtype, &detail, &isi); err != nil {
			return nil, err
		}

		suggestions = append(suggestions, SuggestionsPO{
			Fkid:     fkid,
			ItemType: itemtype,
			Value:    nopocustomer,
			Category: nopocustomer,
			Label:    fmt.Sprintf(`[%s] - %s`, nopocustomer, isi),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suggestions, nil
}

func GetVendor(Id int) ([]PurchaseOrder, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT a.vendor, b.id_vendor, b.id_company, b.nopo, b.po_date, b.note, b.type, b.ppn FROM vendor AS a LEFT JOIN po_customer AS b ON a.id = b.id_vendor WHERE b.id = %d LIMIT 1", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	purchaseorder := []PurchaseOrder{}
	for rows.Next() {
		var vendorid, companyid int
		var vendorname, ponumber, podate, note, ppn, potype string

		if err := rows.Scan(&vendorname, &vendorid, &companyid, &ponumber, &podate, &note, &potype, &ppn); err != nil {
			return nil, err
		}

		purchaseorder = append(purchaseorder, PurchaseOrder{
			PoId:       Id,
			VendorName: vendorname,
			Vendorid:   vendorid,
			Companyid:  companyid,
			PoNumber:   ponumber,
			PoDate:     podate,
			Note:       note,
			PoType:     potype,
			Ppn:        ppn,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(purchaseorder) < 1 {
		return nil, errors.New("invalid ID")
	}

	return purchaseorder, nil
}

func GetItem(Id int) ([]PoItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT a.id, a.detail, a.size, a.price_1, a.price_2, a.qty, a.unit, a.merk, a.type, a.core, a.gulungan, a.bahan, c.value FROM po_item AS a JOIN po_customer AS b ON b.id = a.id_fk JOIN setting AS c ON c.id = b.type WHERE a.id = %d LIMIT 1", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	poitem := []PoItem{}
	for rows.Next() {
		var itemid, qty int
		var detail, size, price1, price2, unit, merk, itemtype, core, roll, material, value string

		if err := rows.Scan(&itemid, &detail, &size, &price1, &price2, &qty, &unit, &merk, &itemtype, &core, &roll, &material, &value); err != nil {
			return nil, err
		}

		var attr Attr
		if err := json.Unmarshal([]byte(value), &attr); err != nil {
			return nil, err
		}

		poitem = append(poitem, PoItem{
			Id:        itemid,
			Detail:    detail,
			Size:      size,
			Price1:    price1,
			Price2:    price2,
			Qty:       qty,
			Unit:      unit,
			Merk:      merk,
			ItemType:  itemtype,
			Core:      core,
			Roll:      roll,
			Material:  material,
			InputAttr: attr.Input,
			PrintAttr: attr.Print,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(poitem) < 1 {
		return nil, errors.New("invalid ID")
	}

	return poitem, nil
}

func UpdateVendor(Sessionid string, Id int, Vendorid int, Companyid int, Podate string, Note string, Ppn string) ([]PurchaseOrder, error) {
	var nopo string
	var Ppns int
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT nopo FROM po_customer WHERE id = '%d' LIMIT 1", Id)
	if err = sql.Connection.QueryRow(query).Scan(&nopo); err != nil {
		if err.Error() == `sql: no rows in result set` {
			nopo = ``
		} else {
			return nil, err
		}
	}

	// User ID validation
	if nopo == "" {
		return nil, errors.New("invalid po_id")
	}

	ppnConv, err := strconv.Atoi(Ppn)
	if err != nil || ppnConv > 1 {
		Ppns = 0
	} else {
		Ppns = ppnConv
	}

	queryUpdate := fmt.Sprintf("UPDATE po_customer SET id_vendor ='%d', id_company ='%d', po_date ='%s', note ='%s', ppn ='%d' WHERE id ='%d'", Vendorid, Companyid, Podate, Note, Ppns, Id)

	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	// Log capture
	utils.Capture(
		`PO Updated [vendor]`,
		fmt.Sprintf(`PO Id: %d - PO No: %s - VendorId: %d - CompanyId: %d - PO Date: %s - PPN: %d - Note: %s`, Id, nopo, Vendorid, Companyid, Podate, Ppns, Note),
		Sessionid,
	)

	return []PurchaseOrder{}, nil
}

func UpdateItem(Sessionid string, Id int, Detail string, Size string, Price1 string, Price2 string, Qty int, Unit string, Merk string, ItemType string, Core string, Roll string, Material string) ([]PoItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query_id := fmt.Sprintf("SELECT id FROM po_item WHERE id = '%d' LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	// User ID validation
	if !rows_id.Next() {
		return nil, errors.New("invalid item_id")
	}

	defer rows_id.Close()

	queryUpdate := fmt.Sprintf("UPDATE po_item SET detail ='%s', size ='%s', price_1 ='%s', price_2 ='%s', qty ='%d', unit ='%s', merk ='%s', type ='%s', core ='%s', gulungan ='%s', bahan ='%s' WHERE id ='%d'", Detail, Size, utils.PriceFilter(Price1), utils.PriceFilter(Price2), Qty, Unit, Merk, ItemType, Core, Roll, Material, Id)

	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	// Log capture
	utils.Capture(
		`PO Updated [item]`,
		fmt.Sprintf(`Detail: %s - Size: %s - Qty: %d - Unit: %s - Price: %s - Price (sec): %s - Merk: %s - Type: %s - Core: %s - Roll: %s - Material: %s`, Detail, Size, Qty, Unit, Price1, Price2, Merk, ItemType, Core, Roll, Material),
		Sessionid,
	)

	return []PoItem{}, nil
}

func Create(Sessionid string, BodyReq []byte) ([]PurchaseOrder, error) {
	var Ppns int
	var purchaseorder PurchaseOrder
	err := json.Unmarshal([]byte(BodyReq), &purchaseorder)
	if err != nil {
		return nil, err
	}

	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	if purchaseorder.Vendorid < 1 || purchaseorder.Companyid < 1 {
		return nil, errors.New("invalid ID")
	}

	var NoPo, LastPoDate, LastNoPo string
	var AUTO_INCREMENT sql.NullInt64

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	queryAI := `SELECT AUTO_INCREMENT AS ai FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'po_customer'`
	err = sql.Connection.QueryRow(queryAI).Scan(&AUTO_INCREMENT)
	if err != nil {
		return nil, err
	}

	queryLastNoPo := `SELECT po_date, nopo FROM po_customer ORDER BY id DESC LIMIT 1`
	err = sql.Connection.QueryRow(queryLastNoPo).Scan(&LastPoDate, &LastNoPo)
	if err != nil {
		return nil, err
	}

	///////////////////////////////// GENERATE PO NUMBER /////////////////////////////////////
	dateNow := time.Now()
	dateNowFormat := dateNow.Format(config.App.DateFormat_Global)
	dateNowConv, err := strconv.Atoi(fmt.Sprintf("%s%s", dateNowFormat[2:4], dateNowFormat[5:7]))
	if err != nil {
		return nil, err
	}

	LastPoDateParse, err := time.Parse(config.App.DateFormat_Global, LastPoDate)
	if err != nil {
		return nil, err
	}
	LastPoDateFormat := LastPoDateParse.Format(config.App.DateFormat_Global)
	LastPoDateConv, err := strconv.Atoi(fmt.Sprintf("%s%s", LastPoDateFormat[2:4], LastPoDateFormat[5:7]))
	if err != nil {
		return nil, err
	}

	queue, err := strconv.Atoi(LastNoPo[7:])
	if err != nil {
		return nil, err
	}

	queueNumber := queue + 1

	fmt.Println(dateNowConv)
	fmt.Println(LastPoDateConv)
	if dateNowConv > LastPoDateConv {
		NoPo = fmt.Sprintf("PO %d1", dateNowConv)
	} else {
		NoPo = fmt.Sprintf("PO %d%d", LastPoDateConv, queueNumber)
	}
	//////////////////////////////////////////////////////////////////////////////////////////

	ppnConv, err := strconv.Atoi(purchaseorder.Ppn)
	if err != nil || ppnConv > 1 {
		Ppns = 0
	} else {
		Ppns = ppnConv
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	queryPoCustomer := fmt.Sprintf("INSERT INTO po_customer (id_vendor, id_company, po_date, nopo, note, ppn, type, input_by) VALUES ('%d', '%d', '%s', '%s', '%s', '%d', '%s', '%s')", purchaseorder.Vendorid, purchaseorder.Companyid, purchaseorder.PoDate, NoPo, purchaseorder.Note, Ppns, purchaseorder.PoType, Sessionid)
	_, err = tx.Exec(queryPoCustomer)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	queryPoItem := `INSERT INTO po_item (id_fk, item_to, detail, size, price_1, price_2, qty, unit, merk, type, core, gulungan, bahan, hidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmtPoItem, err := tx.Prepare(queryPoItem)
	if err != nil {
		return nil, err
	}

	defer stmtPoItem.Close()

	for index, item := range purchaseorder.Items {
		// validate price is float
		_, err = strconv.ParseFloat(item.Price1, 64)
		if err != nil {
			return nil, fmt.Errorf("price1 format not allowed")
		}

		// validate price is float
		_, err = strconv.ParseFloat(item.Price2, 64)
		if err != nil {
			return nil, fmt.Errorf("price2 format not allowed")
		}

		SequenceItem := index + 1
		_, err := stmtPoItem.Exec(AUTO_INCREMENT, SequenceItem, item.Detail, item.Size, utils.PriceFilter(item.Price1), utils.PriceFilter(item.Price2), item.Qty, item.Unit, item.Merk, item.ItemType, item.Core, item.Roll, item.Material, 0)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Log capture
	utils.Capture(
		`PO Created`,
		fmt.Sprintf(`Vendor: %s - PO No: %s - Date: %s`, purchaseorder.VendorName, NoPo, purchaseorder.PoDate),
		Sessionid,
	)

	return []PurchaseOrder{}, nil
}

func AddItem(Sessionid string, BodyReq []byte) ([]PoItem, error) {
	var item_total int
	var Po AddItemPurchaseOrder

	err := json.Unmarshal([]byte(BodyReq), &Po)
	if err != nil {
		return nil, err
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf(`SELECT COUNT(id) as item_total FROM po_item WHERE id_fk = %d`, Po.Fkid)
	if err = sql.Connection.QueryRow(query).Scan(&item_total); err != nil {
		item_total = 0
	}

	tx, err := sql.Connection.Begin()
	if err != nil {
		return nil, err
	}

	queryPoItem := `INSERT INTO po_item (id_fk, item_to, detail, size, price_1, price_2, qty, unit, merk, type, core, gulungan, bahan, hidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmtPoItem, err := tx.Prepare(queryPoItem)
	if err != nil {
		return nil, err
	}

	defer stmtPoItem.Close()

	log := []map[string]string{}
	for index, item := range Po.Items {
		SequenceItem := item_total + index
		_, err := stmtPoItem.Exec(Po.Fkid, SequenceItem, item.Detail, item.Size, utils.PriceFilter(item.Price1), utils.PriceFilter(item.Price2), item.Qty, item.Unit, item.Merk, item.ItemType, item.Core, item.Roll, item.Material, 0)

		if err != nil {
			tx.Rollback()
			return nil, err
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
		`Item PO Created`,
		fmt.Sprintf(`Fkid: %d - PO No: %s - data: %s`, Po.Fkid, Po.PoNumber, datalog),
		Sessionid,
	)

	return []PoItem{}, nil
}

func Delete(Sessionid string, Id int) ([]PoItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	var id, fkid, sequence_item, detail, size, price1, price2, qty, unit, merk, itemtype, core, roll, material string
	query := fmt.Sprintf(`SELECT id, id_fk, item_to, detail, size, price_1, price_2, qty, unit, merk, type, core, gulungan, bahan FROM po_item where id = %d`, Id)
	if err := sql.Connection.QueryRow(query).Scan(&id, &fkid, &sequence_item, &detail, &size, &price1, &price2, &qty, &unit, &merk, &itemtype, &core, &roll, &material); err != nil {
		if err.Error() == `sql: no rows in result set` {
			return nil, errors.New("invalid ID")
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf("UPDATE po_item SET hidden = 1 WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Log capture
	utils.Capture(
		`PO Deleted`,
		fmt.Sprintf(`Itemid: %d - Fkid: %s - sequence_item: %s -  detail: %s - size: %s - price1: %s - price2: %s - qty: %s - unit: %s - merk: %s - itemtype: %s - core: %s - roll: %s - material: %s`, Id, fkid, sequence_item, detail, size, price1, price2, qty, unit, merk, itemtype, core, roll, material),
		Sessionid,
	)

	return []PoItem{}, nil
}

func GetPrintView(Id int) ([]Print, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT a.vendor, a.address, b.id, b.po_date, b.nopo, b.note, b.ppn, c.id, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.merk, c.type, c.core, c.gulungan, c.bahan, d.isi, d.value FROM vendor AS a JOIN po_customer AS b ON a.id = b.id_vendor JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON b.id = c.id_fk JOIN setting AS d ON d.id = b.type WHERE b.id = %d", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	poMap := make(map[int]*Print)
	// print := []Print{}
	for rows.Next() {
		var poId, itemId, qty int
		var vendorname, vendoraddress, ponumber, podate, note, ppn, detail, size, price1, price2, unit, merk, itemtype, core, roll, material, potype, value string

		if err := rows.Scan(&vendorname, &vendoraddress, &poId, &podate, &ponumber, &note, &ppn, &itemId, &detail, &size, &price1, &price2, &qty, &unit, &merk, &itemtype, &core, &roll, &material, &potype, &value); err != nil {
			return nil, err
		}

		var attr Attr
		if err := json.Unmarshal([]byte(value), &attr); err != nil {
			return nil, err
		}

		if _, exists := poMap[poId]; !exists {
			poMap[poId] = &Print{
				VendorName:    vendorname,
				VendorAddress: vendoraddress,
				PoDate:        podate,
				PoNumber:      ponumber,
				PoType:        potype,
				Note:          note,
				Ppn:           ppn,
				Ttd:           "Iskandar Zulkarnain",
				Items:         []PoItem{},
			}
		}

		if itemId != 0 {
			poMap[poId].Items = append(poMap[poId].Items, PoItem{
				Detail:    detail,
				Size:      size,
				Price1:    price1,
				Price2:    price2,
				Qty:       qty,
				Unit:      unit,
				Merk:      merk,
				ItemType:  itemtype,
				Core:      core,
				Roll:      roll,
				Material:  material,
				InputAttr: attr.Input,
			})
		}

		// print = append(print, Print{
		// 	VendorName:    vendorname,
		// 	VendorAddress: vendoraddress,
		// 	PoDate:        podate,
		// 	PoNumber:      ponumber,
		// 	PoType:        potype,
		// 	Note:          note,
		// 	Ppn:           ppn,
		// 	Ttd:           "Iskandar Zulkarnain",
		// 	Items: []PoItem{
		// 		{
		// 			Detail:   detail,
		// 			Size:     size,
		// 			Price1:   price1,
		// 			Price2:   price2,
		// 			Qty:      qty,
		// 			Unit:     unit,
		// 			Merk:     merk,
		// 			ItemType: itemtype,
		// 			Core:     core,
		// 			Roll:     roll,
		// 			Material: material,
		// 		},
		// 	},
		// 	InputAttr: attr.Input,
		// })
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if len(print) < 1 {
	// 	return nil, errors.New("invalid ID")
	// }

	print := make([]Print, 0, len(poMap))
	for _, po := range poMap {
		print = append(print, *po)
	}

	return print, nil
}

func GetPrintNow(Sessionid string, Id int) ([]Print, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed after all operations
	defer sql.Connection.Close()

	query := fmt.Sprintf("SELECT a.vendor, a.address, b.id, b.po_date, b.nopo, b.note, b.ppn, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.id, c.item_to, d.company, d.address AS alamat, d.email, d.phone, d.logo, e.value FROM vendor AS a LEFT JOIN po_customer AS b ON a.id = b.id_vendor LEFT JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON b.id = c.id_fk LEFT JOIN company AS d ON d.id = b.id_company LEFT JOIN setting AS e ON e.id = b.type WHERE b.id = %d", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	// Generate tanggal print invoice
	durationParse := time.Now()
	print_date := durationParse.Format(config.App.DateFormat_Print)

	// print := []Print{}
	log := map[string]string{}
	poMap := make(map[int]*Print)
	for rows.Next() {
		var poId, itemId, qty, sequenceitem int
		var vendorname, vendoraddress, ponumber, podate, note, ppn, detail, size, price1, price2, unit, companyname, companyaddress, companyemail, companyphone, companylogo, value string

		if err := rows.Scan(&vendorname, &vendoraddress, &poId, &podate, &ponumber, &note, &ppn, &detail, &size, &price1, &price2, &qty, &unit, &itemId, &sequenceitem, &companyname, &companyaddress, &companyemail, &companyphone, &companylogo, &value); err != nil {
			return nil, err
		}

		var attr Attr
		if err := json.Unmarshal([]byte(value), &attr); err != nil {
			return nil, err
		}

		if _, exists := poMap[poId]; !exists {
			poMap[poId] = &Print{
				VendorName:     vendorname,
				VendorAddress:  vendoraddress,
				PoDate:         podate,
				PoNumber:       ponumber,
				Note:           note,
				Ppn:            ppn,
				Ttd:            "Iskandar Zulkarnain",
				CompanyName:    companyname,
				CompanyAddress: companyaddress,
				Companyemail:   companyemail,
				CompanyPhone:   companyphone,
				CompanyLogo:    companylogo,
				PrintAttr:      attr.Print,
				PrintDate:      print_date,
				Items:          []PoItem{},
			}
		}

		if itemId != 0 {
			// Parse prices and calculate subtotal
			price2float, _ := strconv.ParseFloat(price2, 64)
			itemSubtotal := price2float * float64(qty)

			// Update subtotal
			subtotal := 0.0
			if poMap[poId].Subtotal != "" {
				subtotal, _ = strconv.ParseFloat(poMap[poId].Subtotal, 64)
			}
			subtotal += itemSubtotal

			// Update poMap with subtotal
			poMap[poId].Subtotal = fmt.Sprintf("%.2f", subtotal)

			poMap[poId].Items = append(poMap[poId].Items, PoItem{
				SequenceItem: sequenceitem,
				Detail:       detail,
				Size:         size,
				Price1:       price1,
				Price2:       price2,
				Qty:          qty,
				Unit:         unit,
				Subtotal:     fmt.Sprintf("%.2f", itemSubtotal),
			})
		}

		log["Poid"] = fmt.Sprintf(`%d`, poId)
		log["Po No"] = ponumber
	}

	for _, po := range poMap {
		subtotal, _ := strconv.ParseFloat(po.Subtotal, 64)
		taxTotal, _ := strconv.ParseFloat(po.Taxtotal, 64)

		if po.Ppn == `0` {
			taxTotal = 0
		} else {
			taxTotal = subtotal * 11 / 100
		}

		// Calculate tax total and final total
		total := subtotal + taxTotal

		// Set the calculated values in the struct
		po.Taxtotal = fmt.Sprintf("%.2f", taxTotal)
		po.Total = fmt.Sprintf("%.2f", total)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if len(print) < 1 {
	// 	return nil, errors.New("invalid ID")
	// }

	print := make([]Print, 0, len(poMap))
	for _, po := range poMap {
		print = append(print, *po)
	}

	// Log capture
	datalog, _ := json.Marshal(log)
	utils.Capture(
		`PO Print`,
		fmt.Sprintf(`Data: %s`, datalog),
		Sessionid,
	)

	return print, nil
}

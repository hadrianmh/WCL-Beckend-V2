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

type PurchaseOrder struct {
	PoId          int      `json:"po_id,omitempty"`
	Vendorid      int      `json:"vendorid,omitempty"`
	Companyid     int      `json:"companyid,omitempty"`
	PoDate        string   `json:"po_date,omitempty"`
	PoNumber      string   `json:"po_number,omitempty"`
	Note          string   `json:"note,omitempty"`
	Ppn           string   `json:"ppn,omitempty"`
	PoType        string   `json:"po_type,omitempty"`
	MonthlyReport string   `json:"monthly_report,omitempty"`
	InputBy       int      `json:"inputby,omitempty"`
	VendorName    string   `json:"vendor,omitempty"`
	VendorAddress string   `json:"vendor_address,omitempty"`
	CompanyName   string   `json:"company,omitempty"`
	UserId        int      `json:"userid,omitempty"`
	UserName      string   `json:"user,omitempty"`
	Items         []PoItem `json:"items,omitempty"`
}

type PoItem struct {
	Id           int    `json:"itemid,omitempty"`
	Fkid         int    `json:"fkid,omitempty"`
	SequenceItem int    `json:"sequence_item,omitempty"`
	Detail       string `json:"detail,omitempty"`
	Size         string `json:"size,omitempty"`
	Price1       string `json:"price1,omitempty"`
	Price2       string `json:"price2,omitempty"`
	Qty          int    `json:"qty,omitempty"`
	Unit         string `json:"unit,omitempty"`
	Merk         string `json:"merk,omitempty"`
	ItemType     string `json:"item_type,omitempty"`
	Core         string `json:"core,omitempty"`
	Roll         string `json:"roll,omitempty"`
	Material     string `json:"material,omitempty"`
	InputAttr    string `json:"inputattr,omitempty"`
	PrintAttr    string `json:"printattr,omitempty"`
}

type Attr struct {
	Input string `json:"input,omitempty"`
	Print string `json:"print,omitempty"`
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
	Price1      string `json:"price1"`
	Price2      string `json:"price2"`
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
}

type Print struct {
	VendorName     string   `json:"vendor,omitempty"`
	VendorAddress  string   `json:"vendor_address,omitempty"`
	PoDate         string   `json:"po_date,omitempty"`
	PoNumber       string   `json:"po_number,omitempty"`
	PoType         string   `json:"po_type,omitempty"`
	Note           string   `json:"note,omitempty"`
	Ppn            string   `json:"ppn,omitempty"`
	Ttd            string   `json:"ttd,omitempty"`
	InputAttr      string   `json:"inputattr,omitempty"`
	PrintAttr      string   `json:"printattr,omitempty"`
	CompanyName    string   `json:"companyname,omitempty"`
	CompanyAddress string   `json:"companyaddress,omitempty"`
	Companyemail   string   `json:"companyemail,omitempty"`
	CompanyPhone   string   `json:"companyphone,omitempty"`
	CompanyLogo    string   `json:"companylogo,omitempty"`
	Items          []PoItem `json:"items,omitempty"`
}

func Get(ctx *gin.Context) ([]DataTables, error) {
	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	var query, monthlyreport string
	// idParam := ctx.DefaultQuery("id", "0")
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

	monthlyreport = MonthlyReportParam
	if MonthlyReportParam == "" {
		DateNow := time.Now()
		monthlyreport = DateNow.Format(config.App.DateFormat_MonthlyReport)
	}

	// id, err := strconv.Atoi(idParam)
	// if err != nil || id < 1 {
	query = fmt.Sprintf(`SELECT a.id AS id_po, a.po_date, a.nopo, a.note, a.ppn, a.input_by, b.vendor, c.id AS id_po_item, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.merk, c.type, c.core, c.gulungan, c.bahan, d.id AS userid, d.name, e.company, f.isi FROM po_customer AS a JOIN vendor AS b ON a.id_vendor = b.id JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON a.id = c.id_fk JOIN user AS d ON a.input_by = d.id JOIN company AS e ON e.id = a.id_company JOIN setting AS f ON f.id = a.type WHERE a.po_date LIKE '%s%%' ORDER BY a.id DESC LIMIT %d OFFSET %d`, strings.ReplaceAll(monthlyreport, "/", "-"), limit, offset)
	// } else {
	// 	query = fmt.Sprintf(`SELECT a.id AS id_po, a.po_date, a.nopo, a.note, a.ppn, a.input_by, b.vendor, c.id AS id_po_item, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.merk, c.type, c.core, c.gulungan, c.bahan, d.id AS userid, d.name, e.company, f.isi FROM po_customer AS a JOIN vendor AS b ON a.id_vendor = b.id JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON a.id = c.id_fk JOIN user AS d ON a.input_by = d.id JOIN company AS e ON e.id = a.id_company JOIN setting AS f ON f.id = a.type WHERE a.id = %d`, id)
	// }

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 2. purchaseorders := []PurchaseOrder{}
	// 1. poMap := make(map[int]*PurchaseOrder)
	datatables := []DataTables{}
	for rows.Next() {
		var poId, inputBy, itemId, userId, qty int
		var poDate, poNumber, note, ppn, vendorName, detail, size, price1, price2, unit, merk, itemType, core, roll, material, userName, companyName, poType string

		if err := rows.Scan(
			&poId, &poDate, &poNumber, &note, &ppn, &inputBy, &vendorName,
			&itemId, &detail, &size, &price1, &price2, &qty, &unit, &merk,
			&itemType, &core, &roll, &material, &userId, &userName, &companyName, &poType,
		); err != nil {
			return nil, err
		}

		datatables = append(datatables, DataTables{
			PoId:        poId,
			PoDate:      poDate,
			CompanyName: companyName,
			VendorName:  vendorName,
			PoNumber:    poNumber,
			PoType:      poType,
			Detail:      detail,
			Size:        size,
			Price1:      price1,
			Price2:      price2,
			Qty:         qty,
			Unit:        unit,
			Merk:        merk,
			ItemType:    itemType,
			Core:        core,
			Roll:        roll,
			Material:    material,
			Note:        note,
			UserName:    userName,
			ItemId:      itemId,
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
		return nil, err
	}

	// 1.
	// Convert PoMap to slice
	// purchaseorders := make([]PurchaseOrder, 0, len(poMap))
	// for _, po := range poMap {
	// 	purchaseorders = append(purchaseorders, *po)
	// }

	return datatables, nil
}

func GetVendor(Id int) ([]PurchaseOrder, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

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

func UpdateVendor(Id int, Vendorid int, Companyid int, Podate string, Note string, Ppn string) ([]PurchaseOrder, error) {
	var Ppns int
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query_id := fmt.Sprintf("SELECT id FROM po_customer WHERE id = '%d' LIMIT 1", Id)
	rows_id, err := sql.Connection.Query(query_id)
	if err != nil {
		return nil, err
	}

	// User ID validation
	if !rows_id.Next() {
		return nil, errors.New("invalid po_id")
	}

	defer rows_id.Close()

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

	return []PurchaseOrder{}, nil
}

func UpdateItem(Id int, Detail string, Size string, Price1 string, Price2 string, Qty int, Unit string, Merk string, ItemType string, Core string, Roll string, Material string) ([]PoItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

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

	queryUpdate := fmt.Sprintf("UPDATE po_item SET detail ='%s', size ='%s', price_1 ='%s', price_2 ='%s', qty ='%d', unit ='%s', merk ='%s', type ='%s', core ='%s', gulungan ='%s', bahan ='%s' WHERE id ='%d'", Detail, Size, Price1, Price2, Qty, Unit, Merk, ItemType, Core, Roll, Material, Id)

	update, err := sql.Connection.Query(queryUpdate)
	if err != nil {
		return nil, err
	}

	defer update.Close()

	return []PoItem{}, nil
}

func Create(BodyReq []byte) ([]PurchaseOrder, error) {
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

	queryPoCustomer := fmt.Sprintf("INSERT INTO po_customer (id_vendor, id_company, po_date, nopo, note, ppn, type, input_by) VALUES ('%d', '%d', '%s', '%s', '%s', '%d', '%s', '%d')", purchaseorder.Vendorid, purchaseorder.Companyid, purchaseorder.PoDate, NoPo, purchaseorder.Note, Ppns, purchaseorder.PoType, purchaseorder.InputBy)
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

	return []PurchaseOrder{}, nil
}

func Delete(Id int) ([]PoItem, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("UPDATE po_item SET hidden = 1 WHERE id = %d", Id)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return []PoItem{}, nil
}

func GetPrintView(Id int) ([]Print, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

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
				Detail:   detail,
				Size:     size,
				Price1:   price1,
				Price2:   price2,
				Qty:      qty,
				Unit:     unit,
				Merk:     merk,
				ItemType: itemtype,
				Core:     core,
				Roll:     roll,
				Material: material,
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

func GetPrintNow(Id int) ([]Print, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT a.vendor, a.address, b.id, b.po_date, b.nopo, b.note, b.ppn, c.detail, c.size, c.price_1, c.price_2, c.qty, c.unit, c.id, c.item_to, d.company, d.address AS alamat, d.email, d.phone, d.logo, e.value FROM vendor AS a LEFT JOIN po_customer AS b ON a.id = b.id_vendor LEFT JOIN (SELECT * FROM po_item WHERE hidden = 0) AS c ON b.id = c.id_fk LEFT JOIN company AS d ON d.id = b.id_company LEFT JOIN setting AS e ON e.id = b.type WHERE b.id = %d", Id)

	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// print := []Print{}
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
				Items:          []PoItem{},
			}
		}

		if itemId != 0 {
			poMap[poId].Items = append(poMap[poId].Items, PoItem{
				SequenceItem: sequenceitem,
				Detail:       detail,
				Size:         size,
				Price1:       price1,
				Price2:       price2,
				Qty:          qty,
				Unit:         unit,
			})
		}
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

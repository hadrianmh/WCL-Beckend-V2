package controllers

import (
	"backend/adapters"
	"backend/services/auth"
	"backend/services/company"
	"backend/services/customer"
	deliveryorder "backend/services/delivery_order"
	invoice "backend/services/invoice"
	"backend/services/metrics"
	purchaseorder "backend/services/purchase_order"
	salesorder "backend/services/sales_order"
	"backend/services/sortdata"
	"backend/services/user"
	"backend/services/vendor"
	"backend/services/workorder"
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BodyRequestAuth struct {
	Action        string `json:"action" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password"`
	ResfreshToken string `json:"refresh_token"`
}

type BodyRequestUser struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role" binding:"required"`
	Status   int    `json:"status,omitempty"`
	Account  int    `json:"account,omitempty"`
	Hidden   int    `json:"hidden,omitempty"`
}

type BodyRequestCompany struct {
	Id          int    `json:"id,omitempty"`
	CompanyName string `json:"company" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Logo        string `json:"logo,omitempty"`
	InputBy     int    `json:"inputby,omitempty"`
	Hidden      int    `json:"hidden,omitempty"`
}

type BodyRequestVendor struct {
	Id         int    `json:"id,omitempty"`
	VendorName string `json:"vendor" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Phone      int    `json:"phone" binding:"required"`
	InputBy    int    `json:"inputby,omitempty"`
	Hidden     int    `json:"hidden,omitempty"`
}

type BodyRequestCustomer struct {
	Id           int    `json:"id,omitempty"`
	CustomerName string `json:"customername,omitempty" binding:"required"`
	Address      string `json:"address,omitempty" binding:"required"`
	City         string `json:"city,omitempty" binding:"required"`
	Country      string `json:"country,omitempty" binding:"required"`
	Province     string `json:"province,omitempty" binding:"required"`
	PostalCode   string `json:"postalcode,omitempty" binding:"required"`
	Phone        string `json:"phone,omitempty" binding:"required"`
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

type BodyRequestPurchaseOrder struct {
	PoId          int      `json:"po_id,omitempty"`
	Vendorid      int      `json:"vendorid,omitempty" binding:"required"`
	Companyid     int      `json:"companyid,omitempty" binding:"required"`
	PoDate        string   `json:"po_date,omitempty" binding:"required"`
	PoNumber      string   `json:"po_number,omitempty"`
	Note          string   `json:"note,omitempty"`
	Ppn           string   `json:"tax,omitempty" binding:"required"`
	PoType        string   `json:"po_type,omitempty" binding:"required"`
	MonthlyReport string   `json:"monthly_report,omitempty"`
	InputBy       int      `json:"inputby,omitempty"`
	Items         []PoItem `json:"po_item,omitempty"`
}

type BodyRequestAddItemPurchaseOrder struct {
	Fkid  int      `json:"fkid,omitempty" binding:"required"`
	Items []PoItem `json:"po_item,omitempty"`
}

type PoItem struct {
	Id           int    `json:"itemid,omitempty"`
	Fkid         int    `json:"fkid,omitempty"`
	SequenceItem int    `json:"sequence_item,omitempty"`
	Detail       string `json:"detail,omitempty"`
	Size         string `json:"size,omitempty"`
	Price1       string `json:"price_1,omitempty" binding:"required"`
	Price2       string `json:"price_2,omitempty"`
	Qty          int    `json:"qty,omitempty" binding:"required"`
	Unit         string `json:"unit,omitempty"`
	Merk         string `json:"merk,omitempty"`
	ItemType     string `json:"item_type,omitempty"`
	Core         string `json:"core,omitempty"`
	Roll         string `json:"roll,omitempty"`
	Material     string `json:"material,omitempty"`
	Hidden       int    `json:"hidden,omitempty"`
}

type BodyRequestSalesOrder struct {
	Id           int              `json:"poid,omitempty"`
	FkId         int              `json:"fkid,omitempty"`
	CompanyId    int              `json:"companyid"`
	CustomerId   int              `json:"customerid,omitempty"`
	CustomerName string           `json:"customer"`
	PoDate       string           `json:"po_date"`
	NoPoCustomer string           `json:"po_customer,omitempty"`
	OrderGrade   int              `json:"order_grade"`
	Ppn          int              `json:"tax"`
	Items        []SalesOrderItem `json:"items,omitempty"`
}

type SalesOrderItem struct {
	PoitemId     int    `json:"poitemid,omitempty"`
	WoitemId     int    `json:"woitemid,omitempty"`
	Detail       int    `json:"detail"`
	Item         string `json:"item"`
	Size         string `json:"size"`
	Qty          int64  `json:"qty"`
	Unit         string `json:"unit"`
	Price        string `json:"price"`
	Qore         string `json:"qore"`
	Lin          string `json:"lin"`
	Roll         string `json:"roll"`
	Material     string `json:"ingredient"`
	Volume       int64  `json:"volume"`
	Porporasi    int    `json:"porporasi"`
	UkBahanBaku  string `json:"uk_bahan_baku"`
	QtyBahanBaku string `json:"qty_bahan_baku"`
	Sources      string `json:"sources"`
	Merk         string `json:"merk"`
	Type         string `json:"type"`
	Note         string `json:"annotation"`
	Etc1         string `json:"etc1,omitempty"`
	Etc2         string `json:"etc2,omitempty"`
}

type SalesOrderShippingCost struct {
	Id        int    `json:"id"`
	Cost      string `json:"cost"`
	Ekspedisi string `json:"ekspedisi"`
	Uom       string `json:"uom"`
	Jml       string `json:"jml"`
}

type BodyRequestWorkOrder struct {
	Id           int    `json:"id"`
	SequenceItem int    `json:"item_to"`
	SpkDate      string `json:"spk_date"`
	OrderStatus  int    `json:"order_status"`
	InputBy      int    `json:"input_by"`
}

type WorkOrder struct {
	DateNow      string  `json:"tgl,omitempty"`
	NoPoCustomer string  `json:"po_customer,omitempty"`
	CustomerName string  `json:"customer,omitempty"`
	NoSpk        string  `json:"no_spk,omitempty"`
	Note         string  `json:"annotation,omitempty"`
	Size         string  `json:"size,omitempty"`
	UkBahanBaku  string  `json:"uk_bahan_baku,omitempty"`
	Material     string  `json:"ingredient,omitempty"`
	Porporasi    string  `json:"porporasi,omitempty"`
	Roll         string  `json:"roll,omitempty"`
	Qore         string  `json:"qore,omitempty"`
	Lin          string  `json:"lin,omitempty"`
	QtyBahanBaku string  `json:"qty_bahan_baku,omitempty"`
	Ttl          float64 `json:"ttl,omitempty"` // qty produksi
	SatuanUnit   string  `json:"satuanunit,omitempty"`
	Isi          int     `json:"isi,omitempty"`
	Ttd          string  `json:"ttd,omitempty"`
}

type BodyRequestDeliveryOrder struct {
	Id      int                  `json:"id"`
	SjDate  string               `json:"sj_date"`
	NoSj    string               `json:"no_sj"`
	Shipto  string               `json:"shipto"`
	Courier string               `json:"courier"`
	Resi    string               `json:"resi"`
	InputBy int                  `json:"input_by"`
	Items   []DeliveryOrdersItem `json:"do_item,omitempty"`
}

type DeliveryOrdersItem struct {
	SequenceItem int   `json:"item_to,omitempty"`
	Qty          int64 `json:"qty,omitempty"`
}

type BodyRequestInvoice struct {
	Id           string  `json:"id"`
	InvoiceDate  string  `json:"invoice_date,omitempty"`
	CustomerName string  `json:"customer,omitempty"`
	CompanyName  string  `json:"company,omitempty"`
	Shipto       string  `json:"shipto,omitempty"`
	NoPoCustomer string  `json:"po_customer,omitempty"`
	NoSo         string  `json:"no_so,omitempty"`
	NoSj         string  `json:"no_sj,omitempty"`
	NoInvoice    string  `json:"no_invoice,omitempty"`
	SendQty      string  `json:"send_qty,omitempty"`
	Item         string  `json:"item,omitempty"`
	Unit         string  `json:"unit,omitempty"`
	Price        float64 `json:"price,omitempty"`
	Bill         string  `json:"bill,omitempty"`
	Ppn          string  `json:"ppn,omitempty"`
	Total        string  `json:"total,omitempty"`
	Cost         string  `json:"cost,omitempty"`
	Address      string  `json:"address,omitempty"`
	SName        string  `json:"sname,omitempty"`
	SPhone       string  `json:"sphone,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Billto       string  `json:"billto,omitempty"`
	Bank         string  `json:"bank,omitempty"`
	Ttd          string  `json:"ttd,omitempty"`
	Note         string  `json:"note,omitempty"`
	InputBy      int     `json:"input_by,omitempty"`
}

type SortDataArchiveRequest struct {
	Data string `json:"data" binding:"required"`
	From string `json:"from" binding:"required"`
}

type SortDataCounterRequest struct {
	Type      string `json:"type" binding:"required"` // single or periode
	StartDate string `json:"startdate" binding:"required"`
	EndDate   string `json:"Enddate"`
}

func Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"message": "Welcome to WCL microservices."}})
}

func GetDashboard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"message": "Welcome to WCL Dashboard."}})
}

func SortData_Archive(ctx *gin.Context) {
	get, err := sortdata.GetArchive(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SortData_Counter(ctx *gin.Context) {
	get, err := sortdata.GetCounter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func Ping(ctx *gin.Context) {
	ping, err := adapters.NewSql()
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	defer ping.Connection.Close()

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"message": "Pinged the database!"}})
}

func Auth(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestAuth
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	authCheck, err := auth.Init(BodyReq.Action, BodyReq.Email, BodyReq.Password, BodyReq.ResfreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": authCheck})
}

func GetUser(ctx *gin.Context) {
	get, err := user.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func CreateUser(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestUser
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	create, err := user.Create(BodyReq.Name, BodyReq.Email, BodyReq.Password, BodyReq.Role, BodyReq.Account, BodyReq.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func UpdateUser(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestUser
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	update, err := user.Update(BodyReq.Id, BodyReq.Name, BodyReq.Email, BodyReq.Role, BodyReq.Status, BodyReq.Account, BodyReq.Hidden)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	delete, err := user.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func GetCompany(ctx *gin.Context) {
	get, err := company.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func CreateCompany(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestCompany
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := company.Create(sessionid.(string), BodyReq.CompanyName, BodyReq.Address, BodyReq.Email, BodyReq.Phone, BodyReq.Logo, BodyReq.InputBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func UpdateCompany(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestCompany
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := company.Update(sessionid.(string), BodyReq.Id, BodyReq.CompanyName, BodyReq.Address, BodyReq.Email, BodyReq.Phone, BodyReq.Logo, BodyReq.InputBy, BodyReq.Hidden)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func DeleteCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	delete, err := company.Delete(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func GetVendor(ctx *gin.Context) {
	get, err := vendor.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func CreateVendor(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestVendor
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := vendor.Create(sessionid.(string), BodyReq.VendorName, BodyReq.Address, BodyReq.Phone, BodyReq.InputBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func UpdateVendor(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestVendor
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := vendor.Update(sessionid.(string), BodyReq.Id, BodyReq.VendorName, BodyReq.Address, BodyReq.Phone, BodyReq.InputBy, BodyReq.Hidden)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func DeleteVendor(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	delete, err := vendor.Delete(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func GetCustomer(ctx *gin.Context) {
	get, err := customer.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func CreateCustomer(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestCustomer
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := customer.Create(sessionid.(string), BodyReq.CustomerName, BodyReq.Address, BodyReq.City, BodyReq.Country, BodyReq.Province, BodyReq.PostalCode, BodyReq.Phone, BodyReq.SName, BodyReq.SAddress, BodyReq.SCity, BodyReq.SCountry, BodyReq.SProvince, BodyReq.SPostalCode, BodyReq.SPhone, BodyReq.InputBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func UpdateCustomer(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestCustomer
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := customer.Update(sessionid.(string), BodyReq.Id, BodyReq.CustomerName, BodyReq.Address, BodyReq.City, BodyReq.Country, BodyReq.Province, BodyReq.PostalCode, BodyReq.Phone, BodyReq.SName, BodyReq.SAddress, BodyReq.SCity, BodyReq.SCountry, BodyReq.SProvince, BodyReq.SPostalCode, BodyReq.SPhone, BodyReq.InputBy, BodyReq.Hidden)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func DeleteCustomer(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	delete, err := customer.Delete(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func GetPurchaseOrder(ctx *gin.Context) {
	get, err := purchaseorder.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func GetPurchaseOrder_SuggestVendor(ctx *gin.Context) {
	get, err := purchaseorder.SuggestVendor(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_SuggestItem(ctx *gin.Context) {
	get, err := purchaseorder.SuggestItem(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_SuggestType(ctx *gin.Context) {
	get, err := purchaseorder.SuggestType(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_SuggestAttr(ctx *gin.Context) {
	get, err := purchaseorder.SuggestAttr(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_SuggestPO(ctx *gin.Context) {
	get, err := purchaseorder.SuggestPO(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_Vendor(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := purchaseorder.GetVendor(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_Item(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := purchaseorder.GetItem(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func CreatePurchaseOrder(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestPurchaseOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := purchaseorder.Create(sessionid.(string), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func AddItemPurchaseOrder(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestAddItemPurchaseOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := purchaseorder.AddItem(sessionid.(string), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func UpdatePurchaseOrder_Vendor(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestPurchaseOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := purchaseorder.UpdateVendor(sessionid.(string), BodyReq.PoId, BodyReq.Vendorid, BodyReq.Companyid, BodyReq.PoDate, BodyReq.Note, BodyReq.Ppn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func UpdatePurchaseOrder_Item(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq PoItem
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := purchaseorder.UpdateItem(sessionid.(string), BodyReq.Id, BodyReq.Detail, BodyReq.Size, BodyReq.Price1, BodyReq.Price2, BodyReq.Qty, BodyReq.Unit, BodyReq.Merk, BodyReq.ItemType, BodyReq.Core, BodyReq.Roll, BodyReq.Material)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func DeletePurchaseOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	delete, err := purchaseorder.Delete(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func GetPurchaseOrder_PrintView(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := purchaseorder.GetPrintView(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetPurchaseOrder_PrintNow(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	get, err := purchaseorder.GetPrintNow(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_Get(ctx *gin.Context) {
	get, err := salesorder.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func SalesOrder_SuggestType(ctx *gin.Context) {
	get, err := salesorder.SuggestType(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_SuggestCustomer(ctx *gin.Context) {
	get, err := salesorder.SuggestCustomer(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_SuggestItem(ctx *gin.Context) {
	get, err := salesorder.SuggestItem(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_SuggestAttr(ctx *gin.Context) {
	get, err := salesorder.SuggestAttr(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_SuggestSO(ctx *gin.Context) {
	get, err := salesorder.SuggestSO(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_Create(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestSalesOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := salesorder.Create(sessionid.(string), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func SalesOrder_AddItem(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestSalesOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := salesorder.AddItem(sessionid.(string), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func SalesOrder_GetCustomer(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := salesorder.GetCustomer(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_GetItem(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := salesorder.GetItem(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_UpdateCustomer(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestSalesOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := salesorder.UpdateCustomer(sessionid.(string), BodyReq.Id, BodyReq.CompanyId, BodyReq.CustomerId, BodyReq.CustomerName, BodyReq.OrderGrade, BodyReq.PoDate, BodyReq.NoPoCustomer, BodyReq.Ppn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func SalesOrder_UpdateItem(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq SalesOrderItem
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	update, err := salesorder.UpdateItem(sessionid.(string), BodyReq.PoitemId, BodyReq.WoitemId, BodyReq.Item, BodyReq.Size, BodyReq.UkBahanBaku, BodyReq.Qore, BodyReq.Lin, BodyReq.QtyBahanBaku, BodyReq.Roll, BodyReq.Material, BodyReq.Unit, BodyReq.Volume, BodyReq.Note, BodyReq.Price, BodyReq.Qty, BodyReq.Sources, BodyReq.Porporasi, BodyReq.Detail, BodyReq.Merk, BodyReq.Type, BodyReq.Etc1, BodyReq.Etc2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func SalesOrder_GetShipCost(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := salesorder.GetShippingCost(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func SalesOrder_UpdateShipCost(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq SalesOrderShippingCost
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	// Id int, detail string, cost string, ekspedisi string, uom string, jml string
	update, err := salesorder.UpdateShipCost(sessionid.(string), BodyReq.Id, BodyReq.Cost, BodyReq.Ekspedisi, BodyReq.Uom, BodyReq.Jml)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": update}})
}

func SalesOrder_Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	delete, err := salesorder.Delete(sessionid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func Workorder_Get(ctx *gin.Context) {
	get, err := workorder.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func Workorder_Create(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestWorkOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := workorder.Create(BodyReq.Id, BodyReq.SequenceItem, BodyReq.SpkDate, BodyReq.OrderStatus, sessionid.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func Workorder_GetProcess(ctx *gin.Context) {
	id := ctx.Param("id")
	sequence_item := ctx.Param("sequence_item")

	if id == `` || sequence_item == `` {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	get, err := workorder.GetProcess(id, sequence_item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func Workorder_Printview(ctx *gin.Context) {
	idStr := ctx.Param("wocusid")
	sequence_itemStr := ctx.Param("sequence_item")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	sequence_item, err := strconv.Atoi(sequence_itemStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := workorder.Printview(id, sequence_item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func Workorder_Printnow(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq WorkOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	create, err := workorder.Printnow(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func DeliveryOrder_Get(ctx *gin.Context) {
	get, err := deliveryorder.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func DeliveryOrder_GetWaiting(ctx *gin.Context) {
	get, err := deliveryorder.GetWaiting(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": get})
}

func DeliveryOrder_GetItem(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := deliveryorder.GetItem(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func DeliveryOrder_GetNumber(ctx *gin.Context) {
	get, err := deliveryorder.GetNumber()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func DeliveryOrder_Create(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestDeliveryOrder
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	create, err := deliveryorder.Create(sessionid.(string), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func DeliveryOrder_Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	delete, err := deliveryorder.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": delete}})
}

func DeliveryOrder_Printview(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := deliveryorder.Printview(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func DeliveryOrder_Printnow(ctx *gin.Context) {
	idStr := ctx.Param("id")
	ttd := ctx.Param("ttd")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := deliveryorder.Printnow(id, ttd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetInvoice(ctx *gin.Context) {
	get, err := invoice.Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func CreateInvoice(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestInvoice
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	create, err := invoice.Create(BodyReq.Id, BodyReq.InvoiceDate, BodyReq.InputBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func GetInvoice_Printview(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	get, err := invoice.Printview(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func GetInvoice_Printnow(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestInvoice
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	create, err := invoice.Printnow(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": create}})
}

func PaidInvoice(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestInvoice
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	paid, err := invoice.Paid(BodyReq.Id, BodyReq.InvoiceDate, BodyReq.Note)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": paid}})
}

func UnpaidInvoice(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequestInvoice
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	unpaid, err := invoice.UnPaid(BodyReq.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": unpaid}})
}

func DeleteInvoice(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	del, err := invoice.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": del}})
}

func Metrics_Notification(ctx *gin.Context) {
	// Validation userid from access_token set in context as uniqid
	sessionid, exists := ctx.Get("uniqid")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusUnauthorized, "status": "error", "response": gin.H{"message": "invalid token"}})
		return
	}

	get, err := metrics.Notification(sessionid.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func Metrics_SoTracking(ctx *gin.Context) {
	get, err := metrics.SoTracking(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

func Metrics_Static(ctx *gin.Context) {
	get, err := metrics.Static(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "response": gin.H{"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "status": "success", "response": gin.H{"data": get}})
}

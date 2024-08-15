package routes

import (
	"backend/controllers"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	route := app

	// General v1
	ApiV1 := route.Group("/api/v1")
	ApiV1.GET("/", controllers.Home)
	ApiV1.GET("/ping", controllers.Ping)
	ApiV1.POST("/auth", controllers.Auth)

	// Dashboard v1
	ApiV1Dashboard := route.Group("api/v1/dashboard")
	ApiV1Dashboard.Use(utils.AuthenticateJWT())
	{
		ApiV1Dashboard.GET("", controllers.GetDashboard)

		// Metrics
		ApiV1Dashboard.GET("/metrics/notification", controllers.Metrics_Notification)
		ApiV1Dashboard.GET("/metrics/so-tracking", controllers.Metrics_SoTracking)
		ApiV1Dashboard.GET("/metrics/static", controllers.Metrics_Static)

		// Sortdata management
		ApiV1Dashboard.GET("/sortdata/archive", controllers.SortData_Archive)
		ApiV1Dashboard.GET("/sortdata/counter", controllers.SortData_Counter)

		// Users management
		ApiV1Dashboard.GET("/user", controllers.GetUser)
		ApiV1Dashboard.POST("/user", controllers.CreateUser)
		ApiV1Dashboard.PUT("/user", controllers.UpdateUser)
		ApiV1Dashboard.DELETE("/user/:id", controllers.DeleteUser)

		// Company management
		ApiV1Dashboard.GET("/company", controllers.GetCompany)
		ApiV1Dashboard.POST("/company", controllers.CreateCompany)
		ApiV1Dashboard.PUT("/company", controllers.UpdateCompany)
		ApiV1Dashboard.DELETE("/company/:id", controllers.DeleteCompany)

		// Vendor management
		ApiV1Dashboard.GET("/vendor", controllers.GetVendor)
		ApiV1Dashboard.POST("/vendor", controllers.CreateVendor)
		ApiV1Dashboard.PUT("/vendor", controllers.UpdateVendor)
		ApiV1Dashboard.DELETE("/vendor/:id", controllers.DeleteVendor)

		// Customer management
		ApiV1Dashboard.GET("/customer", controllers.GetCustomer)
		ApiV1Dashboard.POST("/customer", controllers.CreateCustomer)
		ApiV1Dashboard.PUT("/customer", controllers.UpdateCustomer)
		ApiV1Dashboard.DELETE("/customer/:id", controllers.DeleteCustomer)

		// Purchase Order management
		ApiV1Dashboard.GET("/purchase-order", controllers.GetPurchaseOrder)
		ApiV1Dashboard.GET("/purchase-order/suggest/vendor", controllers.GetPurchaseOrder_SuggestVendor)
		ApiV1Dashboard.GET("/purchase-order/suggest/item", controllers.GetPurchaseOrder_SuggestItem)
		ApiV1Dashboard.GET("/purchase-order/suggest/type", controllers.GetPurchaseOrder_SuggestType)
		ApiV1Dashboard.GET("/purchase-order/suggest/attr", controllers.GetPurchaseOrder_SuggestAttr)
		ApiV1Dashboard.GET("/purchase-order/suggest/po", controllers.GetPurchaseOrder_SuggestPO)
		ApiV1Dashboard.POST("/purchase-order", controllers.CreatePurchaseOrder)
		ApiV1Dashboard.POST("/purchase-order/item", controllers.AddItemPurchaseOrder)
		ApiV1Dashboard.PUT("/purchase-order/vendor", controllers.UpdatePurchaseOrder_Vendor)
		ApiV1Dashboard.GET("/purchase-order/vendor/:id", controllers.GetPurchaseOrder_Vendor)
		ApiV1Dashboard.PUT("/purchase-order/item", controllers.UpdatePurchaseOrder_Item)
		ApiV1Dashboard.GET("/purchase-order/item/:id", controllers.GetPurchaseOrder_Item)
		ApiV1Dashboard.DELETE("/purchase-order/item/:id", controllers.DeletePurchaseOrder)
		ApiV1Dashboard.GET("/purchase-order/printview/:id", controllers.GetPurchaseOrder_PrintView)
		ApiV1Dashboard.GET("/purchase-order/printnow/:id", controllers.GetPurchaseOrder_PrintNow)

		// Sales Order Management
		ApiV1Dashboard.GET("/sales-order", controllers.GetSalesOrder)
		ApiV1Dashboard.POST("/sales-order", controllers.CreateSalesOrder)
		ApiV1Dashboard.GET("/sales-order/customer/:id", controllers.GetSalesOrder_Customer)
		ApiV1Dashboard.GET("/sales-order/item/:id", controllers.GetSalesOrder_Item)
		ApiV1Dashboard.PUT("/sales-order/customer", controllers.UpdateSalesOrder_Customer)
		ApiV1Dashboard.PUT("/sales-order/item", controllers.UpdateSalesOrder_Item)
		ApiV1Dashboard.DELETE("/sales-order/item/:id", controllers.DeleteSalesOrder)
		ApiV1Dashboard.GET("/sales-order/shipping-cost/:id", controllers.GetSalesOrder_ShipCost)
		ApiV1Dashboard.PUT("/sales-order/shipping-cost", controllers.UpdateSalesOrder_ShipCost)
		ApiV1Dashboard.GET("/sales-order/suggest/type", controllers.GetSalesOrder_SuggestType)
		ApiV1Dashboard.GET("/sales-order/suggest/customer", controllers.GetSalesOrder_SuggestCustomer)
		ApiV1Dashboard.GET("/sales-order/suggest/item", controllers.GetSalesOrder_SuggestItem)
		ApiV1Dashboard.GET("/sales-order/suggest/attr", controllers.GetSalesOrder_SuggestAttr)

		// Workorder Management
		ApiV1Dashboard.GET("/workorder", controllers.Workorder_Get)
		ApiV1Dashboard.POST("/workorder", controllers.Workorder_Create)
		ApiV1Dashboard.GET("/workorder/print/:wocusid/:sequence_item", controllers.Workorder_Printview)
		ApiV1Dashboard.POST("/workorder/print", controllers.Workorder_Printnow)
		ApiV1Dashboard.GET("/workorder/process/:id/:sequence_item", controllers.Workorder_GetProcess)

		// Delivery Orders Management
		ApiV1Dashboard.GET("/delivery-order", controllers.GetDeliveryOrder)
		ApiV1Dashboard.GET("/delivery-order/waiting", controllers.GetDeliveryOrder_Waiting)
		ApiV1Dashboard.GET("/delivery-order/item/:id", controllers.GetDeliveryOrder_Item)
		ApiV1Dashboard.GET("/delivery-order/number", controllers.GetDeliveryOrder_Number)
		ApiV1Dashboard.POST("/delivery-order", controllers.CreateDeliveryOrder)
		ApiV1Dashboard.DELETE("/delivery-order/:id", controllers.DeleteDeliveryOrder)
		ApiV1Dashboard.GET("/delivery-order/printview/:id", controllers.GetDeliveryOrder_Printview)
		ApiV1Dashboard.GET("/delivery-order/printnow/:id", controllers.GetDeliveryOrder_Printnow)

		// Invoice Management
		ApiV1Dashboard.GET("/invoice", controllers.GetInvoice)
		ApiV1Dashboard.POST("/invoice", controllers.CreateInvoice)
		ApiV1Dashboard.GET("/invoice/print/:id", controllers.GetInvoice_Printview)
		ApiV1Dashboard.POST("/invoice/print", controllers.GetInvoice_Printnow)
		ApiV1Dashboard.POST("/invoice/paid", controllers.PaidInvoice)
		ApiV1Dashboard.PUT("/invoice/unpaid", controllers.UnpaidInvoice)
		ApiV1Dashboard.DELETE("/invoice/:id", controllers.DeleteInvoice)

	}

}

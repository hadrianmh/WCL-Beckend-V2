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
		ApiV1Dashboard.POST("/purchase-order", controllers.CreatePurchaseOrder)
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
		ApiV1Dashboard.PUT("/sales-order/shipping-cost", controllers.UpdateSalesOrder_Item)

		// Workorder Management
		ApiV1Dashboard.GET("/workorder", controllers.GetWorkorder)
		ApiV1Dashboard.POST("/workorder", controllers.CreateWorkorder)
		ApiV1Dashboard.GET("/workorder/print/:wocusid/:sequence_item", controllers.GetWorkorder_Printview)
		ApiV1Dashboard.POST("/workorder/print", controllers.GetWorkorder_Printnow)

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

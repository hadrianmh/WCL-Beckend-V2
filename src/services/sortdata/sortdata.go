package sortdata

import (
	"backend/adapters"
	"backend/config"
	"backend/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SortDataResponse struct {
	Month []string `json:"month,omitempty"`
	Year  []string `json:"year,omitempty"`
	JmlPo string   `json:"jml_po,omitempty"`
	JmlDo string   `json:"jml_do,omitempty"`
	JmlIn string   `json:"jml_in,omitempty"`
}

func GetArchive(ctx *gin.Context) ([]SortDataResponse, error) {
	Data := ctx.DefaultQuery("data", "")
	From := ctx.DefaultQuery("from", "")

	if From == `` {
		From = `preorder_customer`
	}

	if Data == `` {
		Data = `po_date`
	}

	// Load Config
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return nil, err
	}

	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT DISTINCT %s FROM %s ORDER BY %s DESC`, strings.ToLower(Data), strings.ToLower(From), strings.ToLower(Data))
	fmt.Println(query)
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	archive := []string{}
	month := []string{}
	year := []string{}
	sortdata := SortDataResponse{}
	for rows.Next() {
		var column string
		var columnParse time.Time

		if err := rows.Scan(&column); err != nil {
			return nil, err
		}

		// Generate estimasi +30 hari invoice
		if columnParse, err = time.Parse(config.App.DateFormat_Global, column); err != nil {
			return nil, err
		}

		Month := columnParse.Format(config.App.DateFormat_MonthlyReport)
		Year := columnParse.Format(config.App.DateFormat_Years)

		// Filtering month unique
		if utils.InArray(archive, Month) {
		} else {
			month = append(month, Month)
		}

		// Filtering year unique
		if utils.InArray(archive, Year) {
		} else {
			year = append(year, Year)
		}

		archive = append(archive, Month)
		archive = append(archive, Year)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	sortdata.Month = month
	sortdata.Year = year

	return []SortDataResponse{sortdata}, nil
}

func GetCounter(ctx *gin.Context) ([]SortDataResponse, error) {
	TypeParam := ctx.DefaultQuery("type", "")
	StartdateParam := ctx.DefaultQuery("startdate", "")
	EnddateParam := ctx.DefaultQuery("enddate", "")

	// Preventif type not found , set default must single
	if TypeParam != "single" {
		if TypeParam != "periode" {
			TypeParam = "single"
		}
	}

	// Preventif end date null, set to start date because the param isnt required
	if TypeParam == `periode` && EnddateParam == `` {
		EnddateParam = StartdateParam
	}

	var Querystr string
	var jml_po, jml_do, jml_in string
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	if strings.ToLower(TypeParam) == `periode` {
		Querystr = fmt.Sprintf(`BETWEEN '%s' AND '%s'`, StartdateParam, EnddateParam)
	} else {
		Querystr = fmt.Sprintf(`LIKE '%s%%'`, StartdateParam)
	}

	query := fmt.Sprintf(`SELECT count(id) AS jml_po FROM preorder_customer WHERE po_date %s`, Querystr)
	if err = sql.Connection.QueryRow(query).Scan(&jml_po); err != nil {
		if err.Error() == `sql: no rows in result set` {
			jml_po = `0`
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`SELECT count(a.id) as jml_do FROM delivery_orders_customer AS a LEFT JOIN status AS b ON a.id_fk = b.id_fk WHERE b.order_status >= 6 AND a.sj_date %s GROUP BY a.id_fk`, Querystr)
	if err = sql.Connection.QueryRow(query).Scan(&jml_do); err != nil {
		if err.Error() == `sql: no rows in result set` {
			jml_do = `0`
		} else {
			return nil, err
		}
	}

	query = fmt.Sprintf(`SELECT count(DISTINCT no_invoice) as jml_in FROM invoice WHERE status = 1 AND invoice_date %s`, Querystr)
	if err = sql.Connection.QueryRow(query).Scan(&jml_in); err != nil {
		if err.Error() == `sql: no rows in result set` {
			jml_in = `0`
		} else {
			return nil, err
		}
	}

	sortdata := SortDataResponse{}
	sortdata.JmlDo = jml_do
	sortdata.JmlPo = jml_po
	sortdata.JmlIn = jml_in

	return []SortDataResponse{sortdata}, nil
}

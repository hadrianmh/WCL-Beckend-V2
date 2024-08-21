package setting

import "backend/adapters"

type DataBank struct {
	Id     int    `json:"id,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func Bank() ([]DataBank, error) {
	sql, err := adapters.NewSql()
	if err != nil {
		return nil, err
	}

	query := `SELECT id,isi FROM setting WHERE ket = 'BANK'`
	rows, err := sql.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	databank := []DataBank{}
	for rows.Next() {
		var id int
		var detail string

		if err := rows.Scan(&id, &detail); err != nil {
			return nil, err
		}

		databank = append(databank, DataBank{
			Id:     id,
			Detail: detail,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return databank, nil
}

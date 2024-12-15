package database

import (
	"github.com/uptrace/bun"
)

type Models struct {
	Model interface{}
	Fkey  *string
}

func ModelsToBytes(db *bun.DB, models []Models) []byte {
	var data []byte
	for _, m := range models {
		query := db.NewCreateTable().Model(m.Model).WithForeignKeys()
		if m.Fkey != nil {
			query.ForeignKey(*m.Fkey)
		}
		raw, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, raw...)
		data = append(data, ";\n"...)
	}
	return data
}

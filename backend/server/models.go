package server

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemoDM struct {
	bun.BaseModel `bun:"table:memos"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserId        uuid.UUID `bun:"user_id,type:uuid,notnull"`
	Title         string    `bun:"type:text,notnull"`
	Text          string    `bun:"type:text,notnull"`
}

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

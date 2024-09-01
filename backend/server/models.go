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

func ModelsToBytes(db *bun.DB, models []interface{}) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys()
		raw, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, raw...)
		data = append(data, ";\n"...)
	}
	return data
}

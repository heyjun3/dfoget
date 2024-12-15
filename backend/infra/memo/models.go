package memo

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

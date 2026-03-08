package database

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID        pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string      `json:"name"`
	Username  string      `json:"username"`
	Password  string      `json:"password"`
	Email     string      `json:"email"`
	CreatedAt pgtype.Timestamptz
}

type Tag struct {
	ID     pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId pgtype.UUID `gorm:"type:uuid"`
	Name   string      `json:"name"`
}

type Link struct {
	ID       pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId   pgtype.UUID `gorm:"type:uuid"`
	Link     string      `json:"link"`
	Name     string      `json:"name"`
	LinkDesc string      `json:"link_desc"`

	LinkTags []LinkTags `gorm:"foreignKey:link_id"`
}

type LinkTags struct {
	LinkId pgtype.UUID `gorm:"type:uuid"`
	TagId  pgtype.UUID `gorm:"type:uuid"`
	Tag    Tag         `gorm:"foreignKey:tag_id"`
}

package models

import (
	"server/internal/utils"
)

type ContactInfo struct {
	ClientID     int     `json:"-"`
	Email        *string `json:"email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
}

type Client struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	ContactInfo   *ContactInfo     `json:"contact_info,omitempty"`
	CountProyects *int             `json:"countProyects,omitempty"`
	CreatedAt     utils.CustomTime `json:"created_at"`
	UpdatedAt     utils.CustomTime `json:"updated_at"`
}

package models

import (
	"time"
)

type ContactInfo struct {
	ClientID     int     `json:"-"`
	Email        *string `json:"email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
}

type Client struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	ContactInfo   *ContactInfo `json:"contact_info,omitempty"`
	Countprojects *int         `json:"countprojects,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

package utils

import "time"

type LogActivity struct {
	UserID string `json:"userID"`

	// log activity
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	Params       interface{} `json:"params"`
	QueryParams  interface{} `json:"queryParams"`
	RequestBody  interface{} `json:"requestBody"`
	ResponseData interface{} `json:"responseData"`
	TimeAccessed *time.Time  `json:"timeAccessed"`

	// status
	Status      string `json:"status"`
	ErrorCode   string `json:"errorCode"`
	Message     string `json:"message"`
	Description string `json:"description"`

	UserAgent string `json:"userAgent"`
	IP        string `json:"ip"`
	RequestID string `json:"requestID"`

	Region   string `json:"region"`
	TicketID string `json:"ticketID"`
}

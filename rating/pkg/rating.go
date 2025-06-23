package model

// RecordID defines a record id. Together with Recordtype
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types.
type RecordType string

// Existing record types.
const (
	RecordTypeMovie RecordType = RecordType("movie")
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for
// some record
type Rating struct {
	RecordID   RecordID   `json:"recordID"`
	RecordType RecordType `json:"recordType"`
	UserID     UserID     `json:"userID"`
	Value      RatingValue
}

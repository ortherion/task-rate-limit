package models

type Task struct {
	ID          uint64        `db:"id"          json:"id"`
	CreatorID   uint64        `db:"creator_id"  json:"-"`
	Title       string        `db:"title"       json:"title,omitempty"`
	Body        string        `db:"body"        json:"body,omitempty"`
	IsDeleted   bool          `db:"is_deleted"  json:"-"`
	Stage       Stage         `db:"status_task" json:"status_task,omitempty"`
	Signatories []Signatories `json:"signatories,omitempty"`
	Date        `json:"date,omitempty"`
}

type Stage int

const (
	Undefined Stage = iota
	Accept
	Reject
	InProcess
)

func (s Stage) String() string {
	switch s {
	case Undefined:
		return "undefined"
	case Accept:
		return "accept"
	case Reject:
		return "reject"
	case InProcess:
		return "in process"
	default:
		return "unknown type"
	}
}

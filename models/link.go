package models

type LinkStatusType string

const (
	LinkStatusAvailable    LinkStatusType = "available"
	LinkStatusNotAvailable LinkStatusType = "not_available"
	LinkStatusUnknown      LinkStatusType = "unknown"
)

type Link struct {
	Link   string         `json:"link"`
	Status LinkStatusType `json:"status"`
}

type LinkListStatusType string

const (
	LinkListStatusInProgress LinkListStatusType = "in_progress"
	LinkListStatusDone       LinkListStatusType = "done"
)

type LinkList struct {
	Num    int                `json:"links_num"`
	Links  []Link             `json:"links"`
	Status LinkListStatusType `json:"status"`
}

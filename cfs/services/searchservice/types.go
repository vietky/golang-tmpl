package searchservice

type CFS struct {
	AgencyID      string `json:"agency_id"`
	EventID       string `json:"event_id"`
	EventNumber   string `json:"event_number"`
	EventTypeCode string `json:"event_type_code"`
	EventTime     string `json:"event_time"`
	DispatchTime  string `json:"dispatch_time"`
	Responder     string `json:"responder"`
}

type User struct {
	UserId   int
	AgencyID string
}

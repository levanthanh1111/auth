package model

// WithdrawRequestStatus type def
type RequestStatus uint8

const (
	// StatusUnknown
	StatusUnknown RequestStatus = 0
	// StatusNew
	StatusNew RequestStatus = 1
	// StatusReadToCollect
	StatusReadToCollect RequestStatus = 2
	// StatusCollected
	StatusCollected RequestStatus = 3
	// StatusCancelled
	StatusCancelled RequestStatus = 4
)

var (
	mapStatusString = map[RequestStatus]string{
		StatusNew:           "new",
		StatusReadToCollect: "ready_to_collect",
		StatusCollected:     "collected",
		StatusCancelled:     "cancelled",
	}

	mapStringStatus = map[string]RequestStatus{
		"new":              StatusNew,
		"ready_to_collect": StatusReadToCollect,
		"collected":        StatusCollected,
		"cancelled":        StatusCancelled,
	}
)

func (s RequestStatus) String() string {
	str := mapStatusString[s]
	if str == "" {
		return "unknown"
	}

	return str
}

func WithdrawRequestStatusFromString(s string) RequestStatus {
	return mapStringStatus[s]
}

// WithdrawRequest model
type WithdrawRequest struct {
	ID           uint64 `json:"id"`
	ContractID   uint64 `json:"contract"`
	ContractorID uint64 `json:"project_contractor_id"`
	Quantity     uint64 `json:"quantity"`
}

// WithdrawRequestStatus `json:"status"`
type WithdrawRequestStatus struct {
	WithdrawRequestID uint64        `json:"id"`
	Status            RequestStatus `json:"status"`
}

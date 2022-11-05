package model

type Contract struct {
	ID       uint64 `json:"id,omitempty"`
	VendorID uint64 `json:"-" gorm:"column:supply_vendor_id"`

	StartDate string `json:"start_date,omitempty"`

	EndDate string `json:"end_date,omitempty"`

	BaseAmount int32 `json:"base_amount,omitempty"`

	ActualAmount    int32  `json:"actual_amount,omitempty"`
	Code            string `json:"code,omitempty"`
	SuplyVendor     *Org   `json:"-" gorm:"foreignKey:VendorID;references:ID"`
	SuplyVendorName string `json:"supply_vendor_name,omitempty"`
}

func (c *Contract) UpdateVendorName() {
	if c != nil && c.SuplyVendor != nil {
		c.SuplyVendorName = c.SuplyVendor.Name
	}
}

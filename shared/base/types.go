package base

import (
	"encoding/json"
	"fmt"
)

type Operator string

const (
	OperatorEqual              Operator = "eq"
	OperatorNotEqual           Operator = "ne"
	OperatorLessThan           Operator = "lt"
	OperatorLessThanOrEqual    Operator = "le"
	OperatorGreaterThan        Operator = "gt"
	OperatorGreaterThanOrEqual Operator = "ge"
	OperatorLike               Operator = "like"
)

type Filters = []*Filter

type ListParams struct {
	Limit   int     `form:"limit" schema:"limit" validate:"gt=0"`
	Offset  int     `form:"offset" schema:"offset" validate:"gte=0"`
	Sort    string  `form:"sort" schema:"sort" validate:""`
	Filters Filters `form:"filters" schema:"filters" validate:"dive,required"`

	acceptedFilterKeys []string `form:"-" schema:"-"`
}

type Filter struct {
	Key      string   `form:"key" schema:"key" validate:"required"`
	Operator Operator `form:"operator" schema:"operator" validate:"required,eq=eq|eq=lt|eq=gt|eq=ge|eq=le|eq=ne|eq=like"`
	Value    string   `form:"value" schema:"value" validate:"required"`
}

type Pagination struct {
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Sort      string `json:"sort"`
	Total     int64  `json:"total"`
	TotalPage int64  `json:"total_page"`
	List      any    `json:"list"`
}

func (lp *ListParams) EnsureDefault() {

	if lp.Limit < 10 {
		lp.Limit = 10
	}
	lp.Sort = "id desc"

	if lp.Offset < 0 {
		lp.Offset = 0
	}

}

func (lp *ListParams) RegisterFilterKeys(acceptedKeys []string) {
	lp.acceptedFilterKeys = acceptedKeys
}

func (lp *ListParams) IsValid() error {

	var err error
	if err = v.Struct(lp); err != nil {
		return err
	}

	if len(lp.Filters) == 0 || len(lp.acceptedFilterKeys) == 0 {
		return nil
	}

	for _, filter := range lp.Filters {
		if err = v.Struct(struct {
			AcceptedFilterKeys []string `schema:"-"`
			Key                string   `json:"key" schema:"filters.$i.key" validate:"oneoffield=AcceptedFilterKeys"`
		}{
			Key:                filter.Key,
			AcceptedFilterKeys: lp.acceptedFilterKeys,
		}); err != nil {
			return err
		}
	}

	return nil

}

var mapMsgValidationTag = map[string]func(val any, param string) string{
	"required": func(val any, param string) string {
		return "This field is required"
	},
	"email": func(val any, param string) string {
		return "invalid email format"
	},
	"lt": func(val any, param string) string {
		return fmt.Sprintf("This field value length must be less than %v", param)
	},
	"lte": func(val any, param string) string {
		return fmt.Sprintf("This field value length must be less than or equal to %v", param)
	},
	"gt": func(val any, param string) string {
		return fmt.Sprintf("This field value length must be greater than %v", param)
	},
	"gte": func(val any, param string) string {
		return fmt.Sprintf("This field value length must be greater than or equal to %v", param)
	},
	"eq": func(val any, param string) string {
		return fmt.Sprintf("This field must be equal to '%v'", param)
	},
	"ne": func(val any, param string) string {
		return fmt.Sprintf("This field must not be equal to '%v'", param)
	},
	"oneoffield": func(val any, param string) string {
		return fmt.Sprintf("not support for value '%v'", val)
	},
	"alphanum": func(val any, param string) string {
		return "invalid alphanumeric format"
	},
}

type APIError struct {
	Field   string
	Message string
}

type APIErrors []*APIError

func (aes APIErrors) Error() string {
	bs, _ := json.Marshal(aes)
	return string(bs)
}

func (aes APIErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"errors": []*APIError(aes)})
}

func NewApiErrors(field, msg string) APIErrors {
	return APIErrors{&APIError{Field: field, Message: msg}}
}

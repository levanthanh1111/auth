package base

import "github.com/tpp/msf/shared/validator"

var v = validator.Get()

type Validatee interface {
	IsValid() error
}

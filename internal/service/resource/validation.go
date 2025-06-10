package resource

import (
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
)

var (
	validSchemes = map[string]struct{}{
		string(HTTPSScheme): {},
		string(HTTPScheme):  {},
	}

	validMethods = map[string]struct{}{
		string(GetMethod):    {},
		string(PostMethod):   {},
		string(PutMethod):    {},
		string(DELETEMethod): {},
	}
)

func (dto *ResourceDTO) Validate() *errors.HTTPError {
	if dto.PublicPath == "" {
		return errors.BadRequestError("public path cannot be empty")
	}
	if dto.ServicePath == "" {
		return errors.BadRequestError("service path cannot be empty")
	}
	if _, ok := validSchemes[dto.Scheme]; !ok {
		return errors.BadRequestError("invalid scheme")
	}
	if _, ok := validMethods[dto.Method]; !ok {
		return errors.BadRequestError("invalid scheme")
	}

	return nil
}

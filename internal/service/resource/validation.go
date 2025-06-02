package resource

import (
	"fmt"
)

var (
	validSchemes = map[string]struct{}{
		string(HTTPSScheme): {},
	}

	validMethods = map[string]struct{}{
		string(GetMethod):    {},
		string(PostMethod):   {},
		string(PutMethod):    {},
		string(DELETEMethod): {},
	}
)

func (dto *ResourceDTO) Validate() error {
	if dto.PublicPath == "" {
		return fmt.Errorf("public path cannot be empty")
	}
	if dto.ServicePath == "" {
		return fmt.Errorf("service path cannot be empty")
	}
	if _, ok := validSchemes[dto.Scheme]; !ok {
		return fmt.Errorf("invalid scheme: %s", dto.Scheme)
	}
	if _, ok := validMethods[dto.Method]; !ok {
		return fmt.Errorf("invalid method: %s", dto.Method)
	}
	return nil
}

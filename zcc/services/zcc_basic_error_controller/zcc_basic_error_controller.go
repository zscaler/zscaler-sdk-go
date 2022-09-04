package basic_error_controller

const (
	basicErrorEndpoint = "/error"
)

type BasicErrorController struct {
	Empty     bool                   `json:"empty,omitempty"`
	Model     map[string]interface{} `json:"model,omitempty"`
	ModelMap  map[string]interface{} `json:"modelMap,omitempty"`
	Reference bool                   `json:"reference,omitempty"`
	Status    []string               `json:"status,omitempty"`
	View      []string               `json:"view,omitempty"`
	ViewName  string                 `json:"viewName,omitempty"`
}

func (service *Service) GetAll() ([]BasicErrorController, error) {
	var basicError []BasicErrorController
	err := service.Client.Read(basicErrorEndpoint, &basicError)
	return basicError, err
}

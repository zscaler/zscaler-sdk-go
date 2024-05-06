package get_otp

const (
	getOtpEndpoint = "/public/v1/getOtp"
)

type GetOtp struct {
	Otp string `json:"otp"`
}

type GetOtpQuery struct {
	Udid string `json:"udid,omitempty" url:"udid,omitempty"`
}

func (service *Service) GetOtp(udid string) (*GetOtp, error) {
	var otp GetOtp
	_, err := service.Client.NewRequestDo("GET", getOtpEndpoint, GetOtpQuery{Udid: udid}, nil, &otp)
	if err != nil {
		return nil, err
	}
	return &otp, err
}

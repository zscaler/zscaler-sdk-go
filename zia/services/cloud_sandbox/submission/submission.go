package submission

const (
	submissionEndpoint = "/zscsb/submit"
)

type SandBoxSubmission struct {
	Code              string `json:"code"`
	Message           string `json:"message,omitempty"`
	FileType          string `json:"fileType,omitempty"`
	MD5               string `json:"md5,omitempty"`
	SandBoxSubmission string `json:"sandboxSubmission,omitempty"`
	VirusName         string `json:"virusName,omitempty"`
	VirusType         string `json:"virusType,omitempty"`
}

type ForceQueryParams struct {
	Force string `json:"force,omitempty" url:"force,omitempty"`
}

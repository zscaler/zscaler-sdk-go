package admin_audit_logs

const (
	entryReportEndpoint = "/auditlogEntryReport"
	DownloadEndpoint    = "/auditlogEntryReport/download"
)

type AdminAuditLogs struct {
	Status                string `json:"status,omitempty"`
	ProgressItemsComplete int    `json:"progressItemsComplete,omitempty"`
	ProgressEndTime       int    `json:"progressEndTime,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ErrorCode             string `json:"errorCode,omitempty"`
}

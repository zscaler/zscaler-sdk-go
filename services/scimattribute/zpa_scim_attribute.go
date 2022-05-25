package scimattribute

const (
//	scimAttributeEndpoint = "/idp/{idpId}/scimattribute"
)

type ScimAttribute struct {
	List       []List `json:"list"`
	TotalPages int32  `json:"totalpages"`
}
type List struct {
	CanonicalValues []string `json:"canonicalvalues"`
	CaseSensitive   bool     `json:"caseSensitive"`
	CreationTime    int32    `json:"creationtime"`
	DataType        string   `json:"datatype"`
	Description     string   `json:"description"`
	ID              int64    `json:"id"`
	IdpID           int64    `json:"idpid"`
	ModifiedBy      int64    `json:"modifiedby"`
	ModifiedTime    int32    `json:"modifiedtime"`
	Multivalued     bool     `json:"multivalued"`
	Mutability      string   `json:"mutability"`
	Name            string   `json:"name"`
	Required        bool     `json:"required"`
	Returned        string   `json:"returned"`
	SchemaURI       string   `json:"schemauri"`
	Uniqueness      bool     `json:"uniqueness"`
}

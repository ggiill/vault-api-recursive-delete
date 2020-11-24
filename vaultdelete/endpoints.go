package vaultdelete

type endpointPath struct {
	method string
	path   string
}

var versionPaths = map[string]map[string]*endpointPath{
	"v2": {
		"list": &endpointPath{
			method: "LIST",
			path:   "v1/secret/metadata",
		},
		"get": &endpointPath{
			method: "GET",
			path:   "v1/secret/data",
		},
		"delete": &endpointPath{
			method: "DELETE",
			path:   "v1/secret/data",
		},
		// "destroy": &endpointPath{
		// 	method: "POST",
		// 	path:   "v1/secret/destroy",
		// },
	},
}

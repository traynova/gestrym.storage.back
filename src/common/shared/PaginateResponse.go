package shared

type ResponsePaginate struct {
	// Página actual
	Page int `json:"page"`
	// Tamaño de la página
	PageSize int `json:"pageSize"`
	// Total de elementos
	Total int `json:"total"`
	// Resultados
	Results []interface{} `json:"result"`
}

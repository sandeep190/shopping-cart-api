package dtobjects

type AddCart struct {
	RequestType string `json:"type" binding:"required"`
	ProductId   int    `json:"product_id" binding:"required"`
}

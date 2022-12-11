package dtobjects

type SignupRequestDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"address"`
	Contact  string `json:"contact" binding:"required"`
}

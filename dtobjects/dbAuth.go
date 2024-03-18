package dtobjects

type DBCredentials struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	BbInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

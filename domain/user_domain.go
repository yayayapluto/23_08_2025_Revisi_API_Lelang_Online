package domain

type (
	LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	UpdateRequestUser struct {
		Username *string `json:"username"`
	}
)

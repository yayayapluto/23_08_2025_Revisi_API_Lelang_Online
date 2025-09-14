package domain

type (
	LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	RegisterInput struct {
		Username        string `json:"username"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	CreateOrganizerUserInput struct {
		Username    string
		Email       string
		Password    string
		OrganizerID uint
	}

	UpdateRequestUser struct {
		Username *string `json:"username"`
	}
)

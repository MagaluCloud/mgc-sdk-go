package containerregistry

type (
	Meta struct {
		Page Page `json:"page"`
	}

	Page struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)

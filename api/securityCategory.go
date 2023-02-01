package api

// Security category object
type SecurityCategory struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	SubCategories []SecuritySubCategory `json:"subCategories"`
}

// Security sub-category object
type SecuritySubCategory struct {
	Id                       string `json:"id"`
	Title                    string `json:"title"`
	Description              string `json:"description"`
	ResolutionRecommendation string `json:"resolutionRecommendation"`
}

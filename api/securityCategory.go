package api

// Security category object
type SecurityCategory struct {
	Description   string                `json:"description"`
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	SubCategories []SecuritySubCategory `json:"subCategories"`
}

// Security sub-category object
type SecuritySubCategory struct {
	Description              string `json:"description"`
	Id                       string `json:"id"`
	ResolutionRecommendation string `json:"resolutionRecommendation"`
	Title                    string `json:"title"`
}

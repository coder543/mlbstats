package main

type Teams struct {
	Copyright string `json:"copyright"`
	Teams     []Team `json:"teams"`
}

type Team struct {
	SpringLeague    Identifiable `json:"springLeague"`
	AllStarStatus   string       `json:"allStarStatus"`
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Link            string       `json:"link"`
	Season          int          `json:"season"`
	Venue           Identifiable `json:"venue"`
	SpringVenue     Identifiable `json:"springVenue"`
	TeamCode        string       `json:"teamCode"`
	FileCode        string       `json:"fileCode"`
	Abbreviation    string       `json:"abbreviation"`
	TeamName        string       `json:"teamName"`
	LocationName    string       `json:"locationName"`
	FirstYearOfPlay string       `json:"firstYearOfPlay"`
	League          Identifiable `json:"league"`
	Division        Identifiable `json:"division"`
	Sport           Identifiable `json:"sport"`
	ShortName       string       `json:"shortName"`
	FranchiseName   string       `json:"franchiseName"`
	ClubName        string       `json:"clubName"`
	Active          bool         `json:"active"`
}

package model

type (
	RecordID    string
	RecordType  string
	UserID      string
	RatingValue int
)

const RecordMovieType RecordType = "movie"

type Rating struct {
	RecordID   RecordID    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"ratingValue"`
}

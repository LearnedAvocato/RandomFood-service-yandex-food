package datastruct

import "time"

type RequestLog struct {
	RequestedCardsNum  int64     `db:"requested_cards_num"`
	GotCardsNum        int64     `db:"got_cards_num"`
	Longitude          float64   `db:"longitude"`
	Latitude           float64   `db:"latitude"`
	UsedRestarauntsNum int64     `db:"used_restaraunts_num"`
	CardsPerRestaraunt int64     `db:"cards_per_restaraunt"`
	CreatedAt          time.Time `db:"created_at"`
}

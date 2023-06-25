package repository

import (
	"context"
	"fmt"
	"yandex-food/internal/pkg/datastruct"
)

func (r *Repository) LogRequest(ctx context.Context, log *datastruct.RequestLog) error {

	query := fmt.Sprintf("insert into %s (requested_cards_num, got_cards_num, longitude, latitude, used_restaraunts_num, cards_per_restaraunt, created_at) values ($1, $2, $3, $4, $5, $6, $7)",
		RequestLogsTableName)

	_, err := r.client.Exec(ctx, query,
		log.RequestedCardsNum, log.GotCardsNum, log.Longitude, log.Latitude, log.UsedRestarauntsNum, log.CardsPerRestaraunt, log.CreatedAt)

	return err
}

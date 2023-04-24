package bot

import (
	"context"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/dto"
)

func (api *Api) getUpdates(request telegram.GetUpdatesRequest, ctx context.Context) (updates []dto.Update, err error) {
	response, err := api.RequestContext(request, ctx)
	defer func() {
		if v := recover(); v != nil {
			updates = nil
			err = nil
		}
	}()
	if err != nil {

		if ctx.Err() != nil {
			return nil, nil
		}

		return nil, err
	}
	return telegram.ParseAs[[]dto.Update](response)
}

func (api *Api) Poll(ctx context.Context) <-chan dto.Update {
	updateC := make(chan dto.Update, 100)
	request := telegram.GetUpdatesRequest{
		Offset:         0,
		Limit:          100,
		Timeout:        10,
		AllowedUpdates: []string{"message", "callback_query"},
	}

	go func() {
		defer close(updateC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := api.getUpdates(request, ctx)
				if err != nil {
					continue
				}
				for _, update := range updates {
					if update.UpdateID >= request.Offset {
						updateC <- update
						request.Offset = update.UpdateID + 1
					}
				}
			}
		}
	}()

	return updateC
}

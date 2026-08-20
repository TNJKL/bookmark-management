package queue

import (
	"context"
	"encoding/json"

	"github.com/TNJKL/bookmark-management/pkg/batch"
)

const BatchSize = 20

func (s *service) SendImportBookmarkJob(ctx context.Context, uid string, bookmarkInputs []*ImportBookmarkInput) error {
	//split array into batch
	batches := batch.SplitIntoBatches(bookmarkInputs, BatchSize)
	for _, batch := range batches {
		err := s.sendJob(ctx, uid, batch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) sendJob(ctx context.Context, uid string, bookmarkInputs []*ImportBookmarkInput) error {
	//create ImportMessage struct
	message := ImportMessage{
		UID:       uid,
		Bookmarks: bookmarkInputs,
	}

	//marshal ImportMessage struct to json
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	//push message to redis queue
	return s.q.PushMessage(ctx, messageBytes)
}

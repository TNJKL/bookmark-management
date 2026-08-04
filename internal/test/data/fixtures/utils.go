package fixtures

import (
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/model"
)

var (
	TestTime = time.Date(2018, 2, 18, 1, 2, 3, 4, time.UTC)
)

// GetTestBase returs a test base
func GetTestBase(id string) model.Base {
	return model.Base{
		ID:        id,
		CreatedAt: TestTime,
		UpdatedAt: TestTime,
	}
}

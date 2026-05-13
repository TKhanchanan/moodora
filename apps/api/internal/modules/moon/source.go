package moon

import (
	"context"
	"fmt"
	"net/http"
)

type APODClient interface {
	GetAPOD(ctx context.Context, date string) (APOD, error)
}

type NASAAPODClient struct {
	HTTPClient *http.Client
	APIKey     string
	BaseURL    string
}

func (c NASAAPODClient) GetAPOD(ctx context.Context, date string) (APOD, error) {
	return APOD{}, fmt.Errorf("nasa apod adapter is not enabled yet")
}

func InternalCalculationSource() Source {
	return Source{
		SourceCode: "internal_moon_calculation",
		SourceName: "Moodora Internal Moon Calculation",
		SourceURL:  "internal://moon-phase-v1",
	}
}

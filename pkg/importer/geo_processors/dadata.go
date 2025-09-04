package geo_processors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ekomobile/dadata/v2/api/clean"

	"github.com/EveryHotel/fsa-gov/pkg/importer"

	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/client"
)

const GeoProcessorNameDadata = "dadata"

type dadataProcessor struct {
	client *clean.Api
}

func NewDadataProcessor(
	apiKey string,
	apiSecret string,
) importer.GeoProcessor {
	return &dadataProcessor{
		client: dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
			ApiKeyValue:    apiKey,
			SecretKeyValue: apiSecret,
		})),
	}
}

func (p *dadataProcessor) GetName() string {
	return GeoProcessorNameDadata
}

func (p *dadataProcessor) Process(ctx context.Context, address string) (city string, latitude float64, longitude float64, data json.RawMessage, err error) {
	addresses, err := p.client.Address(ctx, address)
	if err != nil {
		return city, latitude, longitude, data, fmt.Errorf("clean adddress: %w", err)
	}

	if len(addresses) == 0 {
		return city, latitude, longitude, data, fmt.Errorf("couldn't parse address: %s", address)
	}

	data, err = json.Marshal(addresses)
	if err != nil {
		return city, latitude, longitude, data, fmt.Errorf("pack geo data: %w", err)
	}

	// Не всегда определен город, по-этом путаемся достать хотя бы какое-нибудь название чего-нибудь
	city = addresses[0].City
	if city == "" {
		city = addresses[0].Settlement
		if city == "" {
			city = addresses[0].Region
			if city == "" {
				// Если не смогли определить город - это не так критично как координаты, по-этому не возвращаем ошибку, а просто пишем в лог
				slog.WarnContext(ctx, "city was not found",
					slog.String("address", address),
				)
			}
		}
	}

	if addresses[0].GeoLat != "" {
		latitude, err = strconv.ParseFloat(addresses[0].GeoLat, 64)
		if err != nil {
			return city, latitude, longitude, data, fmt.Errorf("couldn't parse latitude: %w", err)
		}
	}

	if addresses[0].GeoLon != "" {
		longitude, err = strconv.ParseFloat(addresses[0].GeoLon, 64)
		if err != nil {
			return city, latitude, longitude, data, fmt.Errorf("couldn't parse longitude: %w", err)
		}
	}

	return city, latitude, longitude, data, nil
}

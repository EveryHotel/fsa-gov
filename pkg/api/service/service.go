package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/EveryHotel/fsa-gov/pkg/api"
	"github.com/EveryHotel/fsa-gov/pkg/api/dto"
)

type GetNamedItemsFunc func(context.Context) ([]dto.NamedItem, error)

type ApiService interface {
	GetResorts(context.Context, string) (dto.ResortListResponse, error)
	GetResort(context.Context, string) (dto.ResortViewResponse, error)
	GetRegions(context.Context) ([]dto.NamedItem, error)
	GetRoomCategories(context.Context) ([]dto.NamedItem, error)
	GetHotelStatuses(context.Context) ([]dto.NamedItem, error)
	GetHotelCategories(context.Context) ([]dto.NamedItem, error)
	GetAccrAreas(context.Context) ([]dto.NamedItem, error)
	GetHotelTypes(context.Context) ([]dto.NamedItem, error)
}

type apiService struct {
	client api.ApiClient
}

func NewApiService(client api.ApiClient) ApiService {
	return &apiService{
		client: client,
	}
}

func (s apiService) GetResorts(ctx context.Context, date string) (res dto.ResortListResponse, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/export/resorts/changes?start=%s", date), nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetResort(ctx context.Context, code string) (res dto.ResortViewResponse, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/export/resorts/%s/get", code), nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetRegions(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/regions/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetRoomCategories(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/roomCategory/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetHotelStatuses(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/hotelStatus/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetHotelCategories(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/hotelCategory/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetAccrAreas(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/accrArea/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

func (s apiService) GetHotelTypes(ctx context.Context) (res []dto.NamedItem, err error) {
	response, err := s.client.MakeRequest(ctx, "GET", "/export/hotelTypes/get", nil)
	if err != nil {
		return res, fmt.Errorf("make request: %w", err)
	}

	if err = json.Unmarshal(response, &res); err != nil {
		return res, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}

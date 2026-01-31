package transaction

import (
	"backer/campaign"
	"errors"
)

// Custom errors
var (
	ErrCampaignNotFound = errors.New("campaign not found")
	ErrNotAuthorized    = errors.New("not authorized")
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.ID == 0 {
		return []Transaction{}, ErrCampaignNotFound
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, ErrNotAuthorized
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

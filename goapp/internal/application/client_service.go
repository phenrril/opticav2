package application

import (
	"errors"

	"opticav2/internal/domain"
)

type ClientService struct {
	ClientRepo domain.ClientRepository
	// UserRepo domain.UserRepository // If needed for validating UserID exists
}

func NewClientService(clientRepo domain.ClientRepository) *ClientService {
	return &ClientService{ClientRepo: clientRepo}
}

func (s *ClientService) CreateClient(req domain.ClientCreateRequest, userID int) (*domain.Client, error) {
	// Check if DNI already exists
	if req.DNI != "" {
		_, err := s.ClientRepo.FindByDNI(req.DNI)
		if err == nil { // err is nil means DNI found
			return nil, domain.ErrClientDNITaken
		}
		if !errors.Is(err, domain.ErrClientNotFound) { // An unexpected error occurred
			return nil, err
		}
	} else { // DNI is required as per struct tags, but defensive check here
		return nil, errors.New("DNI is required")
	}

	// Optionally, check if name already exists if DNI is not a strict unique identifier for some reason
	// For this example, DNI is the primary unique business identifier.

	client := &domain.Client{
		Name:       req.Name,
		Phone:      req.Phone,
		Address:    req.Address,
		DNI:        req.DNI,
		ObraSocial: req.ObraSocial,
		Medico:     req.Medico,
		UserID:     userID, // ID of the user creating the client
		Status:     1,      // Default to active
		HC:         req.HC,
	}

	if err := s.ClientRepo.Create(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) GetClient(id int) (*domain.Client, error) {
	client, err := s.ClientRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return nil, domain.ErrClientNotFound
		}
		return nil, err // Other unexpected error
	}
	return client, nil
}

func (s *ClientService) ListClients() ([]domain.Client, error) {
	clients, err := s.ClientRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (s *ClientService) UpdateClient(id int, req domain.ClientUpdateRequest) (*domain.Client, error) {
	client, err := s.ClientRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}

	// If DNI is being changed, check for conflicts
	if req.DNI != "" && req.DNI != client.DNI {
		existingClient, errDNI := s.ClientRepo.FindByDNI(req.DNI)
		if errDNI == nil && existingClient.ID != id { // DNI exists for another client
			return nil, domain.ErrClientDNITaken
		}
		if !errors.Is(errDNI, domain.ErrClientNotFound) && errDNI != nil { // Unexpected error
			return nil, errDNI
		}
		client.DNI = req.DNI
	}

	// Update fields
	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.ObraSocial != "" {
		client.ObraSocial = req.ObraSocial
	}
	if req.Medico != "" {
		client.Medico = req.Medico
	}
	if req.HC != "" {
		client.HC = req.HC
	}
	if req.Status != nil {
		client.Status = *req.Status
	}

	if errUpdate := s.ClientRepo.Update(client); errUpdate != nil {
		return nil, errUpdate
	}
	return client, nil
}

func (s *ClientService) DeactivateClient(id int) error {
	client, err := s.ClientRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return domain.ErrClientNotFound
		}
		return err
	}
	client.Status = 0 // Assuming 0 means inactive
	return s.ClientRepo.Update(client)
}

func (s *ClientService) ActivateClient(id int) error {
	client, err := s.ClientRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return domain.ErrClientNotFound
		}
		return err
	}
	client.Status = 1 // Assuming 1 means active
	return s.ClientRepo.Update(client)
}

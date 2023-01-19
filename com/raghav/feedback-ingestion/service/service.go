package service

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	name        string
	retriever   Retriever
	transformer Transformer
}

type Retriever interface {
	RetrieveData(startTime, endTime time.Time) ([]byte, error)
}

type Transformer interface {
	TransformData(data []byte) ([]byte, error)
}

func NewService(name string, retriever Retriever, transformer Transformer) *Service {
	return &Service{
		name:        name,
		retriever:   retriever,
		transformer: transformer,
	}
}

func (s *Service) ProcessFeedback(ctx context.Context, startTime, endTime time.Time) ([]byte, error) {
	data, err := s.retriever.RetrieveData(startTime, endTime)
	if err != nil {
		return nil, err
	}

	transformedData, err := s.transformer.TransformData(data)
	if err != nil {
		return nil, err
	}

	return transformedData, nil
}

func (s *Service) Name() string {
	return s.name
}

type ServiceRegistry struct {
	services map[string]*Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*Service),
	}
}

func (r *ServiceRegistry) RegisterService(service *Service) error {
	if _, ok := r.services[service.Name()]; ok {
		return errors.New("Service already registered")
	}
	r.services[service.Name()] = service
	return nil
}

func (r *ServiceRegistry) GetService(name string) (*Service, error) {
	service, ok := r.services[name]
	if !ok {
		return nil, errors.New("Service not found")
	}
	return service, nil
}

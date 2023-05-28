package service

import "github.com/atomicai/whoosh/internal/repository"

type IDijkstraService interface {
}

type DijkstraService struct {
	repository *repository.DijkstraRepository
}

func NewDijkstraService(dbname string) *DijkstraService {
	return &DijkstraService{repository: repository.NewDijkstraRepository(dbname)}
}

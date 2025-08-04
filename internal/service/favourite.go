package service

import "bookshelf/internal/repository"

type FavouriteService interface {
	AddFavourite(userID, bookID uint) error
	RemoveFavourite(userID, bookID uint) error
	GetFavourites(userID uint, page, limit int) ([]BookBrief, int64, error)
}

type favouriteService struct {
	repo repository.FavouriteRepository
}

func NewFavouriteRepository(repo repository.FavouriteRepository) FavouriteService {
	return &favouriteService{repo: repo}
}

func (s *favouriteService) AddFavourite(userID, bookID uint) error {
	return s.repo.AddFavourite(userID, bookID)
}

func (s *favouriteService) RemoveFavourite(userID, bookID uint) error {
	return s.repo.RemoveFavourite(userID, bookID)
}

func (s *favouriteService) GetFavourites(userID uint, page, limit int) ([]BookBrief, int64, error) {
	books, total, err := s.repo.GetFavourites(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	briefs := make([]BookBrief, len(books))
	for i, book := range books {
		briefs[i] = BookBrief{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Genre:  book.Genre,
			Price:  book.Price,
		}
	}
	return briefs, total, nil
}

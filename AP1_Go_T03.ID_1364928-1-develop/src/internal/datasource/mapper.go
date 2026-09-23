package datasource

import "tic-tac-toe/internal/domain"

// распаковывает игру из DTO в Domain из хранилища (между слоями)
func ToDomain(dto *GameDTO) *domain.Game {
	return &domain.Game{
		ID:    dto.ID,
		Board: dto.Board,
	}
}

// упаковывает игру из Domain в DTO для хранилища (между слоями)
func ToDTO(game *domain.Game) *GameDTO {
	return &GameDTO{
		ID:    game.ID,
		Board: game.Board,
	}
}

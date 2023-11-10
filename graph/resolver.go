package graph

import "funovation_23/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	VideoUsecase *usecases.VideoUsecase

	Mapper      *Mapper
	InputMapper *InputMapper
}

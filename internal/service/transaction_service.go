package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
	"github.com/alanzhumalin/bank/pkg/pagination"
)

type transactionService struct {
	repo     repository.TransactionRepository
	userRepo repository.UserRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (ts *transactionService) GetByUserId(ctx context.Context, userId int, cursorValue string, limit int) (dto.CursorResponse[dto.TransactionResponse], error) {

	queryLimit := limit + 1

	cursor, err := pagination.DecodeTransactionCursor(cursorValue)
	if err != nil {
		return dto.CursorResponse[dto.TransactionResponse]{}, err
	}

	trs, err := ts.repo.GetByUserId(ctx, userId, cursor, queryLimit)

	if err != nil {
		return dto.CursorResponse[dto.TransactionResponse]{}, err
	}

	hasnext := len(trs) > limit

	if hasnext {
		trs = trs[:limit]
	}

	nextCursor := ""

	if hasnext {
		newNextCursor, err := pagination.EncodeTransactionCursor(pagination.TransactionCursor{
			Id:        trs[len(trs)-1].Id,
			CreatedAt: trs[len(trs)-1].CreatedAt,
		})

		if err != nil {
			return dto.CursorResponse[dto.TransactionResponse]{}, err
		}

		nextCursor = newNextCursor
	}

	sl := make([]dto.TransactionResponse, 0, len(trs))

	for _, val := range trs {
		sl = append(sl, dto.ToTransactionResponse(val))
	}

	res := dto.CursorResponse[dto.TransactionResponse]{
		Data: sl,
		Meta: dto.CursorMeta{
			HasNext:    hasnext,
			Limit:      limit,
			NextCursor: nextCursor,
		},
	}

	return res, nil

}

func (ts *transactionService) GetByAccountId(ctx context.Context, accountId int, limit int, cursorValue string, currentUserId int) (dto.CursorResponse[dto.TransactionResponse], error) {

	cursor, err := pagination.DecodeTransactionCursor(cursorValue)
	if err != nil {
		return dto.CursorResponse[dto.TransactionResponse]{}, err
	}

	queryLimit := limit + 1

	trs, ownerUserId, err := ts.repo.GetByAccountId(ctx, accountId, queryLimit, cursor)

	if err != nil {
		return dto.CursorResponse[dto.TransactionResponse]{}, err
	}

	if currentUserId != ownerUserId {
		return dto.CursorResponse[dto.TransactionResponse]{}, domain.ErrorForBidden
	}

	hasNext := len(trs) > limit

	if hasNext {
		trs = trs[:limit]
	}

	nextCursor := ""

	if hasNext && len(trs) > 0 {
		newNextCursor, err := pagination.EncodeTransactionCursor(pagination.TransactionCursor{
			Id:        trs[len(trs)-1].Id,
			CreatedAt: trs[len(trs)-1].CreatedAt,
		})

		if err != nil {
			return dto.CursorResponse[dto.TransactionResponse]{}, err
		}
		nextCursor = newNextCursor
	}

	sl := make([]dto.TransactionResponse, 0, len(trs))

	for _, val := range trs {
		sl = append(sl, dto.ToTransactionResponse(val))
	}

	res := dto.CursorResponse[dto.TransactionResponse]{
		Data: sl,
		Meta: dto.CursorMeta{
			Limit:      limit,
			HasNext:    hasNext,
			NextCursor: nextCursor,
		},
	}

	return res, nil
}

func (ts *transactionService) GetAll(ctx context.Context) ([]dto.TransactionResponse, error) {
	trs, err := ts.repo.GetAll(ctx)

	if err != nil {
		return []dto.TransactionResponse{}, err
	}

	sl := make([]dto.TransactionResponse, 0, len(trs))

	for _, val := range trs {
		sl = append(sl, dto.ToTransactionResponse(val))
	}

	return sl, nil
}

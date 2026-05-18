package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type TransactionCursor struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func EncodeTransactionCursor(cursor TransactionCursor) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeTransactionCursor(cursor string) (*TransactionCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var tS TransactionCursor

	if err = json.Unmarshal(payload, &tS); err != nil {
		return nil, err
	}

	if tS.CreatedAt.IsZero() {
		return nil, errors.New("invalid cursor")
	}

	if tS.Id <= 0 {
		return nil, errors.New("invalid cursor")
	}

	return &tS, err

}

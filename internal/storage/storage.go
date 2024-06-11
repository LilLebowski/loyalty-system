package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/LilLebowski/loyalty-system/internal/db"
	"github.com/LilLebowski/loyalty-system/internal/utils"
)

type DBStorage struct {
	db *sql.DB
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   string `json:"userID,omitempty"`
}

type User struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	ID           string `json:"id,omitempty"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type DBOrder struct {
	Number     string
	Status     string
	Accrual    sql.NullFloat64
	UploadedAt time.Time
}

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type Storage interface {
	Register(ctx context.Context, login string, passwordHash string) (string, error)
	GetUserByLogin(ctx context.Context, authData Auth) (User, error)
	AddOrderForUser(ctx context.Context, externalOrderID string, userID string) error
	GetOrdersByUser(ctx context.Context, userID string) ([]Order, error)
	GetUserBalance(ctx context.Context, userID string) (UserBalance, error)
	AddWithdrawalForUser(ctx context.Context, userID string, withdrawal Withdrawal) error
	GetWithdrawalsForUser(ctx context.Context, userID string) ([]Withdrawal, error)
	GetOrdersInProgress(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order Accrual) error
}

func Init(databasePath string) Storage {
	database, err := db.Init(databasePath)
	if err != nil {
		log.Printf("error while starting db %s", err.Error())
		return nil
	}
	return &DBStorage{db: database}
}

func (s *DBStorage) Register(ctx context.Context, login string, passwordHash string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id FROM \"user\" WHERE \"login\" = $1", login)
	var userID sql.NullString
	err := row.Scan(&userID)
	if err != nil && userID.Valid {
		return "", err
	}
	if userID.Valid {
		return "", fmt.Errorf("user with login %s exist", login)
	}
	row = s.db.QueryRowContext(ctx, "INSERT INTO \"user\" (\"login\", password_hash) VALUES ($1, $2) RETURNING id", login, passwordHash)
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	if userID.Valid {
		userIDValue := userID.String
		return userIDValue, nil
	}
	return "", err
}

func (s *DBStorage) GetUserByLogin(ctx context.Context, auth Auth) (User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM \"user\" WHERE login = $1", auth.Login)
	var userData User
	err := row.Scan(&userData.ID, &userData.Login, &userData.PasswordHash)
	if err != nil {
		return userData, err
	}
	return userData, nil
}

func (s *DBStorage) GetOrdersByUser(ctx context.Context, userID string) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT number, status, accrual, uploaded_at FROM \"order\" WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orderList := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal DBOrder
		err := rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual, &orderFromDBVal.UploadedAt)
		if err != nil {
			return nil, err
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status, UploadedAt: orderFromDBVal.UploadedAt}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orderList = append(orderList, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orderList, nil
}

func (s *DBStorage) AddOrderForUser(ctx context.Context, number string, userID string) error {
	row := s.db.QueryRowContext(ctx, "SELECT user_id FROM \"order\" WHERE number = $1", number)
	var orderUserID sql.NullString
	err := row.Scan(&orderUserID)
	if err != nil && orderUserID.Valid {
		return err
	}
	if orderUserID.Valid {
		if orderUserID.String == userID {
			return utils.NewOrderIsExistThisUserError("this order is exist the user", err)
		} else {
			return utils.NewOrderIsExistAnotherUserError("this order is exist another user", err)
		}
	}
	row = s.db.QueryRowContext(
		ctx,
		"INSERT INTO \"order\" (user_id, status, number) VALUES ($1, $2, $3) RETURNING number",
		userID, "NEW", number,
	)
	var orderID string
	err = row.Scan(&orderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetUserBalance(ctx context.Context, userID string) (UserBalance, error) {
	sumOrdersRow := s.db.QueryRowContext(ctx, "SELECT sum(accrual) FROM \"order\" WHERE user_id = $1", userID)
	sumWithdrawalsRow := s.db.QueryRowContext(ctx, "SELECT sum(sum) FROM withdrawal WHERE user_id = $1", userID)
	var sumOrders sql.NullFloat64
	var sumWithdrawals sql.NullFloat64
	err := sumOrdersRow.Scan(&sumOrders)
	if err != nil && sumOrders.Valid {
		return UserBalance{0, 0}, err
	}
	var resultBalance UserBalance
	if !sumOrders.Valid {
		resultBalance.Current = 0
	} else {
		resultBalance.Current = sumOrders.Float64
	}
	err = sumWithdrawalsRow.Scan(&sumWithdrawals)

	if err != nil && sumWithdrawals.Valid {
		return UserBalance{0, 0}, err
	}
	if !sumWithdrawals.Valid {
		resultBalance.Withdrawn = 0
	} else {
		resultBalance.Withdrawn = sumWithdrawals.Float64
	}
	resultBalance.Current -= resultBalance.Withdrawn
	return resultBalance, nil
}

func (s *DBStorage) AddWithdrawalForUser(ctx context.Context, userID string, withdrawal Withdrawal) error {
	userBalance, err := s.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if userBalance.Current < withdrawal.Sum {
		return utils.NewLessBonusErrorError("Got less bonus points than expected", err)
	}
	var withdrawalID string
	row := s.db.QueryRowContext(ctx,
		"INSERT INTO withdrawal (user_id, sum, external_order_id) VALUES ($1, $2, $3) RETURNING id",
		userID, withdrawal.Sum, withdrawal.Order,
	)
	err = row.Scan(&withdrawalID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetWithdrawalsForUser(ctx context.Context, userID string) ([]Withdrawal, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT external_order_id, sum, processed_at FROM withdrawal WHERE user_id = $1", userID)
	if err != nil {
		return make([]Withdrawal, 0), err
	}
	defer rows.Close()
	withdrawalList := make([]Withdrawal, 0)
	for rows.Next() {
		var withdrawal Withdrawal
		err = rows.Scan(&withdrawal.Sum, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return make([]Withdrawal, 0), err
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return withdrawalList, nil
}

func (s *DBStorage) GetOrdersInProgress(ctx context.Context) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT number, status, accrual from \"order\" where status not in ('INVALID', 'PROCESSED')")

	if err != nil {
		return make([]Order, 0), err
	}
	defer rows.Close()
	orderList := make([]Order, 0)
	for rows.Next() {
		var orderFromDBVal DBOrder
		err = rows.Scan(&orderFromDBVal.Number, &orderFromDBVal.Status, &orderFromDBVal.Accrual)
		if err != nil {
			return make([]Order, 0), err
		}
		order := Order{Number: orderFromDBVal.Number, Status: orderFromDBVal.Status}
		if orderFromDBVal.Accrual.Valid {
			order.Accrual = orderFromDBVal.Accrual.Float64
		}
		orderList = append(orderList, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orderList, nil
}

func (s *DBStorage) UpdateOrder(ctx context.Context, order Accrual) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE \"order\" SET status = $1, accrual = $2 where number = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(order.Status, order.Accrual, order.Order); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/database/SQLC"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/helper"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/helper/validator"
)

type OrderRepository struct {
	Db      *sql.DB
	Queries *SQLC.Queries
}

func NewOrderRepository(db *sql.DB, q *SQLC.Queries) *OrderRepository {
	return &OrderRepository{Db: db, Queries: q}
}

func (r *OrderRepository) List(ctx context.Context, page, limit int, sort string) ([]entities.Order, error) {
	query := "SELECT id, price, tax, final_price, issue_date, type_requisition, delete_at FROM orders WHERE delete_at IS NULL"
	args := []interface{}{}

	if sort == "asc" || sort == "desc" {
		query += " ORDER BY id " + sort
	}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query += " OFFSET ?"
		args = append(args, offset)
	}

	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	for rows.Next() {
		var i entities.Order
		var issueDateBytes []byte
		if err := rows.Scan(
			&i.ID,
			&i.Price,
			&i.Tax,
			&i.FinalPrice,
			&issueDateBytes,
			&i.TypeRequisition,
			&i.DeleteAt,
		); err != nil {
			return nil, err
		}

		if issueDateBytes != nil {
			issueDateStr := string(issueDateBytes)
			i.IssueDate, err = time.Parse("2006-01-02 15:04:05", issueDateStr)
			if err != nil {
				return nil, err
			}
		}
		orders = append(orders, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) Save(ctx context.Context, arg *entities.Order) error {
	err := validator.ValidateOrder(arg)
	if err != nil {
		return err
	}
	var deleteAt sql.NullTime
	if arg.DeleteAt != nil {
		deleteAt = sql.NullTime{Time: *arg.DeleteAt, Valid: !arg.DeleteAt.IsZero()}
	} else {
		deleteAt = sql.NullTime{Valid: false}
	}

	err = r.Queries.Save(ctx, SQLC.SaveParams{
		ID:              helper.NewID().String(),
		Price:           arg.Price,
		Tax:             arg.Tax,
		FinalPrice:      arg.FinalPrice,
		IssueDate:       sql.NullTime{Time: arg.IssueDate, Valid: !arg.IssueDate.IsZero()},
		TypeRequisition: arg.TypeRequisition,
		DeleteAt:        deleteAt,
	})
	return err
}

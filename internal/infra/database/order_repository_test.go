package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/database/SQLC"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB
var repository *OrderRepository

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de testes: %v", err)
	}

	createTable := `
	CREATE TABLE orders (
		id TEXT NOT NULL PRIMARY KEY,
		price REAL NOT NULL,
		tax REAL NOT NULL,
		final_price REAL NOT NULL,
		issue_date DATETIME NOT NULL,
		type_requisition VARCHAR(10) NOT NULL,
		delete_at DATETIME NULL
	);`
	_, err = testDB.Exec(createTable)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	repository = NewOrderRepository(testDB, SQLC.New(testDB))

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func TestSaveIntegration(t *testing.T) {
	order := entities.Order{
		Price:           150.75,
		Tax:             15.00,
		FinalPrice:      165.75,
		IssueDate:       time.Now(),
		TypeRequisition: "rest",
		DeleteAt:        nil,
	}

	err := repository.Save(context.Background(), &order)
	assert.NoError(t, err)

	var uuid string
	err = testDB.QueryRow("SELECT id FROM orders").Scan(&uuid)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuid)

	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", uuid).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	_, err = testDB.Exec("DELETE FROM orders")
	assert.NoError(t, err)
}

func TestListOrders_PaginationAndSorting(t *testing.T) {
	orders := []entities.Order{
		{ID: "1", Price: 100.0, Tax: 10.0, FinalPrice: 110.0, IssueDate: time.Now(), TypeRequisition: "rest", DeleteAt: nil},
		{ID: "2", Price: 200.0, Tax: 20.0, FinalPrice: 220.0, IssueDate: time.Now(), TypeRequisition: "rest", DeleteAt: nil},
		{ID: "3", Price: 300.0, Tax: 30.0, FinalPrice: 330.0, IssueDate: time.Now(), TypeRequisition: "rest", DeleteAt: nil},
		{ID: "4", Price: 400.0, Tax: 40.0, FinalPrice: 440.0, IssueDate: time.Now(), TypeRequisition: "rest", DeleteAt: nil},
		{ID: "5", Price: 500.0, Tax: 50.0, FinalPrice: 550.0, IssueDate: time.Now(), TypeRequisition: "rest", DeleteAt: nil},
	}

	for _, order := range orders {
		_, err := testDB.Exec(`INSERT INTO orders (id, price, tax, final_price, issue_date, type_requisition, delete_at) 
				VALUES (?,?,?,?,?,?,?)`, order.ID, order.Price, order.Tax, order.FinalPrice, order.IssueDate, order.TypeRequisition, order.DeleteAt)
		assert.NoError(t, err)
	}

	// Testar limit = 3, página 1, ordenação ascendente
	results, err := repository.List(context.Background(), 1, 3, "asc")
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, "1", results[0].ID)
	assert.Equal(t, "2", results[1].ID)
	assert.Equal(t, "3", results[2].ID)
	assert.Panics(t, func() {
		_ = results[3].ID
	}, "Panic: index out of range")

	// Testar limit = 2, página 2, ordenação descendente
	results, err = repository.List(context.Background(), 2, 2, "desc")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "3", results[0].ID)
	assert.Equal(t, "2", results[1].ID)

	// Testar limit = 5 (todos os registros), ordenação descendente
	results, err = repository.List(context.Background(), 1, 5, "desc")
	assert.NoError(t, err)
	assert.Len(t, results, 5)
	assert.Equal(t, "5", results[0].ID)
	assert.Equal(t, "4", results[1].ID)
	assert.Equal(t, "3", results[2].ID)
	assert.Equal(t, "2", results[3].ID)
	assert.Equal(t, "1", results[4].ID)

	_, err = testDB.Exec("DELETE FROM orders")
	assert.NoError(t, err)
}

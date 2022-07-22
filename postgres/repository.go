package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"store/product"
	"store/store"
	"store/user"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) ListUsers(ctx context.Context) ([]user.UserCard, error) {
	ucSlice := []user.UserCard{}
	rows, err := r.db.Query(`
		SELECT 
		u.full_name, 
		u.password, 
		u.email,
		c.card_number,
		c.balance FROM users AS u 
		JOIN card AS c ON u.card_id=c.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: %w", err)
	}
	
	for rows.Next(){
		uc := user.UserCard{}
		err := rows.Scan(
			&uc.FullName,
			&uc.Password,
			&uc.Email,
			&uc.CardNumber,
			&uc.Balance,
		)
		if err != nil {
			return nil, fmt.Errorf("ListUser (rows.scan): %w", err)
		}
		ucSlice = append(ucSlice, uc)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ListUser (rows.Err): %w", err)
	}
	defer rows.Close()
	return ucSlice, nil
}

func (r *PostgresRepository) AddCard(ctx context.Context, c user.Card) error {
	_, err := r.db.Exec(`
		INSERT INTO card 
		(id, card_number, balance)
		VALUES ($1, $2, $3)
	`, c.ID, c.CardNumber, c.Balance,
	)

	if err != nil {
		return fmt.Errorf("AddCard: %w", err)
	}

	return nil
}

func (r *PostgresRepository) AddUser(ctx context.Context, u user.User, c user.Card) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("AddCard: %w", err)
	}
	_, err = tx.Exec(`
	INSERT INTO card 
	(id, card_number, balance)
	VALUES ($1, $2, $3)
	`, c.ID, c.CardNumber, c.Balance,
	)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("AddCard: %w", err)
	}

	_, err = tx.Exec(`
	INSERT INTO users 
	(id, full_name, password, email, card_id)
	VALUES ($1, $2, $3, $4, $5)
	`, u.CardID,
		u.FullName,
		u.Password,
		u.Email,
		u.CardID,
	)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("AddUser: %w", err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetUser(ctx context.Context, fn, email, pw string) (user.User, error) {
	u := user.User{}
	err := r.db.QueryRow(`
		SELECT * FROM users WHERE full_name=$1 AND email=$2 AND password=$3
	`, fn, email, pw).Scan(&u.ID, &u.FullName, &u.Password, &u.Email, &u.CardID, &u.IsAdmin)

	if err != nil {
		return u, fmt.Errorf("GetUser: %w", err)
	}

	return u, nil
}

func (r *PostgresRepository) AddProduct(ctx context.Context, p product.Product) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO product (id, name, quantity, price, original_price)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, p.ID, p.Name, p.Description, p.Quantity, p.Price, p.OriginalPrice)

	if err != nil {
		return fmt.Errorf("AddProduct: %w", err)
	} 
	return nil
}

func (r *PostgresRepository) AddProducts(ctx context.Context, ps []product.Product) error {
	for _, p := range ps {
		_, err := r.db.ExecContext(ctx, `
		INSERT INTO product (id, name, quantity, price, original_price)
		VALUES ($1, $2, $3, $4, $5, $6)
		`, p.ID, p.Name, p.Description, p.Quantity, p.Price, p.OriginalPrice)

		if err != nil {
			continue
		}
	}

	return nil
}

func (r *PostgresRepository) GetProduct(ctx context.Context, name string) (product.Product, error) {
	p := product.Product{}
	err := r.db.QueryRow(`
		SELECT * FROM product WHERE name=$1 
	`, name).Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.OriginalPrice)
	if err != nil {
		return p, fmt.Errorf("GetProduct: %w", err)
	}

	return p, nil
}

func (r *PostgresRepository) ListProducts(ctx context.Context) ([]product.Product, error) {
	products := []product.Product{}
	rows, err := r.db.QueryContext(ctx, `
		SELECT * FROM product
	`)
	if err != nil {
		return nil, fmt.Errorf("ListProduct: %w", err)
	}

	for rows.Next() {
		p := product.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.OriginalPrice)

		if err != nil {
			return nil, fmt.Errorf("ListProduct: %w", err)
		}
		products = append(products, p)
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ListProduct: %w", err)
	}

	return products, nil
}

func (r *PostgresRepository) SellProduct(ctx context.Context, sale store.Sales, product product.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE product SET quantity = quantity-$1 WHERE id=$2
	`, sale.SoldQuantity, sale.ProductID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SellProduct(1): %w", err)
	}

	customer := user.User{}
	err = tx.QueryRowContext(ctx, `
		SELECT * FROM users WHERE id=$1
	`, sale.CustomerID).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.Password,
		&customer.Email,
		&customer.CardID,
		&customer.IsAdmin,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SellProduct(2): %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE card SET balance = balance - $1 WHERE id=$2
	`, sale.Profit, customer.CardID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SellProduct(3): %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO sales_history VALUES ($1, $2, $3, $4, $5)
	`, sale.CustomerID, sale.ProductID,
		sale.SoldQuantity, sale.Profit,
		sale.Time,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SellProduct(4): %w", err)
	}

	tx.Commit()
	return nil
}

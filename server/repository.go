package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NeverlandMJ/arzon-market/product"
	"github.com/NeverlandMJ/arzon-market/store"
	"github.com/NeverlandMJ/arzon-market/user"

	"golang.org/x/crypto/bcrypt"
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
		JOIN card AS c ON u.id=c.owner
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
		(id, card_number, balance, owner)
		VALUES ($1, $2, $3, $4)
	`, c.ID, c.CardNumber, c.Balance, c.OwnerID,
	)

	if err != nil {
		return fmt.Errorf("AddCard: %w", err)
	}
	return nil
}

func (r *PostgresRepository) AddUser(ctx context.Context, u user.User) error {
	bp, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
	INSERT INTO users 
	(id, full_name, password, email)
	VALUES ($1, $2, $3, $4)
	`, u.ID,
		u.FullName,
		string(bp),
		u.Email,
	)

	if err != nil {
		return fmt.Errorf("AddUser: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetUser(ctx context.Context, email, pw string) (user.User, error) {
	
	u := user.User{}
	err := r.db.QueryRow(`
		SELECT * FROM users WHERE email=$1
	`, email).Scan(&u.ID, &u.FullName, &u.Password, &u.Email, &u.IsAdmin)

	if err != nil {
		return user.User{}, fmt.Errorf("GetUser: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (r *PostgresRepository) AddProduct(ctx context.Context, p product.Product) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO product (id, name, description, quantity,  price, original_price, img)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, p.ID, p.Name, p.Description, p.Quantity, p.Price, p.OriginalPrice, p.ImageLink)

	if err != nil {
		return fmt.Errorf("AddProduct: %w", err)
	} 
	return nil
}

func (r *PostgresRepository) AddProducts(ctx context.Context, ps []product.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddProduct: %w", err)
	}
	for _, p := range ps {
		_, err := tx.ExecContext(ctx, `
		INSERT INTO product (id, name, description, quantity, price, original_price, img)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, p.ID, p.Name, p.Description, p.Quantity, p.Price, p.OriginalPrice, p.ImageLink)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("AddProduct: %w", err)
		}
	}
	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetProduct(ctx context.Context, id string) (product.Product, error) {
	p := product.Product{}
	err := r.db.QueryRow(`
		SELECT * FROM product WHERE id=$1 
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.OriginalPrice, &p.ImageLink)
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
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.OriginalPrice, &p.ImageLink)

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
		&customer.IsAdmin,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SellProduct(2): %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE card SET balance = balance - $1 WHERE owner=$2
	`, sale.Profit, customer.ID)

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


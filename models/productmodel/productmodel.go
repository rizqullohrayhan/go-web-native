package productmodel

import (
	"go-web-native/config"
	"go-web-native/entities"
)

func GetAll() []entities.Product {
	rows, err := config.DB.Query(`
		SELECT
			products.id AS product_id,
			products.name AS product_name,
			products.category_id AS product_category_id,
			products.image AS product_image,
			products.stock AS product_stock,
			products.description AS product_description,
			products.created_at AS product_created_at,
			products.updated_at AS product_updated_at,
			categories.name AS category_name
		FROM products
		LEFT JOIN categories ON products.category_id = categories.id
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.CategoryId,
			&product.Image,
			&product.Stock,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category.Name,
		)
		if err != nil {
			panic(err)
		}
		products = append(products, product)
	}
	return products
}

func GetById(id string) entities.Product {
	rows := config.DB.QueryRow(`
		SELECT
			products.id AS product_id,
			products.name AS product_name,
			products.category_id AS product_category_id,
			products.image AS product_image,
			products.stock AS product_stock,
			products.description AS product_description,
			products.created_at AS product_created_at,
			products.updated_at AS product_updated_at,
			categories.name AS category_name
		FROM products
		LEFT JOIN categories ON products.category_id = categories.id
		WHERE products.id=?
	`, id)

	var product entities.Product
	err := rows.Scan(
		&product.Id,
		&product.Name,
		&product.CategoryId,
		&product.Image,
		&product.Stock,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Category.Name,
	)
	if err != nil {
		panic(err)
	}
	return product
}

func Create(product entities.Product) error {
	_, err := config.DB.Exec(
		"INSERT INTO products (name, category_id, image, stock, description) VALUES (?, ?, ?, ?, ?)",
		product.Name, product.CategoryId, product.Image, product.Stock, product.Description,
	)
	if err != nil {
		return err
	}
	return nil
}

func Edit(product entities.Product, id string) error {
	_, err := config.DB.Exec(
		"UPDATE products SET name=?, category_id=?, image=?, stock=?, description=? WHERE id=?",
		product.Name, product.CategoryId, product.Image, product.Stock, product.Description, id,
	)
	if err != nil {
		return err
	}
	return nil
}

func Delete(id string) error {
	_, err := config.DB.Exec(
		"DELETE FROM products WHERE id=?",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
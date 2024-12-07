package categorymodel

import (
	"go-web-native/config"
	"go-web-native/entities"
)

func GetAll() []entities.Category {
	rows, err := config.DB.Query(`SELECT * FROM categories`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}

	return categories
}

func GetByID(id string) entities.Category {
	var category entities.Category
	row := config.DB.QueryRow("Select * FROM categories WHERE id=?", id)
	if err := row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
		panic(err)
	}
	return category
}

func Create(name string) error {
	_, err := config.DB.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		return err
	}
	return nil
}

func Edit(name string, id string) error {
	_, err := config.DB.Exec("UPDATE categories SET name=? WHERE id=?", name, id)
	if err != nil {
		return err
	}
	return nil
}

func Delete(id string) error {
	_, err := config.DB.Exec("DELETE FROM categories WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}
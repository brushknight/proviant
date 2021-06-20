package category

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Id int `json:"id"`
	Title string `json:"title"`
}

type DTO struct {
	Id int `json:"id"`
	Title string `json:"title"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) Category{

	p := &Category{}

	r.db.Connection().First(p, "id = ?", id)

	return *p
}

func (r *Repository) GetAll() []Category{

	var categories []Category
	r.db.Connection().Find(&categories)

	return categories
}

func (r *Repository) Delete(id int){

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Delete(&Category{}, id)
}

func (r *Repository) Create(dto DTO){

	model := Category{
		Title: dto.Title,
	}

	r.db.Connection().Create(&model)
}

func (r *Repository) Update(id int, dto DTO){

	model := Category{
		Title: dto.Title,
	}

	r.db.Connection().Model(&Category{Id: id}).Updates(model)
}

func ModelToDTO(m Category) DTO {
	return DTO{
		Id: m.Id,
		Title: m.Title,
	}
}


func (r *Repository) Migrate() error{
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Category{})
	if err != nil{
		return fmt.Errorf("migration of Product table failed: %v", err)
	}
	return nil
}

func Setup(d db.DB) (*Repository, error) {

	repo := &Repository{}

	repo.db = d

	err := repo.Migrate()
	if err != nil{
		return nil, err
	}

	return repo, nil

}

package woocommerce

import (
	"net/http"
)

// Category service
type CategoriesService service

// Category object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#product-category-properties
type Category struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
	Parent int `json:"parent,omitempty"`
	Description string `json:"description,omitempty"`
	Display string `json:"display,omitempty"`
	Image *CategoryImage `json:"image,omitempty"`
	MenuOrder int `json:"menu_order,omitempty"`
	ProductCount int `json:"count,omitempty"`
}

// Category image object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#product-category-image-properties
type CategoryImage struct {
	Id int `json:"id,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	DateCreatedGmt string `json:"date_created_gmt,omitempty"`
	DateModified string `json:"date_modified,omitempty"`
	DateModifiedGmt string `json:"date_modified_gmt,omitempty"`
	Source string `json:"src,omitempty"`
	Name string `json:"name,omitempty"`
	Alt string `json:"alt,omitempty"`
}

type ListCategoryParams struct {
	Context string `url:"context,omitempty"`
	Page int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
	Search string `url:"search,omitempty"`
	Exclude *[]int `url:"exclude,omitempty"`
	Include *[]int `url:"include,omitempty"`
	Order string `url:"order,omitempty"`
	OrderBy string `url:"order_by,omitempty"`
	HideEmpty bool `url:"hide_empty,omitempty"`
	Parent int `url:"parent,omitempty"`
	Product int `url:"product,omitempty"`
	Slug string `url:"slug,omitempty"`
}

type DeleteCategoryParams struct {
	Force string `json:"force,omitempty"`
}

type BatchCategoryUpdate struct {
	Create *[]Category `json:"create,omitempty"`
	Update *[]Category `json:"update,omitempty"`
	Delete *[]int      `json:"delete,omitempty"`
}

type BatchCategoryUpdateResponse struct {
	Create *[]Category `json:"create,omitempty"`
	Update *[]Category `json:"update,omitempty"`
	Delete *[]Category `json:"delete,omitempty"`
}

// Create a category. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#create-a-product-category
func (service *CategoriesService) Create(category *Category) (*Category, *http.Response, error) {
	_url := "/products/categories"
	req, _ := service.client.NewRequest("POST", _url, nil, category)

	createdCategory := new(Category)
	response, err := service.client.Do(req, createdCategory)

	if err != nil {
		return nil, response, err
	}

	return createdCategory, response, nil
}

// Get a category. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#retrieve-a-product-category
func (service *CategoriesService) Get(categoryID string) (*Category, *http.Response, error) {
	_url := "/products/categories/" + categoryID
	req, _ := service.client.NewRequest("GET", _url, nil, nil)

	category := new(Category)
	response, err := service.client.Do(req, category)

	if err != nil {
		return nil, response, err
	}

	return category, response, nil
}

// List categories. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#list-all-product-categories
func (service *CategoriesService) List(opts *ListCategoryParams) (*[]Category, *http.Response, error) {
	_url := "/products/categories"
	req, _ := service.client.NewRequest("GET", _url, opts, nil)

	categories := new([]Category)
	response, err := service.client.Do(req, categories)

	if err != nil {
		return nil, response, err
	}

	return categories, response, nil
}

// Update a category. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#update-a-product-category
func (service *CategoriesService) Update(categoryID string, category *Category) (*Category, *http.Response, error) {
	_url := "/products/categories/" + categoryID
	req, _ := service.client.NewRequest("PUT", _url, nil, category)

	updatedCategory := new(Category)
	response, err := service.client.Do(req, updatedCategory)

	if err != nil {
		return nil, response, err
	}

	return updatedCategory, response, nil
}

//Delete a category. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#delete-a-product-category
func (service *CategoriesService) Delete(categoryID string, opts *DeleteCategoryParams) (*Category, *http.Response, error) {
	_url := "/products/categories/" + categoryID
	req, _ := service.client.NewRequest("DELETE", _url, opts, nil)

	category := new(Category)
	response, err := service.client.Do(req, category)

	if err != nil {
		return nil, response, err
	}

	return category, response, nil
}

// Batch update categories. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#batch-update-product-categories
func (service *CategoriesService) Batch(opts *BatchCategoryUpdate) (*BatchCategoryUpdateResponse, *http.Response, error) {
	_url := "/products/categories/batch"
	req, _ := service.client.NewRequest("POST", _url, nil, opts)

	categories := new(BatchCategoryUpdateResponse)
	response, err := service.client.Do(req, categories)

	if err != nil {
		return nil, response, err
	}

	return categories, response, err
}

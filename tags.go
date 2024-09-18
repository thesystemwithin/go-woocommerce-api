package woocommerce

import (
	"net/http"
)

// Tag service
type TagsService service

// Tag object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#product-tag-properties
type Tag struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	ProductCount int `json:"count,omitempty"`
}

type ListTagParams struct {
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

type DeleteTagParams struct {
	Force string `json:"force,omitempty"`
}

type BatchTagUpdate struct {
	Create *[]Tag `json:"create,omitempty"`
	Update *[]Tag `json:"update,omitempty"`
	Delete *[]int      `json:"delete,omitempty"`
}

type BatchTagUpdateResponse struct {
	Create *[]Tag `json:"create,omitempty"`
	Update *[]Tag `json:"update,omitempty"`
	Delete *[]Tag `json:"delete,omitempty"`
}

// Create a tag. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#create-a-product-tag
func (service *TagsService) Create(tag *Tag) (*Tag, *http.Response, error) {
	_url := "/products/tags"
	req, _ := service.client.NewRequest("POST", _url, nil, tag)

	createdTag := new(Tag)
	response, err := service.client.Do(req, createdTag)

	if err != nil {
		return nil, response, err
	}

	return createdTag, response, nil
}

// Get a tag. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#retrieve-a-product-tag
func (service *TagsService) Get(tagID string) (*Tag, *http.Response, error) {
	_url := "/products/tags/" + tagID
	req, _ := service.client.NewRequest("GET", _url, nil, nil)

	tag := new(Tag)
	response, err := service.client.Do(req, tag)

	if err != nil {
		return nil, response, err
	}

	return tag, response, nil
}

// List tags. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#list-all-product-tags
func (service *TagsService) List(opts *ListTagParams) (*[]Tag, *http.Response, error) {
	_url := "/products/tags"
	req, _ := service.client.NewRequest("GET", _url, opts, nil)

	tags := new([]Tag)
	response, err := service.client.Do(req, tags)

	if err != nil {
		return nil, response, err
	}

	return tags, response, nil
}

// Update a tag. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#update-a-product-tag
func (service *TagsService) Update(tagID string, tag *Tag) (*Tag, *http.Response, error) {
	_url := "/products/tags/" + tagID
	req, _ := service.client.NewRequest("PUT", _url, nil, tag)

	updatedTag := new(Tag)
	response, err := service.client.Do(req, updatedTag)

	if err != nil {
		return nil, response, err
	}

	return updatedTag, response, nil
}

//Delete a tag. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#delete-a-product-tag
func (service *TagsService) Delete(tagID string, opts *DeleteTagParams) (*Tag, *http.Response, error) {
	_url := "/products/tags/" + tagID
	req, _ := service.client.NewRequest("DELETE", _url, opts, nil)

	tag := new(Tag)
	response, err := service.client.Do(req, tag)

	if err != nil {
		return nil, response, err
	}

	return tag, response, nil
}

// Batch update tags. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/?shell#batch-update-product-tags
func (service *TagsService) Batch(opts *BatchTagUpdate) (*BatchTagUpdateResponse, *http.Response, error) {
	_url := "/products/tags/batch"
	req, _ := service.client.NewRequest("POST", _url, nil, opts)

	tags := new(BatchTagUpdateResponse)
	response, err := service.client.Do(req, tags)

	if err != nil {
		return nil, response, err
	}

	return tags, response, err
}

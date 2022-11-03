package models

import (
	"errors"
	"m/v2/utils"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty" query:"limit"`
	Page       int         `json:"page,omitempty" query:"page"`
	Sort       string      `json:"sort,omitempty"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id asc" // TODO: fix, this only works if model has ID
	}
	return p.Sort
}

// Sets Sort using json tab name validated against the supplied model
// It is recommended that a reponse model be used as a whitelist.
func (p *Pagination) SetSort(tag string, desc bool, model interface{}) error {
	if len(tag) == 0 {
		p.Sort = ""
		return nil
	}
	var sd string
	sf := utils.FieldForJsonTag(tag, model)
	if sf == nil {
		return errors.New("Sort field does not exist")
	}
	sfs := utils.ToSnakeCase(*sf)
	if desc {
		sd = sfs + " desc"
	} else {
		sd = sfs + " asc"
	}
	p.Sort = sd
	return nil
}

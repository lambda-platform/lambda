package utils

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	ShowSQL bool
}

// Paginator
type Paginator struct {
	//To int         `json:"to"`
	Total       int64       `json:"total"`
	LastPage    int         `json:"last_page"`
	Data        interface{} `json:"data"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	CurrentPage int         `json:"current_page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

type gridResult struct {
	CurrentPage string      `json:"current_page"`
	Data        interface{} `json:"data"`
	//Data map[string]map[string]string `json:"data"`
	FirstPageUrl string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageUrl  string `json:"last_page_url"`
	NextPageUrl  string `json:"next_page_url"`
	Path         string `json:"path"`
	PerPage      string `json:"per_page"`
	PrevPageUrl  string `json:"prev_page_url"`
	To           int    `json:"to"`
	Total        int    `json:"total"`
}

// Paging
func Paging(p *Param, result interface{}) *Paginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	var paginator Paginator
	var count int64
	var offset int

	countRecords(db, result, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)

	paginator.Total = count
	paginator.Data = result
	paginator.CurrentPage = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.LastPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.LastPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, count *int64) {
	db.Model(anyType).Count(count)
}

func ExtractAllowedSorts(model interface{}) map[string]bool {
	allowed := make(map[string]bool)

	t := reflect.TypeOf(model)

	// Unwrap pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// If it's a slice, use its element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	// We only handle structs
	if t.Kind() != reflect.Struct {
		return allowed
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("gorm")
		if strings.Contains(tag, "column:") {
			parts := strings.Split(tag, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "column:") {
					col := strings.TrimPrefix(part, "column:")
					col = strings.TrimSpace(col)
					if col != "" {
						allowed[col] = true
					}
				}
			}
		}
	}

	return allowed
}

func Sort(model interface{}, sortCol, sortDir string, db *gorm.DB) *gorm.DB {

	if sortCol == "" || sortCol == "null" || sortCol == "undefined" {
		return db
	}

	allowedSorts := ExtractAllowedSorts(model)
	if !allowedSorts[sortCol] {

		return db
	}

	sd := strings.ToLower(sortDir)
	if sd != "asc" && sd != "desc" {
		sd = "asc"
	}

	orderExpr := fmt.Sprintf("%s %s", sortCol, sd)
	return db.Order(orderExpr)
}

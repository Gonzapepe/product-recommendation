package entities

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Product struct {
	ID          string               `json:"id" bson:"id,omitempty"`
	StoreID     string               `json:"storeId,omitempty" bson:"storeId"`
	Categories  []CategoryForProduct `json:"categories,omitempty" bson:"categories"`
	Description Description          `json:"description,omitempty" bson:"description"`
	Images      []Image              `json:"images,omitempty" bson:"images"`
	Name        Name                 `json:"name,omitempty" bson:"name" validate:"required"`
	Published   bool                 `json:"published,omitempty" bson:"published"`
	Urls        Urls                 `json:"urls,omitempty" bson:"urls"`
	Variants    []Variant            `json:"variants,omitempty" bson:"variants"`
	SoldCount   int                  `json:"soldCount,omitempty" bson:"soldCount" validate:"gte=0"`
	ClickCount  int                  `json:"clickCoçunt,omitempty" bson:"clickCount" validate:"gte=0"`
	CreatedAt   time.Time            `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type CategoryForProduct struct {
	ID            string       `json:"id" bson:"id"`
	Name          CategoryName `json:"name" bson:"name"`
	Subcategories []string     `json:"subcategories" bson:"subcategories"`
}

type CategoryName struct {
	LocalizedString `bson:",inline"`
}

type Description struct {
	LocalizedString `bson:",inline"`
}

type Image struct {
	ID       int    `json:"id" bson:"id"`
	Src      string `json:"src" bson:"src"`
	Position int    `json:"position" bson:"position"`
	Alt      []Alt  `json:"alt" bson:"alt"`
}

type Name struct {
	LocalizedString `bson:",inline"`
}

type Urls struct {
	CanonicalURL string  `json:"canonicalURL" bson:"canonicalURL"`
	VideoURL     *string `json:"videoURL" bson:"videoURL"`
}

type Variant struct {
	ID    string  `json:"id" bson:"id"`
	Value string  `json:"value" bson:"value"`
	Stock int     `json:"stock" bson:"stock"`
	Price float64 `json:"price" bson:"price"`
}

type LocalizedString struct {
	En *string `json:"en,omitempty" bson:"en,omitempty"`
	Es *string `json:"es,omitempty" bson:"es,omitempty"`
	Pt *string `json:"pt,omitempty" bson:"pt,omitempty"`
}

type Alt struct {
	LocalizedString `bson:",inline"`
}

type CustomTime struct {
	time.Time
}

// MarshalJSON converts the CustomTime to JSON format.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(time.RFC3339Nano))), nil
}

// UnmarshalJSON converts JSON format to CustomTime.
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"2006-01-02T15:04:05.000Z"`, string(data))
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// UnmarshalBSON converts BSON format to CustomTime.
func (ct *CustomTime) UnmarshalBSON(data []byte) error {
	// BSON Date is int64 representing milliseconds since Unix epoch
	var t time.Time
	if err := bson.Unmarshal(data, &t); err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// MarshalBSON converts CustomTime to BSON format.
func (ct CustomTime) MarshalBSON() ([]byte, error) {
	return bson.Marshal(ct.Time)
}



type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag string `json:"tag"`
	Value string `json:"value"`
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

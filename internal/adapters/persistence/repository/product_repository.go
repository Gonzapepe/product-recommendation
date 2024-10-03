package repository

import (
	"backend-challenge/internal/domain/entities"
	"backend-challenge/internal/domain/repositories"
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Client, dbName, collectionName string) repositories.ProductRepository {
    return &productRepository{
        collection: db.Database(dbName).Collection(collectionName),
    }
}

func (r *productRepository) GetByID(id string) (*entities.Product, error) {
    var product entities.Product
    err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&product)

    if err != nil {
        return nil, err
    }

    return &product, nil
}

func (r *productRepository) GetPaginated(offset, limit int) ([]*entities.Product, error) {
    var products []*entities.Product

    findOptions := options.Find()
    findOptions.SetSkip(int64(offset))
    findOptions.SetLimit(int64(limit))
    
    cursor, err := r.collection.Find(context.TODO(), bson.M{}, findOptions)

    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var product entities.Product
        err := cursor.Decode(&product)
        if err != nil {
            return nil, err
        }
        products = append(products, &product)
    }
    
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return products, nil
}

func (r *productRepository) GetAll() ([]*entities.Product, error) {
    var products []*entities.Product
    
    cursor, err := r.collection.Find(context.TODO(), bson.M{})

    if err != nil {
        return nil, err
    }

    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var elem *entities.Product
        err := cursor.Decode(&elem)
        if err != nil {
            return nil, err
        }

        products = append(products, elem)
    }

    return products, nil
}

func (r *productRepository) Create(product *entities.Product) error {
    _, err := r.collection.InsertOne(context.TODO(), product)
    return err
}

func (r *productRepository) Update(product *entities.Product) error {
    updateData := bson.M{}

    // Helper function to add non-zero values to updateData
    addIfNotZero := func(key string, value interface{}) {
        if !reflect.ValueOf(value).IsZero() {
            updateData[key] = value
        }
    }

    // Check and add each field individually
    addIfNotZero("storeId", product.StoreID)
    addIfNotZero("categories", product.Categories)
    addIfNotZero("description", product.Description)
    addIfNotZero("images", product.Images)
    addIfNotZero("name", product.Name)
    addIfNotZero("published", product.Published)
    addIfNotZero("urls", product.Urls)
    addIfNotZero("variants", product.Variants)
    addIfNotZero("soldCount", product.SoldCount)
    addIfNotZero("clickCount", product.ClickCount)

    updateData["updatedAt"] = time.Now()

    // Only perform the update if there are fields to update
    if len(updateData) > 1 { // > 1 because updatedAt is always present
        _, err := r.collection.UpdateOne(context.TODO(), bson.M{"id": product.ID}, bson.M{"$set": updateData})
        return err
    }

    return nil
}

func (r *productRepository) Delete(id string) error {
    _, err := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
    return err
}
package routes

import (
	"context"
	"main/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemsHandler struct {
	Collection *mongo.Collection
}

func NewItemsHandler(client *mongo.Client) *ItemsHandler {
	collection := client.Database("fantaisie").Collection("items")
	return &ItemsHandler{Collection: collection}
}

func (h *ItemsHandler) CreateItem(c echo.Context) error {
	item := new(models.Item)

	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err := h.Collection.InsertOne(context.TODO(), item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *ItemsHandler) GetItem(c echo.Context) error {
	idParam := c.Param("id")

	if idParam != "" {
		var item models.Item
		err := h.Collection.FindOne(context.TODO(), bson.M{"_id": idParam}).Decode(&item)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, "Item not found")
			}
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		switch item.Type {
		case "weapon":
			var weaponData models.WeaponData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &weaponData)
			item.Data = weaponData
		case "food":
			var foodData models.FoodData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &foodData)
			item.Data = foodData
		case "material":
			var materialData models.MaterialData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &materialData)
			item.Data = materialData
		}

		return c.JSON(http.StatusOK, item)
	}

	cursor, err := h.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve items")
	}
	defer cursor.Close(context.TODO())

	var items []models.Item
	if err := cursor.All(context.TODO(), &items); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse items")
	}

	for i := range items {
		switch items[i].Type {
		case "weapon":
			var weaponData models.WeaponData
			dataBytes, _ := bson.Marshal(items[i].Data)
			bson.Unmarshal(dataBytes, &weaponData)
			items[i].Data = weaponData
		case "food":
			var foodData models.FoodData
			dataBytes, _ := bson.Marshal(items[i].Data)
			bson.Unmarshal(dataBytes, &foodData)
			items[i].Data = foodData
		case "material":
			var materialData models.MaterialData
			dataBytes, _ := bson.Marshal(items[i].Data)
			bson.Unmarshal(dataBytes, &materialData)
			items[i].Data = materialData
		}
	}

	return c.JSON(http.StatusOK, items)
}

func (h *ItemsHandler) UpdateItem(c echo.Context) error {
	id := c.Param("id")

	var item models.Item
	err := h.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, "Item not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = h.Collection.UpdateOne(context.TODO(), bson.M{"_id": id}, map[string]interface{}{"$set": item})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ItemsHandler) DeleteItem(c echo.Context) error {
	id := c.Param("id")

	_, err := h.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Item deleted successfully")
}

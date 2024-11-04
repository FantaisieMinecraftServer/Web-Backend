package routes

import (
	"context"
	"main/models"
	"net/http"
	"strings"

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

func (h *ItemsHandler) GetItems(c echo.Context) error {
	typesParam := c.QueryParam("types")

	filter := bson.M{}
	if typesParam != "" {
		types := strings.Split(typesParam, ",")
		filter["type"] = bson.M{"$in": types}
	}

	cursor, err := h.Collection.Find(context.TODO(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve items")
	}
	defer cursor.Close(context.TODO())

	var items []models.Item
	if err := cursor.All(context.TODO(), &items); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse items")
	}

	categoryMap := make(map[string]models.Category)
	for _, item := range items {
		switch item.Type {
		case models.CategoryWeapon:
			var weaponData models.WeaponData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &weaponData)
			item.Data = weaponData

			groupID := weaponData.Group

			if _, ok := categoryMap[groupID]; !ok {
				categoryMap[groupID] = models.Category{
					ID:   groupID,
					Name: models.WeaponsGroup[groupID],
					DisplayItem: models.DisplayItem{
						ID:              "red_dye",
						Name:            "武器",
						CustomModelData: 1000,
					},
					Items: []models.Item{},
				}
			}

			category := categoryMap[groupID]
			category.Items = append(category.Items, item)
			categoryMap[groupID] = category

		case models.CategoryFood:
			var foodData models.FoodData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &foodData)
			item.Data = foodData

			groupID := "food"

			if _, ok := categoryMap[groupID]; !ok {
				categoryMap[groupID] = models.Category{
					ID:   groupID,
					Name: "Food",
					DisplayItem: models.DisplayItem{
						ID:              "minecraft",
						CustomModelData: 1001,
					},
					Items: []models.Item{},
				}
			}

			category := categoryMap[groupID]
			category.Items = append(category.Items, item)
			categoryMap[groupID] = category

		case models.CategoryMaterial:
			var materialData models.MaterialData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &materialData)
			item.Data = materialData

			groupID := "material"

			if _, ok := categoryMap[groupID]; !ok {
				categoryMap[groupID] = models.Category{
					ID:   groupID,
					Name: "Material",
					DisplayItem: models.DisplayItem{
						ID:              "minecraft",
						CustomModelData: 1002,
					},
					Items: []models.Item{},
				}
			}

			category := categoryMap[groupID]
			category.Items = append(category.Items, item)
			categoryMap[groupID] = category
		}
	}

	orderedGroups := []string{
		models.Dagger,
		models.Sword,
		models.Spear,
		models.Hammer,
		models.Wand,
		models.Bow,
	}

	var data []models.Category
	for _, groupID := range orderedGroups {
		if category, ok := categoryMap[groupID]; ok {
			data = append(data, category)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *ItemsHandler) GetItem(c echo.Context) error {
	idParam := c.Param("id")

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

	updatedItem := new(models.Item)
	if err := c.Bind(updatedItem); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Ensure the correct Data type based on item.Type
	switch item.Type {
	case "weapon":
		var weaponData models.WeaponData
		dataBytes, _ := bson.Marshal(updatedItem.Data)
		bson.Unmarshal(dataBytes, &weaponData)
		updatedItem.Data = weaponData
	case "food":
		var foodData models.FoodData
		dataBytes, _ := bson.Marshal(updatedItem.Data)
		bson.Unmarshal(dataBytes, &foodData)
		updatedItem.Data = foodData
	case "material":
		var materialData models.MaterialData
		dataBytes, _ := bson.Marshal(updatedItem.Data)
		bson.Unmarshal(dataBytes, &materialData)
		updatedItem.Data = materialData
	}

	_, err = h.Collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updatedItem})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedItem)
}

func (h *ItemsHandler) DeleteItem(c echo.Context) error {
	id := c.Param("id")

	_, err := h.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Item deleted successfully")
}

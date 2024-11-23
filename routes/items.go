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

func (h *ItemsHandler) GetItems(c echo.Context) error {
	typeParam := c.QueryParam("type")

	filter := bson.M{}
	if typeParam != "" {
		filter["type"] = typeParam
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

	classMap := make(map[string]models.Class)
	var materials []models.Item

	for _, item := range items {
		switch item.Type {
		case models.TypeWeapon:
			var weaponData models.WeaponData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &weaponData)
			item.Data = weaponData

			groupID := weaponData.Group
			class, exists := classMap[groupID]
			if !exists {
				class = models.Class{
					ID:   groupID,
					Name: models.WeaponClasses[groupID],
					DisplayItem: &models.DisplayItem{
						ID:              "barrier",
						Name:            models.WeaponClasses[groupID],
						CustomModelData: models.WeaponClassesCMD[groupID],
					},
					Items: []models.Item{},
				}
			}

			class.Items = append(class.Items, item)
			classMap[groupID] = class

		case models.TypeMaterial:
			var materialData models.MaterialData
			dataBytes, _ := bson.Marshal(item.Data)
			bson.Unmarshal(dataBytes, &materialData)
			item.Data = materialData

			materials = append(materials, item)
		}
	}

	orderedWeaponClasses := []string{
		models.Dagger,
		models.Sword,
		models.Spear,
		models.Hammer,
		models.Wand,
		models.Bow,
	}

	var weapons []models.Class
	for _, groupID := range orderedWeaponClasses {
		if category, ok := classMap[groupID]; ok {
			weapons = append(weapons, category)
		}
	}

	types := []models.Type{
		{
			ID:   "weapons",
			Name: "武器",
			DisplayItem: &models.DisplayItem{
				ID:              "red_dye",
				Name:            "武器",
				CustomModelData: 10000,
			},
			Classes: weapons,
		},
		{
			ID:          "materials",
			Name:        "素材",
			DisplayItem: nil,
			Classes: []models.Class{
				{
					ID:          "material",
					Name:        "素材",
					DisplayItem: nil,
					Items:       materials,
				},
			},
		},
	}

	if typeParam == "weapon" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": weapons,
		})
	}
	if typeParam == "material" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": materials,
		})
	}

	return c.JSON(http.StatusOK, models.API{
		Status: 200,
		Types:  types,
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

package handler

import (
	"fmt"
	"linksaver/server/database"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type PaginationQueryString struct {
	Page    int
	Limit   int
	Keyword string
}

type TagData struct {
	Tag string `json:"tag"`
}

func HandleCreateTag(gc *gin.Context) {
	var newTagName database.Tag
	var tagData TagData
	var userId pgtype.UUID

	ctxId, err := gc.Cookie("session_id")
	if ctxId == "" {
		fmt.Println(ctxId)
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "User ID is not present",
		})
		return
	}

	err = userId.Scan(ctxId)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to associate User ID to tag",
			"error":   fmt.Sprintf("%#v/n", err),
			"id":      ctxId,
		})
		return
	}

	gc.ShouldBindJSON(&tagData)

	newTagName.Name = tagData.Tag
	newTagName.UserId = userId
	err = gorm.G[database.Tag](database.DB).Create(gc, &newTagName)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create user, please try again later [2]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	gc.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "New tag : " + newTagName.Name + " created",
		"tag":     newTagName,
	})

}

func HandleDeleteTags(gc *gin.Context) {
	var tagId string
	tagId = gc.Param("id")

	_, err := gorm.G[database.Tag](database.DB).Where("id = ?", tagId).Delete(gc)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create user, please try again later [2]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Tag deleted",
	})

}

func HandleGetTags(gc *gin.Context) {
	var queryString PaginationQueryString
	var userId pgtype.UUID

	err := gc.ShouldBindQuery(&queryString)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Invalid parameter received",
		})
		return
	}

	basequery := database.DB.Model(&database.Tag{})

	if queryString.Keyword != "" {
		pattern := "%" + queryString.Keyword + "%"
		basequery.Where("name ILIKE ?", pattern)
	}

	ctxId, err := gc.Cookie("session_id")
	if ctxId == "" || err != nil {
		fmt.Println(ctxId)
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "User ID is not present",
		})
		return
	}

	err = userId.Scan(ctxId)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to associate User ID to tag",
			"error":   fmt.Sprintf("%#v/n", err),
			"id":      ctxId,
		})
		return
	}

	basequery.Where("user_id = ?", userId)

	if queryString.Page < 1 {
		queryString.Page = 1
	}
	if queryString.Limit < 1 {
		queryString.Limit = 20
	}
	offset := (queryString.Page - 1) * queryString.Limit

	var total int64
	var tags []database.Tag
	basequery.Count(&total)
	basequery.Offset(offset).
		Limit(queryString.Limit).
		Find(&tags)

	var pages float64 = 0
	if total > 0 && queryString.Limit > 0 {
		pages = math.Ceil(float64(total+int64(queryString.Limit)-1) / float64(queryString.Limit))
	}

	gc.JSON(http.StatusOK, gin.H{
		"status": true,
		"tags":   tags,
		"metadata": gin.H{
			"page":    queryString.Page,
			"limit":   queryString.Limit,
			"keyword": queryString.Keyword,
			"pages":   pages,
			"total":   total,
		},
	})

}

func HandleSearchTags(gc *gin.Context) {
	var keyword string
	var userId pgtype.UUID
	keyword = gc.Param("query")
	if len(keyword) < 1 {
		return
	}

	ctxId, err := gc.Cookie("session_id")
	if ctxId == "" || err != nil {
		fmt.Println(ctxId)
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "User ID is not present",
		})
		return
	}

	err = userId.Scan(ctxId)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to associate User ID to tag",
			"error":   fmt.Sprintf("%#v/n", err),
			"id":      ctxId,
		})
		return
	}

	var tags []database.Tag
	pattern := "%" + keyword + "%"
	result := database.DB.Model(&database.Tag{}).
		Where("user_id = ?", userId).
		Where("name ILIKE ?", pattern).
		Limit(10).
		Find(&tags)

	if result.Error != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to search tags",
			"error":   fmt.Sprintf("%#v\n", result.Error),
		})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status": true,
		"tags":   tags,
	})

}

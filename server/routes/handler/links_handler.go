package handler

import (
	"fmt"
	"linksaver/server/database"
	"linksaver/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type LinksParams struct {
	Page    int
	Limit   int
	Keyword string
	Total   int64
}

func HandleGetLinks(gc *gin.Context) {
	var queryString LinksParams
	err := gc.ShouldBindQuery(&queryString)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Invalid parameter received",
		})
		return
	}

	fmt.Println("query string", queryString)

	basequery := database.DB.Model(&database.Link{})
	if queryString.Keyword != "" {
		pattern := "%" + queryString.Keyword + "%"
		basequery.Where(
			database.DB.Where("link ILIKE ?", pattern).Or("name ILIKE ?", pattern).Or("link_desc ILIKE ?", pattern),
		)
	}

	var userId pgtype.UUID
	ctxId, _ := gc.Get("ID")
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

	offset := (queryString.Page - 1) * queryString.Limit

	var total int64
	var links []database.Link
	basequery.Count(&total)
	_ = basequery.Preload("LinkTags").Offset(offset).Limit(queryString.Limit).Find(&links)

	gc.JSON(http.StatusOK, gin.H{
		"status": true,
		"links":  links,
		"metadata": gin.H{
			"page":    queryString.Page,
			"limit":   queryString.Limit,
			"total":   queryString.Total,
			"keyword": queryString.Keyword,
		},
	})

}

type CreateLinkInput struct {
	Link     string   `json:"link"`
	Name     string   `json:"name"`
	LinkDesc string   `json:"link_desc"`
	LinkTags []string `json:"link_tags"`
}

func HandleCreateLink(gc *gin.Context) {
	var input CreateLinkInput
	err := gc.ShouldBindJSON(&input)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid input, please try again [1]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	var userId pgtype.UUID
	ctxId, _ := gc.Get("ID")
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

	newLink := database.Link{
		Link:     input.Link,
		Name:     input.Name,
		LinkDesc: input.LinkDesc,
		UserId:   userId,
	}

	err = gorm.G[database.Link](database.DB).Create(gc, &newLink)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create link, please try again later [3]",
			"error":   fmt.Sprintf("%#v\n", err),
		})
		return
	}

	// attaching link tags if exist
	if len(input.LinkTags) > 0 {
		type linkTagRow struct {
			LinkId pgtype.UUID `gorm:"column:link_id"`
			TagId  pgtype.UUID `gorm:"column:tag_id"`
		}
		rows := make([]linkTagRow, len(input.LinkTags))
		for i, tagID := range input.LinkTags {
			var tid pgtype.UUID
			tid.Scan(tagID)
			rows[i] = linkTagRow{LinkId: newLink.ID, TagId: tid}
		}
		err = database.DB.Table("link_tags").CreateInBatches(&rows, len(rows)).Error
		if err != nil {
			gc.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "failed to attach tags to link [4]",
				"error":   fmt.Sprintf("%#v\n", err),
			})
			return
		}
	}

	gc.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Link saved",
		"link":    newLink,
	})

}

type LinksScrollParams struct {
	Offset  int    `form:"offset"`
	Limit   int    `form:"limit"`
	Keyword string `form:"keyword"`
}

func HandleGetLinksScroll(gc *gin.Context) {
	var queryString LinksScrollParams
	err := gc.ShouldBindQuery(&queryString)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Invalid parameter received",
		})
		return
	}

	if queryString.Limit <= 0 {
		queryString.Limit = 20
	}

	basequery := database.DB.Model(&database.Link{})
	if queryString.Keyword != "" {
		pattern := "%" + queryString.Keyword + "%"
		basequery.Where(
			database.DB.Where("link ILIKE ?", pattern).Or("name ILIKE ?", pattern).Or("link_desc ILIKE ?", pattern),
		)
	}

	var userId pgtype.UUID
	ctxId, _ := gc.Get("ID")
	err = userId.Scan(ctxId)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to associate User ID",
			"error":   fmt.Sprintf("%#v/n", err),
		})
		return
	}
	basequery.Where("user_id = ?", userId)

	var links []database.Link
	_ = basequery.Preload("LinkTags").Offset(queryString.Offset).Limit(queryString.Limit).Find(&links)

	hasMore := len(links) == queryString.Limit

	gc.JSON(http.StatusOK, gin.H{
		"status":   true,
		"links":    links,
		"has_more": hasMore,
		"metadata": gin.H{
			"offset":  queryString.Offset,
			"limit":   queryString.Limit,
			"keyword": queryString.Keyword,
		},
	})
}

func HandleGetWebsiteMeta(gc *gin.Context) {
	url := gc.Query("url")
	if url == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "url is required"})
		return
	}

	meta, err := utils.GetWebsiteMeta(url)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status": true,
		"metadata": gin.H{
			"title":       meta.Title,
			"description": meta.Description,
		},
	})
}

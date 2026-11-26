package products

import (
	"ArtiVana/internal/modules/product"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type ProductResponse = product.ProductResponse

// Products godoc
// @Summary Get products
// @Tags Products
// @Produce json
// @Success 200 {array} ProductResponse
// @Param ids query []int true "List of product IDs" example(1, 2, 3)
// @Param subcategory query int false "Subcategory ID" example(10)
// @Param subsubcategory query int false "Subsubcategory ID" example(12)
// @Router /api/products [get]
func (h *Handler) GetProducts(c *gin.Context) {
	idsStrings := c.Request.URL.Query()["ids"]
	ids := make([]int, 0, len(idsStrings))
	for _, id := range idsStrings {
		id, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Invalid products IDs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid products ids"})
			return
		}

		ids = append(ids, id)
	}

	var subcategory int
	var err error
	subcategoryStr := c.Request.URL.Query().Get("subcategory")
	if subcategoryStr != "" {
		subcategory, err = strconv.Atoi(subcategoryStr)
		if err != nil {
			log.Printf("Invalid product subcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid product subcategory"})
			return
		}
	}

	var subsubcategory int
	subsubcategoryStr := c.Request.URL.Query().Get("subcategory")
	if subsubcategoryStr != "" {
		subsubcategory, err = strconv.Atoi(subcategoryStr)
		subsubcategory, err = strconv.Atoi(subcategoryStr)
		if err != nil {
			log.Printf("Invalid product subsubcategory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid product subsubcategory"})
			return
		}
	}

	products, err := h.service.GetProducts(c.Request.Context(), ids, subcategory, subsubcategory)
	if err != nil {
		log.Println("failed to get products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Filter godoc
// @Summary Get products
// @Tags Filter
// @Produce json
// @Success 200 {array} ProductResponse
// @Param brands query []int false "Brand IDs" example(1, 2, 3)
// @Param subcategory query int true "Subcategory ID" example(5)
// @Param subsubcategory query int false "Subsubcategory ID" example(12)
// @Param subsubcategories query []int false "Subsubcategory ID" example(10, 11, 12)
// @Param min_price query int false "Minimum price" example(100)
// @Param max_price query int false "Maximum price" example(500)
// @Param rating query int false "Rating" example(4)
// @Param characteristics query string false "Characteristics" example("color": [1, 2, 3],"size": [10, 20])
// @Param sort query string false "Sorting" example(price_asc)
// @Param page query int false "Page" example(1)
// @Router /api/products/filter [get]
func (h *Handler) GetFilterProducts(c *gin.Context) {
	filter, err := parseCreateFilter(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid filter parameters"})
		return
	}

	products, err := h.service.GetFilterProducts(c.Request.Context(), filter)
	if err != nil {
		log.Println("failed to get products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Filter godoc
// @Summary Get filters and total products
// @Tags Filter
// @Produce json
// @Success 200 {array} FilterMeta
// @Param brands query []int false "Brand IDs" example(1, 2, 3)
// @Param subcategory query int true "Subcategory ID" example(5)
// @Param subsubcategory query int false "Subsubcategory ID" example(12)
// @Param subsubcategories query []int false "Subsubcategory ID" example(10, 11, 12)
// @Param min_price query int false "Minimum price" example(100)
// @Param max_price query int false "Maximum price" example(500)
// @Param rating query int false "Rating" example(4)
// @Param characteristics query string false "Characteristics" example("color": [1, 2, 3],"size": [10, 20])
// @Param sort query string false "Sorting" example(price_asc)
// @Param page query int false "Page" example(1)
// @Router /api/products/filter-meta [get]
func (h *Handler) GetFilterMeta(c *gin.Context) {
	filter, err := parseCreateFilter(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid filter parameters"})
		return
	}

	filters, err := h.service.GetFilterMeta(c.Request.Context(), filter)
	if err != nil {
		log.Println("failed to get filters:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get filters"})
		return
	}

	c.JSON(http.StatusOK, filters)
}

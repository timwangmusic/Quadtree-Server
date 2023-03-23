package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weihesdlegend/quadtree-server/place"
	"github.com/weihesdlegend/quadtree-server/quadtree"
	"go.uber.org/zap"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	MaxNumPlacesPerNode = 100
	MaxTreeDepth        = 5
)

type Server struct {
	quadTree *quadtree.QuadTree
	logger   *zap.Logger
}

func (server *Server) Init() {
	server.quadTree = &quadtree.QuadTree{}
	server.quadTree.Init(MaxNumPlacesPerNode, MaxTreeDepth)
	server.logger, _ = zap.NewProduction()
}

func (server *Server) AddPlace(c *gin.Context) {
	var p place.Place
	if c.BindJSON(&p) == nil {
		server.logger.Info("add place success",
			zap.String("id", p.Id),
			zap.String("name", p.Name),
			zap.String("type", p.Type),
			zap.String("geolocation latitude", fmt.Sprintf("%.4f", p.Location.Lat)),
			zap.String("geolocation longitude", fmt.Sprintf("%.4f", p.Location.Lng)),
		)
		server.quadTree.Insert(p)
		c.String(http.StatusOK, "success")
	} else {
		c.String(http.StatusBadRequest, "bad request")
	}
}

func (server *Server) SearchPlaces(c *gin.Context) {
	lat, latParsingErr := strconv.ParseFloat(c.DefaultQuery("lat", "0.0"), 64)
	if latParsingErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot parse latitude: %s", c.Query("lat"))})
	}
	if math.Abs(lat) > 90.0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("latitude %s is out of range of -90.0 to 90.0", c.Query("lat"))})
	}
	lng, lngParsingErr := strconv.ParseFloat(c.DefaultQuery("lng", "0.0"), 64)
	if lngParsingErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot parse longitude: %s", c.Query("lng"))})
	}
	if math.Abs(lng) > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("longitude %s is out of range of -180.0 to 180.0", c.Query("lng"))})
	}
	geoLocation := place.GeoLocation{
		Lat: lat,
		Lng: lng,
	}
	radius, radParsingErr := strconv.ParseFloat(c.DefaultQuery("radius", "200.0"), 64)
	if radParsingErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot parse radius: %s", c.Query("radius"))})
	}
	server.logger.Info("range search request",
		zap.Float64("lat", lat),
		zap.Float64("lng", lng),
		zap.Float64("radius", radius),
	)

	res := server.quadTree.RangeSearch(&geoLocation, radius)
	c.String(http.StatusOK, fmt.Sprintf("%+v", res))
}

func (server *Server) Run() {
	server.Init()

	router := gin.Default()

	// group endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/places/add", server.AddPlace)
		v1.GET("/places/search", server.SearchPlaces)
	}

	// automatically looking for environment variable PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = ":10086"
	}
	err := router.Run(port)
	if err != nil {
		server.logger.Fatal(err.Error())
	}
}

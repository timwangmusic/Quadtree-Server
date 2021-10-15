package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weihesdlegend/quadtree-server/place"
	"github.com/weihesdlegend/quadtree-server/quadtree"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
)

const (
	MaxNumPlacesPerNode = 100
)

type Server struct {
	quadTree *quadtree.QuadTree
	logger   *zap.Logger
}

func (server *Server) Init() {
	server.quadTree = &quadtree.QuadTree{}
	server.quadTree.Init(MaxNumPlacesPerNode)
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

// TODO: add path parameter validations
func (server *Server) SearchPlaces(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.DefaultQuery("lat", "0.0"), 64)
	lng, _ := strconv.ParseFloat(c.DefaultQuery("lng", "0.0"), 64)
	geoLocation := place.GeoLocation{
		Lat: lat,
		Lng: lng,
	}
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "200"), 64)
	server.logger.Info("range search request",
		zap.Float64("lat", lat),
		zap.Float64("lng", lng),
		zap.Float64("radius", radius),
	)

	res := server.quadTree.RangeSearch(&geoLocation, radius)
	c.String(http.StatusOK, fmt.Sprintf("%+v", res))
}

func (server Server) Run() {
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

package main

// webservice.go
import (
	"fmt"
	"net/http"

	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type FibonacciService struct {
}

// GetFibonacci method returns Fibonacci as float64.
func (fs *FibonacciService) GetFibonacci(ctx *gin.Context) {
	ordinal, err := strconv.Atoi(ctx.Param("ordinal"))
	if err != nil {
		ordinal = 0
	}

	f := fibonacci()
	var result float64
	for iter := 0; iter < ordinal; iter++ {
		result = f()
	}

	ctx.JSON(200, gin.H{
		"Fibonacci": result,
	})
}

// GetMemoizedResults method returns NumberMemoizedResults as int.
func (fs *FibonacciService) GetMemoizedResults(ctx *gin.Context) {
	fibLimit, err := strconv.ParseFloat(ctx.Param("value"), 64)
	if err != nil {
		fibLimit = 0
	}

	fibList, err := GetMemoizedResults(fibLimit)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"NumberMemoizedResults": len(fibList),
	})
}

// ClearDataStore method
func (fs *FibonacciService) ClearDataStore(ctx *gin.Context) {
	err := ClearDataStore()
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ClearDataStore": true,
	})
}

/*************************************************************************************/

// GetHost func returns full hostname from .environment (do NOT include :port)
func GetHost() string {
	return os.Getenv("FIB_API_DOMAIN")
}

// GetPort func
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	return port
}

// ContextOptions func matches axios headers configuration.
// For the server to allow CORS, catch all Preflight OPTIONS requests that the client browser sends before the real query is sent to the SAME URL.
// In general, the pre-flight OPTIONS request doesn't like 301 redirects where the server is caching at different levels.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Clear-Site-Data : Clear-Site-Data: "cache", "cookies", "storage"
// When using a cache policy on the API proxy, ensure that the response of the CORS policy is not cached!
func ContextOptions(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") // GetHost() Don't use Access-Control-Allow-Origin: * if your server is trying to set cookie and you use withCredentials = true.
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// Cache-Control: private, max-age=3600 X-CSRF-Token, Authorization,
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, X-Requested-With, Accept-Encoding, Cache-Control")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, POST, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Accept", "application/json")
	ctx.Writer.Header().Set("Content-Type", "application/json; application/x-www-form-urlencoded; charset=utf-8")

	if ctx.Request.Method != "OPTIONS" {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusOK)
	}
}

func InitializeRoutes(fs *FibonacciService) *gin.Engine { // (client m.ClientCredentials)
	gin.SetMode(gin.ReleaseMode) // Switch to "release" mode in production; or export GIN_MODE=release
	router := gin.Default()

	// Credential is not supported if the CORS header ‘Access-Control-Allow-Origin’ is ‘*’
	// The wildcard asterisk only works for AllowedOrigins. Using the asterisk in AllowedMethods and AllowedHeaders will have no affect.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"User-Agent", "Referrer", "Host", "Token", "Accept", "Content-Type", "Origin", "Content-Length", "X-Requested-With", "Accept-Encoding"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return true // origin == hostName
		},
		MaxAge: 86400,
	}))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./build/static")
	// copy the React build directory to the golang folder when deploying to Docker!
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

	// Direct all routes to index.html:
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", true)
	})

	fib := router.Group("/fib")
	{
		fib.GET("/:ordinal", fs.GetFibonacci)
		fib.GET("/upper/:value", fs.GetMemoizedResults)

		fib.OPTIONS("/clear", ContextOptions)
		fib.POST("/clear", fs.ClearDataStore)
	}

	apiPort := GetPort()
	api := "Handling API calls on " + GetHost() + ":" + apiPort
	fmt.Println(api)
	router.Run(":" + apiPort)
	return router
}

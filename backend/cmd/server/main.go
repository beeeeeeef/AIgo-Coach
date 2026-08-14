package main
import(
	"aigo-coach/backend/internal/analyze"
	"github.com/gin-gonic/gin"
)

func main() {
	r:= gin.Default()
	r.POST("/api/analyze",analyze.HandleAnalyzeMock)
	r.Run(":8080")	
}
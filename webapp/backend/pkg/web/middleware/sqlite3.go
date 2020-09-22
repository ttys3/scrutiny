package middleware

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

func DatabaseMiddleware(appConfig config.Interface, logger logrus.FieldLogger) gin.HandlerFunc {
	//var database *gorm.DB
	fmt.Printf("Trying to connect to database stored: %s\n", appConfig.GetString("web.database.location"))
	database, err := gorm.Open("sqlite3", appConfig.GetString("web.database.location"))
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.SetLogger(&GormLogger{Logger: logger})
	database.AutoMigrate(&db.Device{})
	database.AutoMigrate(&db.SelfTest{})
	database.AutoMigrate(&db.Smart{})
	database.AutoMigrate(&db.SmartAtaAttribute{})
	database.AutoMigrate(&db.SmartNvmeAttribute{})
	database.AutoMigrate(&db.SmartScsiAttribute{})

	//TODO: detrmine where we can call defer database.Close()
	return func(c *gin.Context) {
		c.Set("DB", database)
		c.Next()
	}
}

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct{ Logger logrus.FieldLogger }

// Print handles log events from Gorm for the custom logger.
func (gl *GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		gl.Logger.WithFields(
			logrus.Fields{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Debug(v[3])
	case "log":
		gl.Logger.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

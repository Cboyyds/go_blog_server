package global

// 全局对象
import (
	"server/config"
	
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 全局对象
var (
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	ESClient   *elasticsearch.TypedClient
	Redis      redis.Client
	BlackCache local_cache.Cache
)

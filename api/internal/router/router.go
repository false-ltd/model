package router

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/false-ltd/model/api/internal/config"
	"github.com/false-ltd/model/api/internal/handler"
	"github.com/false-ltd/model/api/internal/middleware"
	"github.com/false-ltd/model/api/frontend"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/false-ltd/model/api/docs/swagger"
)

type Handlers struct {
	Model    *handler.ModelHandler
	Provider *handler.ProviderHandler
	Stats    *handler.StatsHandler
	Atlas    *handler.AtlasHandler
	Compare  *handler.CompareHandler
	Sync     *handler.SyncHandler
	Sitemap  *handler.SitemapHandler
}

func Setup(cfg *config.Config, h *Handlers) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Trust no proxy by default: behind the k8s nginx ingress,
	// MODEL_TRUSTED_PROXIES must list the proxy CIDRs, otherwise
	// X-Forwarded-For could be spoofed to bypass the rate limiter.
	if err := r.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
		panic(err)
	}

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	// Compress JSON/sitemap/text responses ≥1KB (atlas: 913KB → ~142KB).
	// Already-compressed formats are excluded.
	r.Use(gzip.Gzip(
		gzip.DefaultCompression,
		gzip.WithMinLength(1024),
		gzip.WithExcludedExtensions([]string{".png", ".gif", ".jpg", ".jpeg", ".webp", ".woff", ".woff2"}),
	))

	limiter := middleware.NewRateLimiter()

	r.GET("/health", handler.Health)
	r.GET("/sitemap.xml", h.Sitemap.Sitemap)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(limiter.Middleware(60))
	{
		api.GET("/models", h.Model.List)
		api.GET("/models/:id", h.Model.Get)

		api.GET("/providers", h.Provider.List)

		api.GET("/stats", h.Stats.Get)

		api.GET("/atlas", h.Atlas.Get)

		api.GET("/compare", h.Compare.Compare)
	}

	// Sync gets its own, stricter limiter — mounting it outside the api
	// group avoids two simultaneous token buckets per IP.
	sync := r.Group("/api/v1/sync")
	sync.Use(limiter.Middleware(10))
	{
		sync.POST("", middleware.APIKeyAuth(cfg.Auth.APIKeys, cfg.Sync.PublicTrigger), h.Sync.Trigger)
		sync.GET("/status", h.Sync.Status)
	}

	return r
}

// SetupSPA serves embedded frontend static files with SPA fallback.
// Only whitelisted SPA routes get index.html; everything else returns 404.
func SetupSPA(r *gin.Engine) {
	sub, err := fs.Sub(frontend.DistFS, "dist")
	if err != nil {
		return
	}
	fileServer := http.FileServer(http.FS(sub))

	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		trimmed := strings.TrimPrefix(reqPath, "/")

		// Serve known static files (_nuxt/*, favicon, etc.)
		if trimmed != "" {
			if f, err := fs.Stat(sub, trimmed); err == nil && !f.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// Only whitelisted SPA routes get index.html fallback
		if isSPARoute(path.Clean(reqPath)) {
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.Status(http.StatusNotFound)
	})
}

// isSPARoute checks if the path is a known SPA route.
// Pages: /, /catalog, /compare, /providers, /model/:id
// i18n prefix_except_default: /zh prefix for Chinese locale.
func isSPARoute(p string) bool {
	p = path.Clean(p)
	allowed := []string{
		"/",
		"/catalog",
		"/compare",
		"/providers",
		"/zh",
		"/zh/catalog",
		"/zh/compare",
		"/zh/providers",
	}
	for _, a := range allowed {
		if p == a {
			return true
		}
	}
	// Dynamic routes: /model/:id, /zh/model/:id
	if strings.HasPrefix(p, "/model/") || p == "/model" ||
		strings.HasPrefix(p, "/zh/model/") || p == "/zh/model" {
		return true
	}
	return false
}

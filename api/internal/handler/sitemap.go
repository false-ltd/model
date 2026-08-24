package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/repository"
	"github.com/gin-gonic/gin"
)

const sitemapCacheKey = "sitemap"
const sitemapCacheTTL = 5 * time.Minute

type SitemapHandler struct {
	modelRepo *repository.ModelRepo
	siteURL   string
	cache     *cache.Cache
}

func NewSitemapHandler(modelRepo *repository.ModelRepo, siteURL string, c *cache.Cache) *SitemapHandler {
	return &SitemapHandler{modelRepo: modelRepo, siteURL: siteURL, cache: c}
}

func (h *SitemapHandler) Sitemap(c *gin.Context) {
	v, err := h.cache.GetOrCompute(sitemapCacheKey, sitemapCacheTTL, h.build)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to generate sitemap")
		return
	}
	c.Data(http.StatusOK, "application/xml; charset=utf-8", v.([]byte))
}

func (h *SitemapHandler) build() (any, error) {
	ids, err := h.modelRepo.FindAllIDs()
	if err != nil {
		return nil, err
	}

	now := time.Now().Format("2006-01-02")
	locales := []struct{ prefix, lang string }{
		{"", "en"},
		{"/zh", "zh"},
	}

	var b strings.Builder

	// Static pages
	staticPages := []struct {
		path       string
		changefreq string
		priority   string
	}{
		{"/", "daily", "1.0"},
		{"/catalog", "daily", "0.9"},
		{"/providers", "weekly", "0.8"},
		{"/compare", "weekly", "0.6"},
	}
	for _, p := range staticPages {
		h.urlEntry(&b, p.path, now, p.changefreq, p.priority, locales)
	}

	// Model pages
	for _, id := range ids {
		h.urlEntry(&b, fmt.Sprintf("/model/%d", id), now, "weekly", "0.7", locales)
	}

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
%s</urlset>`, b.String())

	return []byte(xml), nil
}

func (h *SitemapHandler) urlEntry(b *strings.Builder, path, lastmod, changefreq, priority string, locales []struct{ prefix, lang string }) {
	fmt.Fprintf(b, "  <url>\n    <loc>%s%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n",
		h.siteURL, path, lastmod, changefreq, priority)
	for _, l := range locales {
		fmt.Fprintf(b, "    <xhtml:link rel=\"alternate\" hreflang=\"%s\" href=\"%s%s%s\"/>\n",
			l.lang, h.siteURL, l.prefix, path)
	}
	fmt.Fprintf(b, "    <xhtml:link rel=\"alternate\" hreflang=\"x-default\" href=\"%s%s\"/>\n", h.siteURL, path)
	b.WriteString("  </url>\n")
}

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/false-ltd/model/api/internal/repository"
	"github.com/gin-gonic/gin"
)

type SitemapHandler struct {
	modelRepo    *repository.ModelRepo
	providerRepo *repository.ProviderRepo
	siteURL      string
}

func NewSitemapHandler(modelRepo *repository.ModelRepo, providerRepo *repository.ProviderRepo, siteURL string) *SitemapHandler {
	return &SitemapHandler{modelRepo: modelRepo, providerRepo: providerRepo, siteURL: siteURL}
}

func (h *SitemapHandler) Sitemap(c *gin.Context) {
	ids, err := h.modelRepo.FindAllIDs()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to generate sitemap")
		return
	}

	now := time.Now().Format("2006-01-02")
	urlset := ""
	locales := []struct{ prefix, lang string }{
		{"", "en"},
		{"/zh", "zh"},
	}

	// Static pages
	staticPages := []struct {
		path      string
		changefreq string
		priority   string
	}{
		{"/", "daily", "1.0"},
		{"/catalog", "daily", "0.9"},
		{"/providers", "weekly", "0.8"},
		{"/compare", "weekly", "0.6"},
	}

	for _, p := range staticPages {
		urlset += h.urlEntry(p.path, now, p.changefreq, p.priority, locales)
	}

	// Model pages
	for _, id := range ids {
		urlset += h.urlEntry(fmt.Sprintf("/model/%d", id), now, "weekly", "0.7", locales)
	}

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
%s</urlset>`, urlset)

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(xml))
}

func (h *SitemapHandler) urlEntry(path, lastmod, changefreq, priority string, locales []struct{ prefix, lang string }) string {
	s := fmt.Sprintf("  <url>\n    <loc>%s%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n",
		h.siteURL, path, lastmod, changefreq, priority)

	for _, l := range locales {
		s += fmt.Sprintf("    <xhtml:link rel=\"alternate\" hreflang=\"%s\" href=\"%s%s%s\"/>\n",
			l.lang, h.siteURL, l.prefix, path)
	}
	s += fmt.Sprintf("    <xhtml:link rel=\"alternate\" hreflang=\"x-default\" href=\"%s%s\"/>\n", h.siteURL, path)
	s += "  </url>\n"
	return s
}

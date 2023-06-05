package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iKayrat/website_test/checker"
)

func main() {
	check := checker.New()

	websites := []string{
		"https://www.example.com",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.youtube.com",
		"https://www.facebook.com",
		"https://www.baidu.com",
		"https://www.wikipedia.org",
		"https://www.qq.com",
		"https://www.taobao.com",
		"https://www.yahoo.com",
		"https://www.tmall.com",
		"https://www.amazon.com",
		"https://www.google.co.in",
		"https://www.twitter.com",
		"https://www.sohu.com",
		"https://www.jd.com",
		"https://www.live.com",
		"https://www.instagram.com",
		"https://www.sina.com.cn",
		"https://www.weibo.com",
		"https://www.google.co.jp",
		"https://www.reddit.com",
		"https://www.vk.com",
		"https://www.360.cn",
		"https://www.login.tmall.com",
		"https://www.blogspot.com",
		"https://www.yandex.ru",
		"https://www.google.com.hk",
		"https://www.netflix.com",
		"https://www.linkedin.com",
		"https://www.pornhub.com",
		"https://www.google.com.br",
		"https://www.twitch.tv",
		"https://www.pages.tmall.com",
		"https://www.csdn.net",
		"https://www.yahoo.co.jp",
		"https://www.mail.ru",
		"https://www.aliexpress.com",
		"https://www.alipay.com",
		"https://www.office.com",
		"https://www.google.fr",
		"https://www.google.ru",
		"https://www.google.co.uk",
		"https://www.microsoftonline.com",
		"https://www.google.de",
		"https://www.ebay.com",
		"https://www.microsoft.com",
		"https://www.livejasmin.com",
		"https://www.t.co",
		"https://www.bing.com",
		"https://www.xvideos.com",
		"https://www.google.ca",
	}

	app := fiber.New()

	go func() {
		for {
			check.CheckWebsite(context.Background(), websites)
			time.Sleep(time.Minute)
		}
	}()

	app.Get("/access", func(c *fiber.Ctx) error {
		url := c.Query("url")

		if url == "" {
			msg := "Missing 'url' parameter"

			return c.Status(fiber.StatusBadRequest).JSON(wrap(msg))
		}

		site, found := check.GetAccessTime(url)
		if found {
			check.IncCounter("/access")
			return c.JSON(site)

		} else {

			msg := fmt.Sprintf("Site %s is not accessible", url)

			return c.JSON(wrap(msg))
		}

	})

	app.Get("/min", func(c *fiber.Ctx) error {
		status := check.GetMinAccessTime()

		if status == nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap("No accessible sites found"))
		}

		check.IncCounter("/min")
		return c.Status(fiber.StatusOK).JSON(status)
	})

	app.Get("/max", func(c *fiber.Ctx) error {
		status := check.GetMaxAccessTime()

		if status == nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap("No accessible sites found"))
		}

		check.IncCounter("/max")

		return c.Status(fiber.StatusOK).JSON(status)
	})

	// Middleware
	app.Use(isAdmin)

	app.Get("/counts", func(c *fiber.Ctx) error {

		counts := check.GetCounts()

		if len(counts) <= 0 {

			return c.Status(fiber.StatusOK).JSON(wrap("There is no stats yet"))
		}

		return c.Status(fiber.StatusOK).JSON(counts)

	})

	log.Fatal(app.Listen(":8080"))

}

// Middleware
func isAdmin(c *fiber.Ctx) error {
	role := c.Get("Role")
	if role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(wrap("Admin required"))
	}

	return c.Next()
}

// Response
type message struct {
	Msg string `json:"message"`
}

func wrap(msg string) message {
	return message{
		msg,
	}
}

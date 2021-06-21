package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var (
	store                   CacheStorage
	bypassCacheHeader       = "Bypass-Cache"
	serverCacheStatusHeader = "S-Cache"
)

func main() {
	store = NewCacheStorage()
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(store.db)

	// Create a new html engine
	engine := html.New("./views", ".html")

	// start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// default to no store
	app.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})

	app.Get("/cache", allowCache(5), func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})

	app.Get("/server-cache", allowServerCache(5), func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})

	app.Get("/server-and-cdn-cache", allowServerCache(30), allowCache(5), func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})

	app.Get("/timeout", func(c *fiber.Ctx) error {
		timeout, _ := strconv.Atoi(c.Query("t", "100"))
		dur := time.Duration(timeout) * time.Millisecond
		time.Sleep(dur)
		return c.JSON(fiber.Map{"waited": dur, "hint": "use query t to indicate timeout in ms"})
	})

	app.Get("/error/:statusCode", func(c *fiber.Ctx) error {
		code, _ := strconv.Atoi(c.Params("statusCode"))
		return c.Render("error", fiber.Map{
			"errorCode": code,
			"errorText": http.StatusText(code),
		})
	})

	app.Get("/server-1/example", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"server": 1, "status": "ok"})
	})

	app.Get("/server-2/example", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"server": 2, "status": "ok"})
	})

	dbRoutes := app.Group("/db")
	dbRoutes.Get("/all", func(c *fiber.Ctx) error {
		_ = store.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			it := txn.NewIterator(opts)
			defer it.Close()
			fmt.Printf("----Query at %s------\n", time.Now().Format("2-Jan-06 3:04:05 pm"))
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				err := item.Value(func(v []byte) error {
					fmt.Printf("key=%s, value=%s, exp=%d, isExp=%t\n", k, v, item.ExpiresAt(), item.IsDeletedOrExpired())
					return nil
				})
				if err != nil {
					return err
				}
			}
			fmt.Println("----End Query---")
			return nil
		})
		return c.JSON(fiber.Map{"result": "view logs for result"})
	})

	log.Fatal(app.Listen(":3000"))
}

// allowCache set Cache-Control for client-side and CDN caching
func allowCache(dur int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", dur))
		return c.Next()
	}
}

// allowServerCache uses server-side caching
func allowServerCache(ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// bypass cache
		if bypassCache(c) {
			c.Response().Header.Set(serverCacheStatusHeader, "unreachable")
			return c.Next()
		}

		// Only cache GET methods
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Get key from request
		key := c.Path()

		// get record from storage
		get, err := store.Get(key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var item *CacheItem

		// has record
		if get != nil {
			if err := json.Unmarshal(get, &item); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			c.Response().SetBodyRaw(item.Body)
			c.Response().SetStatusCode(item.Status)
			c.Response().Header.SetContentTypeBytes(item.Ctype)
			if len(item.Cencoding) > 0 {
				c.Response().Header.SetBytesV(fiber.HeaderContentEncoding, item.Cencoding)
			}
			c.Response().Header.Set(serverCacheStatusHeader, "hit")
			return nil
		}
		// continue stack
		c.Response().Header.Set(serverCacheStatusHeader, "miss")
		if err := c.Next(); err != nil {
			return err
		}

		// Cache response
		item = new(CacheItem)
		item.Body = utils.CopyBytes(c.Response().Body())
		item.Status = c.Response().StatusCode()
		item.Ctype = utils.CopyBytes(c.Response().Header.ContentType())
		item.Cencoding = utils.CopyBytes(c.Response().Header.Peek(fiber.HeaderContentEncoding))

		bytes, err := json.Marshal(item)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := store.Set(key, bytes, ttl*time.Second); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Finish response
		return nil
	}
}

func bypassCache(c *fiber.Ctx) bool {
	if c.Get(bypassCacheHeader) == "1" || c.Get(bypassCacheHeader) == "true" ||
		c.Query("refresh") == "true" || c.Query("refresh") == "1" {
		return true
	}
	return false
}

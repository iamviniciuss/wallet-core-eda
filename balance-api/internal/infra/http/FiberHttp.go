package infra

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/errors"
	"github.com/valyala/fasthttp"
)

type FiberHttp struct {
	app          *fiber.App
	CustomParams QueryParams
}

func NewFiberHttp() *FiberHttp {
	f := new(FiberHttp)
	f.app = fiber.New()

	f.app.Use(cors.New())

	f.app.Use(func(c *fiber.Ctx) error {
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("Access-Control-Allow-Origin", "*")

		return c.Next()
	})

	return f
}

func (f *FiberHttp) Get(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError)) {
	f.app.Get(path, func(c *fiber.Ctx) error {
		result, err := callback(c.AllParams(), c.Body(), f.CustomParams)

		if err != nil {
			c.Status(200)
			return c.JSON(err)
		}

		return c.JSON(result)
	})
}

func (f *FiberHttp) Post(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError)) {
	f.app.Post(path, func(c *fiber.Ctx) error {
		result, err := callback(c.AllParams(), c.Body(), f.CustomParams)

		if err != nil {
			c.Status(200)
			return c.JSON(err)
		}

		return c.JSON(result)
	})
}

func (f *FiberHttp) Put(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError)) {
	f.app.Put(path, func(c *fiber.Ctx) error {
		result, err := callback(c.AllParams(), c.Body(), f.CustomParams)

		if err != nil {
			c.Status(err.StatusCode)
			return c.JSON(err)
		}

		return c.JSON(result)
	})
}

func (f *FiberHttp) ListenAndServe(port string) error {
	return f.app.Listen(port)
}

type FiberQueryParams struct {
	Args *fasthttp.Args
}

func NewFiberQueryParams(args *fasthttp.Args) *FiberQueryParams {
	return &FiberQueryParams{
		Args: args,
	}
}

func (fqp *FiberQueryParams) GetParam(key string) []byte {
	return fqp.Args.Peek(key)
}

func (fqp *FiberQueryParams) AddParam(key string, value string) {
	fqp.Args.Add(key, value)
}

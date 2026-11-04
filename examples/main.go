package main

import (
	discovery "github.com/zeiss/fiber-discovery"

	"github.com/gofiber/fiber/v3"
)

type response struct {
	Endpoint string `json:"endpoint" xml:"endpoint" msgpack:"endpoint" cbor:"endpoint"`
}

func main() {
	app := fiber.New()

	app.Get(discovery.WellknownURL("builder-configuration"), discovery.New(response{Endpoint: "http://localhost:3000"}))

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

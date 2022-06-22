package main

import (
	"context"
	"log"

	"github.com/adduc/terraform-provider-woodpecker/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/adduc/woodpecker",
	}

	err := providerserver.Serve(context.Background(), internal.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

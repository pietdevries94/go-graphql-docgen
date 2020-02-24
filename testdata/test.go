package main

import "github.com/pietdevries94/go-graphql-docgen/testoutput"

import "fmt"

func main() {
	c := testoutput.NewClient("https://favware.tech/api")

	pika, err := c.GetPikachu(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(pika.GetPokemonDetails.Species)

	charm, err := c.GetExact(testoutput.PokemonCharmander, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(charm.GetPokemonDetailsByName.Species)

	morePika, err := c.GetMorePikachu(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(morePika.GetPokemonDetails.GetFullType().FlavorTexts[0].Flavor)
}

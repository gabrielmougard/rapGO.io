package main

// go build -o ../../bin/musicParser/parser

func main() {

	//Other interesting genre :
	// "Soundtrack", "Nerdcore", "Wonky"
	baseURLs := []string{"Hip-Hop_Beats","Hip-Hop"}
	var mp MusicParser
	mp.Seed(baseURLs,true,23)
	mp.Start()
}
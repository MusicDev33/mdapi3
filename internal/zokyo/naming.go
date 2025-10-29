package zokyo

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var adjectives = []string{
	"adorable", "brave", "careless", "dazzling", "eager", "fearless", "gigantic",
	"hilarious", "icy", "jolly", "keen", "lazy", "mysterious", "naughty", "optimistic",
	"perfect", "quaint", "radiant", "silly", "terrific", "unusual", "vibrant", "witty",
	"zealous", "alert", "beautiful", "charming", "delightful", "elegant", "fascinating",
	"graceful", "handsome", "innovative", "joyful", "knowledgeable", "lovely", "majestic",
	"natural", "original", "proud", "quick", "resilient", "steadfast", "tender", "upbeat",
	"versatile", "whimsical", "youthful", "zany",
}

var nouns = []string{
	"aardvark", "butterfly", "chameleon", "dolphin", "echidna", "flamingo", "gazelle",
	"hedgehog", "impala", "jaguar", "koala", "lemming", "marmoset", "narwhal", "ocelot",
	"panther", "quokka", "rhinoceros", "swordfish", "tapir", "urchin", "vulture", "wombat",
	"xenops", "yak", "zebra", "airplane", "backpack", "cactus", "dandelion", "eggplant",
	"fireplace", "guitar", "honey", "igloo", "jacket", "kangaroo", "lighthouse", "mask",
	"nightingale", "orchid", "pinecone", "quartz", "rainbow", "saxophone", "trampoline",
	"umbrella", "violin", "whirlpool", "xylophone", "yacht", "zeppelin",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func GenerateName() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s %s", capitalize(adj), capitalize(noun))
}

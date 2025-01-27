package name

import "math/rand"

var syllables = []string{
	"ab", "ac", "ad", "af", "ag", "ah", "aj", "ak", "al", "am", "an", "ap", "aq", "ar", "as", "at", "av", "aw", "ax", "ay", "az",
	"ba", "be", "bi", "bo", "bu",
	"ca", "ce", "ci", "co", "cu",
	"da", "de", "di", "do", "du",
	"ea", "ee", "ei", "eo", "eu",
	"fa", "fe", "fi", "fo", "fu",
	"ga", "ge", "gi", "go", "gu",
	"ha", "he", "hi", "ho", "hu",
	"ia", "ie", "ii", "io", "iu",
	"ja", "je", "ji", "jo", "ju",
	"ka", "ke", "ki", "ko", "ku",
	"la", "le", "li", "lo", "lu",
	"ma", "me", "mi", "mo", "mu",
	"na", "ne", "ni", "no", "nu",
	"oa", "oe", "oi", "oo", "ou",
	"pa", "pe", "pi", "po", "pu",
	"qa", "qe", "qi", "qo", "qu",
	"ra", "re", "ri", "ro", "ru",
	"sa", "se", "si", "so", "su",
	"ta", "te", "ti", "to", "tu",
	"ua", "ue", "ui", "uo", "uu",
	"va", "ve", "vi", "vo", "vu",
	"wa", "we", "wi", "wo", "wu",
	"xa", "xe", "xi", "xo", "xu",
	"ya", "ye", "yi", "yo", "yu",
	"za", "ze", "zi", "zo", "zu",
}

func Generate(numberOfSyllables int) string {

	n := ""

	for i := 0; i < numberOfSyllables; i++ {
		n += syllables[rand.Intn(len(syllables))]
	}

	return n
}

var nickNames = []string{
	"great",
	"mighty",
	"powerful",
	"strong",
	"weak",
	"small",
	"tiny",
	"big",
	"large",
	"massive",
	"kind",
	"evil",
	"good",
	"bad",
	"ugly",
}

func GenerateNickName() string {
	return nickNames[rand.Intn(len(nickNames))]
}

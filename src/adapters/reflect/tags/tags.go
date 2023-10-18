package tags

import (
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/src/helpers"
)

type Tag struct {
	Name  string   `required:"true"`
	Rules []string `required:"true"` //split by comma
}

type Tags []Tag

func isSpace(c string) bool {
	return c == " " || c == "\t"
}

func isAlphanumeric(c string) bool {
	return c >= "a" && c <= "z" ||
		c >= "A" && c <= "Z" ||
		c >= "0" && c <= "9"
}

func isKeyValSeparator(c string) bool {
	return c == ":"
}

// there might be some confusion between the tag value separator and the tag value delimiter
// the tag value separator will be used to split the tag value into multiple values (e.g. "a,b,c" -> ["a", "b", "c"])
// the tag value delimiter will be used to escape the tag value separator (e.g. "a,b,c" -> ["a,b,c"])
// I'm open to suggestions on how to make this less confusing naming-wise
func isTagValueSeparator(c string) bool {
	return c == ","
}

func isTagValueDelimiter(c string) bool {
	return c == "\""
}

func bannedTagNames() []string {
	return []string{
		"json",
		"xml",
		"yaml",
		"bson",
	}
}

func isBannedTagName(name string) bool {
	for _, n := range bannedTagNames() {
		if n == name {
			return true
		}
	}

	return false
}

func readToken(raw string, i *int) string {
	token := ""

	for *i < len(raw) && isAlphanumeric(string(raw[*i])) {
		token += string(raw[*i])
		*i++
	}

	return token
}

func skipSpaces(raw string, i *int) {
	for *i < len(raw) && isSpace(string(raw[*i])) {
		*i++
	}
}

func GetTags(raw string) Tags {
	if len(raw) == 0 {
		return nil
	}

	tagList := Tags{}

	i := 0
	for i < len(raw) {
		skipSpaces(raw, &i)
		name := readToken(raw, &i)
		skipSpaces(raw, &i)

		if i >= len(raw) || !isKeyValSeparator(string(raw[i])) {
			log.Warn().Str("raw", raw).Str("tag", name).Str("reason", "tag must be followed by a colon").Msg("failed to parse tag")
			return nil
		}
		i++

		if i >= len(raw) || !isTagValueDelimiter(string(raw[i])) {
			log.Warn().Str("raw", raw).Str("tag", name).Str("reason", "tag value must be enclosed in double quotes").Msg("failed to parse tag")
			return nil
		}
		i++

		skipSpaces(raw, &i)

		tag := Tag{}
		tag.Name = name

		for i < len(raw) && !isTagValueDelimiter(string(raw[i])) {
			for i < len(raw) {
				skipSpaces(raw, &i)
				token := readToken(raw, &i)
				tag.Rules = append(tag.Rules, token)

				if i >= len(raw) || !isTagValueSeparator(string(raw[i])) {
					break
				}
				i++
			}

			if !isTagValueDelimiter(string(raw[i])) {
				log.Warn().Str("raw", raw).Str("tag", name).Str("reason", "tag value must be enclosed in double quotes").Msg("failed to parse tag")
				return nil
			}

			skipSpaces(raw, &i)
		}

		if !isBannedTagName(tag.Name) {
			tagList = append(tagList, tag)
		}

		i++
	}

	return helpers.ArrayOrNil(tagList)
}

func (ts Tags) Required() bool {
	for _, t := range ts {
		if t.Name == "required" {
			if len(t.Rules) > 0 {
				return t.Rules[0] == "true"
			}
		}
	}

	return false
}

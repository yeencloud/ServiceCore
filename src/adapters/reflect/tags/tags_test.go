package tags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoTags(t *testing.T) {
	type ExampleModule struct {
	}

	tagList := GetTags("")

	assert.Len(t, tagList, 0)
}

func TestOneTag(t *testing.T) {
	tagList := GetTags(`validate:"example"`)

	success := assert.Len(t, tagList, 1)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
}

func TestTagWithMultipleValues(t *testing.T) {
	tagList := GetTags(`validate:"example,example2"`)

	success := assert.Len(t, tagList, 1)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example", "example2"})
}

func TestMultipleTags(t *testing.T) {
	tagList := GetTags(`greaterThan:"0" lesserThan:"10"`)

	success := assert.Len(t, tagList, 2)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "greaterThan")
	assert.Equal(t, tagList[0].Rules, []string{"0"})
	assert.Equal(t, tagList[1].Name, "lesserThan")
	assert.Equal(t, tagList[1].Rules, []string{"10"})
}

func TestMultipleTagsWithMultipleValues(t *testing.T) {
	tagList := GetTags(`not:"debug,dev" excludes:"preprod"`)

	success := assert.Len(t, tagList, 2)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "not")
	assert.Equal(t, tagList[0].Rules, []string{"debug", "dev"})
	assert.Equal(t, tagList[1].Name, "excludes")
	assert.Equal(t, tagList[1].Rules, []string{"preprod"})
}

func TestBannedTags(t *testing.T) {
	tagList := GetTags(`json:"example" validate:"example"`)

	success := assert.Len(t, tagList, 1)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
}

func TestTagValueWithMultipleSpaces(t *testing.T) {
	tagList := GetTags(`validate:"  example  "`)

	success := assert.Len(t, tagList, 1)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
}

func TestTagNameWithMultipleSpaces(t *testing.T) {
	tagList := GetTags(`  validate  :"example"`)

	success := assert.Len(t, tagList, 1)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
}

func TestSpacingBetweenTwoTags(t *testing.T) {
	tagList := GetTags(`validate:"example"    		excludes:"preprod"`)

	success := assert.Len(t, tagList, 2)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
	assert.Equal(t, tagList[1].Name, "excludes")
	assert.Equal(t, tagList[1].Rules, []string{"preprod"})
}

func TestMultipleTagsWithMultipleSpaces(t *testing.T) {
	tagList := GetTags(`  validate  :"example"		  excludes  :  "preprod"  `)

	success := assert.Len(t, tagList, 2)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example"})
	assert.Equal(t, tagList[1].Name, "excludes")
	assert.Equal(t, tagList[1].Rules, []string{"preprod"})
}

func TestMultipleTagsWithMultipleValuesAndSpaces(t *testing.T) {
	tagList := GetTags(`  validate  :" example ,  example2  "		  excludes  :  "preprod , preprod2  "  `)

	success := assert.Len(t, tagList, 2)
	if !success {
		return
	}
	assert.Equal(t, tagList[0].Name, "validate")
	assert.Equal(t, tagList[0].Rules, []string{"example", "example2"})
	assert.Equal(t, tagList[1].Name, "excludes")
	assert.Equal(t, tagList[1].Rules, []string{"preprod", "preprod2"})
}
package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"unicode"
)

const (
	equalName   = "equal"
	greaterName = "greater"
	lowerName   = "lower"
)

func dataSourceVersionValidator() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Version validator datasource.",

		ReadContext: dataSourceVersionValidatorRead,

		Schema: map[string]*schema.Schema{
			"current_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Currently deployed version.",
			},
			"required_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required version.",
			},
		},
	}
}

func dataSourceVersionValidatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(".")
	if err != nil {
		return diags
	}

	tagrefs, err := r.Tags()
	if err != nil {
		return diags
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t)
		return nil
	})
	if err != nil {
		return diags
	}
	tags, err := r.TagObjects()
	if err != nil {
		return diags
	}
	var array []string
	err = tags.ForEach(func(t *object.Tag) error {
		array = append(array, t.String())
		return nil
	})
	if err != nil {
		return diags
	}
	diags = diag.Errorf(fmt.Sprintf("Tags: %v", array))
	return diags
}

func compare(currentVersionArrayInt []int, requiredVersionArrayInt []int) string {
	if currentVersionArrayInt[0] > requiredVersionArrayInt[0] {
		return greaterName
	} else if currentVersionArrayInt[0] == requiredVersionArrayInt[0] {
		if currentVersionArrayInt[1] > requiredVersionArrayInt[1] {
			return greaterName
		} else if currentVersionArrayInt[1] == requiredVersionArrayInt[1] {
			if currentVersionArrayInt[2] > requiredVersionArrayInt[2] {
				return greaterName
			} else if currentVersionArrayInt[2] == requiredVersionArrayInt[2] {
				return equalName
			}
		} else {
			return lowerName
		}
	} else {
		return lowerName
	}
	return ""
}

func DigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}

func DigitPostfix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[i:]
		}
	}
	return s
}

func fillZeroes(arr []string) (resultArr []string) {
	resultArr = []string{"0", "0", "0"}
	for index, element := range arr {
		resultArr[index] = element
	}
	return
}

func stringArrToNumbers(stringArr []string) (intArr []int) {
	intArr = []int{0, 0, 0}
	for i, v := range stringArr {
		err := errors.New("")
		intArr[i], err = strconv.Atoi(v)
		if err != nil {
			diag.Errorf(fmt.Sprintf("Error: %s", err))
		}
	}
	return
}

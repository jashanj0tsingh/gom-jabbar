package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"lapoutinemtl.com/utils"
	"net/http"
)

// add more choices if required!
var oilChoices = [3]string{"vegetable", "canola", "olive"}

func validateOilChoice(oilChoice string) string {
	for i := range oilChoices {
		if oilChoice == oilChoices[i] {
			return oilChoice
		}
	}
	return ""
}

// todo: sick of writing the same client code again and again, refactor it out
func fryPotatoes(context echo.Context) error {

	qp := context.QueryParams()
	oc := qp.Get("oilChoice")

	// validate oil choice
	oc = validateOilChoice(oc)
	if oc == "" {
		return context.JSON(http.StatusBadRequest, fmt.Sprintf(`%s not a valid oil choice`, oc))
	}
	// request nordo for cooked potatoes
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("http://host.docker.internal:5144/nordo/potatoes/boil")
	if err != nil {
		fmt.Errorf("err is: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var potatoes utils.CutPotatoes
	var fries utils.PotatoFries
	err = json.Unmarshal(body, &potatoes)
	if err != nil {}

	if potatoes.Portion != "" && potatoes.AreCooked && potatoes.DippedIn == "maple syrup" {
		fries.AreHotAndCrispy = true
		return context.JSON(http.StatusAccepted, struct {
			utils.PotatoFries `json:"potatoFries"`
		}{fries})
	} else {
		return context.JSON(http.StatusBadRequest, nil)
	}
}


package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func apiGet(link string, obj interface{}) error {
	response, _ := http.Get(link)
	if response.StatusCode != 200 {
		return errors.New("statusCode = " + strconv.Itoa(response.StatusCode))
	}
	responseData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseData, obj)
	return nil
}

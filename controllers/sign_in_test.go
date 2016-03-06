package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSignIn(t *testing.T) {
	bodyStr := "{\"Email\":\"work_test_b@163.com\", \"Password\":\"123456\"}"
	body := ioutil.NopCloser(strings.NewReader(bodyStr))
	client := &http.Client{}

	reqs, _ := http.NewRequest("POST", url+"/uauth/signin", body)
	reqs.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(reqs)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	remap := make(map[string]interface{})
	json.Unmarshal(data, &remap)

	if fmt.Sprint(remap["Status"]) != "200" {
		t.Error("Status Expected:200 Actual:", remap["Status"])
	}
	if remap["Message"] != "OK" {
		t.Error("Message Expected:OK Actual:", remap["Message"])
	}

}

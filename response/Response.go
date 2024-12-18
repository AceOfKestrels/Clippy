package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func Deserialize(resp *http.Response) (Response, error) {
	var response Response

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal([]byte(string(body)), &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

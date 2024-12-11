package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	ModelVersion string `json:"modelVersion"`
}

func Deserialize(resp *http.Response) (Response, error) {
	var response Response

	body, err := ioutil.ReadAll(resp.Body)
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

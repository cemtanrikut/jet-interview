package newcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const assetsURL = "https://YOUR_SUBDOMAIN.rest.marketingcloudapis.com/asset/v1/content/assets/query"

type Asset struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Content      string    `json:"content"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

// FetchUpdatedAssets fetches updated or new assets from Salesforce.
func FetchUpdatedAssets(accessToken string, since time.Time) ([]Asset, error) {
	query := fmt.Sprintf(`{
        "query": {
            "leftOperand": {
                "property": "modifiedDate",
                "simpleOperator": "greaterThan",
                "value": "%s"
            }
        }
    }`, since.Format(time.RFC3339))

	req, err := http.NewRequest("POST", assetsURL, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Err req: %s", resp.Status)
	}

	var result struct {
		Items []Asset `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

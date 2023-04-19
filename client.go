package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

type BitwardenClient struct {
	BaseURL      string
	SessionToken string
	Client       *http.Client
}

func NewBitwardenClient(conf Config) (BitwardenClient, error) {
	bwClient := BitwardenClient{
		BaseURL: conf.Addr,
	}

	if os.Getenv("BW_PASSWORD") == "" {
		return BitwardenClient{}, fmt.Errorf("bitwarden password not set in BW_PASSWORD")
	}
	loginCred := map[string]string{
		"password": os.Getenv("BW_PASSWORD"),
	}

	json_data, err := json.Marshal(loginCred)
	if err != nil {
		return BitwardenClient{}, fmt.Errorf("failed to encode login creds: %w", err)
	}

	resp, err := http.Post(bwClient.BaseURL+"/unlock", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return BitwardenClient{}, fmt.Errorf("unable to authenticate to BW host: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp UnlockResp
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	bwClient.SessionToken = tokenResp.Data["raw"].(string)
	bwClient.Client = &http.Client{}
	return bwClient, nil
}

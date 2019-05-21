package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/VideoCoin/common/proto"
	"github.com/sirupsen/logrus"
)

type CloudManagerConfig struct {
	ManagerAddr string
	Logger      *logrus.Entry
}

type cloudManager struct {
	httpClient  http.Client
	managerAddr string
	logger      *logrus.Entry
}

type Job struct {
	Status    string `json:"status"`
	OutputURL string `json:"output_url"`
	Profile   string `json:"profile"`
}

func NewCloudManager(c CloudManagerConfig) *cloudManager {
	return &cloudManager{
		httpClient:  http.Client{Timeout: time.Second * 5},
		managerAddr: c.ManagerAddr,
		logger:      c.Logger.WithField("component", "cloud"),
	}
}

func (c *cloudManager) CreateJob(streamID *big.Int, address string) (string, error) {
	addr := fmt.Sprintf("%s/api/v1/job", c.managerAddr)

	jobRequest := &proto.AddJobRequest{
		StreamId:      streamID.Int64(),
		WalletAddress: address,
		ProfileId:     1,
	}

	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(jobRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, addr, buff)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("create job request failed with bad %d status", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	jobResponse := new(proto.AddJobResponse)
	err = json.Unmarshal(body, jobResponse)
	if err != nil {
		return "", err
	}

	return jobResponse.RtmpInputUrl, nil
}

func (c *cloudManager) GetJob(streamID *big.Int) (*Job, error) {
	addr := fmt.Sprintf("%s/api/v1/stream/%s", c.managerAddr, streamID.String())

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("get job request failed with bad %d status", res.StatusCode)
	}

	defer res.Body.Close()

	jsonBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	job := new(Job)
	err = json.Unmarshal(jsonBody, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (c *cloudManager) UpdateJobContractAddress(streamID *big.Int, contractAddress string) error {
	addr := fmt.Sprintf("%s/api/v1/contract_address/%s/%s", c.managerAddr, streamID.String(), contractAddress)

	req, err := http.NewRequest(http.MethodPost, addr, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("update job contract address request failed with bad %d status", res.StatusCode)
	}

	defer res.Body.Close()

	return nil
}

func (c *cloudManager) CancelJob(streamID *big.Int) error {
	addr := fmt.Sprintf("%s/api/v1/stream/stop/%s", c.managerAddr, streamID.String())

	req, err := http.NewRequest(http.MethodPost, addr, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("cancel job request failed with bad %d status", res.StatusCode)
	}

	defer res.Body.Close()

	return nil
}

func (c *cloudManager) AwaitJobStatus(streamID *big.Int, status proto.WorkOrderStatus) (*Job, error) {
	for timeout := time.After(time.Minute); ; {
		select {
		case <-timeout:
			return nil, errors.New("request timed out")
		default:
			job, err := c.GetJob(streamID)
			if err != nil {
				return nil, err
			}

			c.logger.Infof("received a job status: %s", job.Status)

			if job.Status != status.String() {
				time.Sleep(5 * time.Second)
				continue
			}

			return job, nil
		}
	}
}

package hybrik

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	// ErrGopSizeNan is an error returned when the GopSize field of db.Preset is not a valid number
	ErrGopSizeNan = fmt.Errorf("bitrate non a number")
)

// ErrNotDeleted is the error returned when a job is not successfully deleted from Hybrik
type ErrNotDeleted struct {
	JobID string
}

func (e ErrNotDeleted) Error() string {
	return fmt.Sprintf("Job of ID %s not deleted", e.JobID)
}

type transcodeErrResult struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
	Error   string `json:"error"`
}

type JobStatus struct {
	Errors []interface{} `json:"errors"`
	Job    struct {
		ID             int         `json:"id"`
		IsAPIJob       int         `json:"is_api_job"`
		Priority       int         `json:"priority"`
		CreationTime   time.Time   `json:"creation_time"`
		ExpirationTime time.Time   `json:"expiration_time"`
		UserTag        interface{} `json:"user_tag"`
		Status         string      `json:"status"`
		RenderStatus   string      `json:"render_status"`
		TaskCount      int         `json:"task_count"`
		Progress       float64     `json:"progress"`
		Name           string      `json:"name"`
		FirstStarted   time.Time   `json:"first_started"`
		LastCompleted  time.Time   `json:"last_completed"`
	} `json:"job"`
	Tasks []struct {
		ID         int       `json:"id"`
		Priority   int       `json:"priority"`
		Name       string    `json:"name"`
		RetryCount int       `json:"retry_count"`
		Status     string    `json:"status"`
		Assigned   time.Time `json:"assigned"`
		Completed  time.Time `json:"completed"`
	} `json:"tasks"`
}

type JobInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Progress  int       `json:"progress"`
}

// QueueJob takes a hybrik job json and submits a new job
func (c *Client) QueueJob(jobJSON string) (string, error) {

	resp, err := c.client.CallAPI("POST", "/jobs", nil, strings.NewReader(jobJSON))
	if err != nil {
		return "", err
	}

	var ter transcodeErrResult
	err = json.Unmarshal([]byte(resp), &ter)
	if err == nil {
		if !ter.Success && ter.Error != "" {
			return "", fmt.Errorf("%s", ter.Error)
		}
	}

	var job struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal([]byte(resp), &job)
	if err != nil {
		return "", err
	}

	return job.ID, err
}

// StopJob takes a jobID and deletes is from Hybrik
func (c *Client) StopJob(jobID string) error {
	resp, err := c.client.CallAPI("PUT", fmt.Sprintf("/jobs/%s/stop", jobID), nil, nil)
	if err != nil {
		return err
	}
	var job struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal([]byte(resp), &job)
	if job.ID != jobID || err != nil {
		return ErrNotDeleted{JobID: jobID}
	}

	return nil
}

// GetJobInfo takes a jobID and obtains basic status details about the corresponding job
func (c *Client) GetJobInfo(jobID string) (JobInfo, error) {
	values := url.Values{}
	values.Add("fields[]", "id")
	values.Add("fields[]", "name")
	values.Add("fields[]", "progress")
	values.Add("fields[]", "status")
	values.Add("fields[]", "start_time")
	values.Add("fields[]", "end_time")

	resp, err := c.client.CallAPI("GET", fmt.Sprintf("/jobs/%s/info", jobID), values, nil)
	if err != nil {
		return JobInfo{}, err
	}

	var ji JobInfo
	err = json.Unmarshal([]byte(resp), &ji)
	if err != nil {
		return JobInfo{}, err
	}

	return ji, nil
}

package horde

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/jaszczw/fiber/pkg/redis"
)

type RequestError struct {
	// The error message for this status code.
	Message *string `json:"message,omitempty"`
}
type RequestResult struct {
	RequestStatusCheck
	Generations []Generation `json:"generations,omitempty"`
}

type Generation struct {
	WorkerID   string `json:"worker_id,omitempty"`   // Worker ID
	WorkerName string `json:"worker_name,omitempty"` // Worker Name
	Model      string `json:"model,omitempty"`       // Generation Model
	Img        string `json:"img,omitempty"`         // Generated Image
	Seed       string `json:"seed,omitempty"`        // Generation Seed
}

type RequestStatusCheck struct {
	Finished      int32   `json:"finished,omitempty"`
	Processing    int32   `json:"processing,omitempty"`
	Restarted     int32   `json:"restarted,omitempty"`
	Waiting       int32   `json:"waiting,omitempty"`
	Done          bool    `json:"done,omitempty"`
	Faulted       bool    `json:"faulted,omitempty"`
	WaitTime      int32   `json:"wait_time,omitempty"`
	QueuePosition int32   `json:"queue_position,omitempty"`
	Kudos         float64 `json:"kudos,omitempty"`
	IsPossible    bool    `json:"is_possible,omitempty"`
}

func CheckImageStatus(requestId string) (*RequestStatusCheck, error) {
	response, err := http.Get("https://stablehorde.net/api/v2/generate/check/" + requestId)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("%s failed: %s | status: %d", requestId, response.Body, response.StatusCode)
		var result RequestError
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error: %s", *result.Message)
	}

	var result RequestStatusCheck
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetImageResult(requestId string) (*RequestResult, error) {
	response, err := http.Get("https://stablehorde.net/api/v2/generate/status/" + requestId)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("%s failed: %s | status: %d", requestId, response.Body, response.StatusCode)
		var result RequestError
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error: %s", *result.Message)
	}

	var result RequestResult
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckImageStatusLoop(requestId string) error {
	for {
		result, err := CheckImageStatus(requestId)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return err
		}

		fmt.Printf("Looping: %s", requestId)
		jsonString, _ := json.Marshal(result)

		redis.RedisClient.Set(context.Background(), requestId, jsonString, time.Minute)
		if result.Done {
			fmt.Printf("Done: %s", requestId)
			return handleDone(requestId)
		}

		// Wait 2 seconds
		time.Sleep(2 * time.Second)
	}
}

func handleDone(requestId string) error {
	finalResult, err := GetImageResult(requestId)
	if err != nil {
		return err
	}

	if finalResult == nil {
		return fmt.Errorf("requestId: %s, is not a valid image", requestId)
	}

	jsonString, _ := json.Marshal(finalResult)

	redis.RedisClient.Set(context.Background(), requestId, jsonString, time.Minute)

	http.Get(os.Getenv("DONE_CALLBACK_URL") + requestId)

	return nil
}

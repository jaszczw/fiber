package horde

type ModelPayloadRootStable struct {
	SamplerName       *string  `json:"sampler_name,omitempty"`
	Toggles           []int    `json:"toggles,omitempty"`
	CfgScale          *float64 `json:"cfg_scale,omitempty"`
	DenoisingStrength *float64 `json:"denoising_strength,omitempty"`
	Seed              *string  `json:"seed,omitempty"`
	Height            *int     `json:"height,omitempty"`
	Width             *int     `json:"width,omitempty"`
	SeedVariation     *int     `json:"seed_variation,omitempty"`
	PostProcessing    []string `json:"post_processing,omitempty"`
	Karras            *bool    `json:"karras,omitempty"`
}

type ModelGenerationInputStable struct {
	*ModelPayloadRootStable
	Steps *int `json:"steps,omitempty"`
	N     *int `json:"n,omitempty"`
}

type GenerationInput struct {
	Prompt           *string                     `json:"prompt,omitempty"`
	Params           *ModelGenerationInputStable `json:"params,omitempty"`
	Nsfw             *bool                       `json:"nsfw,omitempty"`
	TrustedWorkers   *bool                       `json:"trusted_workers,omitempty"`
	CensorNsfw       *bool                       `json:"censor_nsfw,omitempty"`
	SourceImage      *string                     `json:"source_image,omitempty"`
	SourceProcessing *string                     `json:"source_processing,omitempty"`
	SourceMask       *string                     `json:"source_mask,omitempty"`
	R2               *bool                       `json:"r2,omitempty"`
	Workers          []string                    `json:"workers,omitempty"`
	Models           []string                    `json:"models,omitempty"`
}

type RequestAsync struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

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

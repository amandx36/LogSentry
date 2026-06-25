package config

type Config struct {
    LogFile    string `json:"log_file"`
    OutputDir  string `json:"output_dir"`
    BufferSize int    `json:"buffer_size"`
	InputDir   string  `json:"input_dir"`
}
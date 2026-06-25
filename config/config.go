package config

import (
	"encoding/json"
	"os"
)

type Config struct {
    InputDir    string `json:"input_dir"`
    OutputDir  string `json:"output_dir"`
    BufferSize int    `json:"buffer_size"`
	
}

func Loadconfig (path string)(Config , error ){
	var cfg Config;
	// opening json 
	file , err := os.Open(path)
	if err != nil{
		return Config{} , err 
	}

	defer file.Close()
	// decode the  json and put it into the 
	err = json.NewDecoder(file).Decode(&cfg)
	if err !=nil{
		return  Config{},err 
	}
	return cfg , nil;
}
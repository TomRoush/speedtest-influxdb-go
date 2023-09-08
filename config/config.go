package config

import (
    "os"

    "gopkg.in/yaml.v3"

    "github.com/TomRoush/speedtest-influxdb-go/log"
)

type Config struct {
    General struct {
        Delay int `yaml:"delay"`
    } `yaml:"general"`
    Logging struct {
        Level string `yaml:"level"`
    } `yaml:"logging"`
    Speedtest struct {
        Server int `yaml:"server"`
    } `yaml:"speedtest"`
    InfluxDB struct {
        Address string `yaml:"address"`
        Port int `yaml:"port"`
        Organization string `yaml:"organization"`
        Bucket string `yaml:"bucket"`
        Token string `yaml:"token"`
        SSL bool `yaml:"ssl"`
    } `yaml:"influxdb"`
}

func ReadConfig() Config {
    var config Config
    ReadFile(&config)
    return config
}

func ReadFile(config *Config) {
    f, err := os.Open("config.yaml")
    HandleError(err)
    defer f.Close()

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(config)
    HandleError(err)
}

func HandleError(err error) {
    if err != nil {
        log.Error("%s", err)
        panic(err)
    }
}

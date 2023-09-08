package main

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	influxdb2write "github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/showwin/speedtest-go/speedtest"

	"github.com/TomRoush/speedtest-influxdb-go/config"
	"github.com/TomRoush/speedtest-influxdb-go/log"
)

func main() {
	log.Debug("Starting application")
	cfg := config.ReadConfig()
	log.SetLevel(log.ParseLevel(cfg.Logging.Level))

	time.Sleep(60 * time.Second)
	speedtestClient := speedtest.New()
	client, writeAPI := GetInfluxConnection(cfg)

	for {
		s := RunSpeedtest(speedtestClient, cfg)
		p := FormatResults(s)
		WriteInfluxData(writeAPI, p)

		if cfg.General.Delay == 0 {
			client.Close()
			return
		}
		sleepTime := time.Duration(cfg.General.Delay) * time.Second
		log.Info("Waiting %s seconds until next test", sleepTime)
		time.Sleep(sleepTime)
	}
}

func GetInfluxConnection(cfg config.Config) (influxdb2.Client, influxdb2api.WriteAPIBlocking) {
	protocol := "http"
	if cfg.InfluxDB.SSL {
		protocol += "s"
	}
	serverURL := fmt.Sprintf("%s://%s:%d", protocol, cfg.InfluxDB.Address, cfg.InfluxDB.Port)
	client := influxdb2.NewClient(serverURL, cfg.InfluxDB.Token)

	log.Debug("Testing connection to InfluxDB")
	_, err := client.Health(context.Background())
	HandleError(err)
	log.Debug("Successful connection to InfluxDB")

	writeAPI := client.WriteAPIBlocking(cfg.InfluxDB.Organization, cfg.InfluxDB.Bucket)

	return client, writeAPI
}

func SetupSpeedtest(speedtestClient *speedtest.Speedtest, cfg config.Config) *speedtest.Server {
	log.Debug("Setting up SpeedTest.net client")

	serverID := []int{}
	if cfg.Speedtest.Server != 0 {
		serverID = append(serverID, cfg.Speedtest.Server)
	}

	log.Debug("Picking the closest server")

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer(serverID)
	s := targets[0]

	log.Info("Selected server %s in %s at distance %.0fkm", s.ID, s.Name, s.Distance)

	return s
}

func RunSpeedtest(speedtestClient *speedtest.Speedtest, cfg config.Config) *speedtest.Server {
	s := SetupSpeedtest(speedtestClient, cfg)

	log.Info("Starting download test")
	s.DownloadTest()
	log.Info("Starting upload test")
	s.UploadTest()

	log.Info("Download: %.2fMbps - Upload: %.2fMbps - Latency: %s", s.DLSpeed, s.ULSpeed, s.Latency)

	return s
}

func FormatResults(s *speedtest.Server) *influxdb2write.Point {
	return influxdb2.NewPoint(
		"speed_test_results",
		map[string]string{
			"server":         s.ID,
			"server_name":    s.Name,
			"server_country": s.Country,
		},
		map[string]interface{}{
			"download": s.DLSpeed * 1000000.0,
			"upload":   s.ULSpeed * 1000000.0,
			"ping":     float64(s.Latency.Microseconds()) / 1000.0,
		},
		time.Now(),
	)
}

func WriteInfluxData(writeAPI influxdb2api.WriteAPIBlocking, point *influxdb2write.Point) {
	err := writeAPI.WritePoint(context.Background(), point)
	HandleError(err)
	log.Debug("Data written to InfluxDB")
}

func HandleError(err error) {
	if err != nil {
		log.Error("%s", err)
		panic(err)
	}
}

package main

type Config struct {
	HTTP struct {
		Port          int    `json:"port"`
		HttpsKeyFile  string `json:"https_key_file"`
		HttpsCertFile string `json:"https_cert_file"`
	} `json:"http"`
	Database struct {
		ConnectURI string `json:"connect_uri"`
		DBName     string `json:"db_name"`
	} `json:"database"`
}

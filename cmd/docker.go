package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"Qpay/config"
)

// seedCmd represents the seed command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "create docker-compose and doc_config",
	Long: `create docker-compose and set up config for using docker 
		create docker composer for postgres, alpine and build, after running this, you can up this docker compose and start the project`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerize(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}

func dockerize(configPath string) {
	var cfg *config.Config
	if _, err := os.Stat(configPath); err != nil {
		cfg = config.InitConfig("sample_config.yaml")
	} else {
		cfg = config.InitConfig(configPath)
	}
	dockerConfig := fmt.Sprintf(`database:
  host: qpay-postgres
  port: 5432
  username: %s
  password: %s
  db: %s

server:
  port: 80

auth:
  secret-key: %s
  expiresIn: %v

admin:
  name: %s
  username: %s
  email: %s
  password: %s`,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DB,
		cfg.JWT.SecretKey,
		int(cfg.JWT.ExpirationTime),
		cfg.Admin.Name,
		cfg.Admin.Username,
		cfg.Admin.Email,
		cfg.Admin.Password,
	)
	err := ioutil.WriteFile("doc_config.yaml", []byte(dockerConfig), 0644)
	if err != nil {
		log.Fatal(err)
	}
	dockerCompose := fmt.Sprintf(`version: "3.7"

services:
  # -----------------------------
  # qpay application
  # -----------------------------
  qpay-app:
    container_name: qpay-api
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped
    volumes:
      - ./doc_config.yaml:/app/config.yaml
    ports:
      - %v:80
    networks:
      - qpay
    depends_on:
      - postgres
    command: ["sh", "-c", "./Qpay migrate && ./Qpay seed && ./Qpay"]


  # -----------------------------
  # postgres database
  # -----------------------------
  postgres:
    container_name: qpay-postgres
    image: postgres:13.3
    restart: unless-stopped
    volumes:
      - qpay:/var/lib/postgresql/data
    ports:
      - %v:5432
    environment:
      - POSTGRES_PASSWORD=%s
      - POSTGRES_USER=%s
      - POSTGRES_DB=%s
      - TZ=Asia/Tehran

    networks:
      - qpay

# -----------------------------
# networks
# -----------------------------
networks:
  qpay:
    external: true

# -----------------------------
# volumes
# -----------------------------
volumes:
  qpay:
    name: qpay
    driver: local
`,
		cfg.Server.Port,
		cfg.Database.Port,
		cfg.Database.Password,
		cfg.Database.Username,
		cfg.Database.DB,
	)

	err = ioutil.WriteFile("docker-compose.yml", []byte(dockerCompose), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

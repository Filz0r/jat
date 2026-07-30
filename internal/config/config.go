package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/services"
	"github.com/filz0r/jat/internal/utils"
	"gorm.io/gorm"
)

type JatMode string

const (
	StandaloneMode JatMode = "standalone"
	ServerMode     JatMode = "server"
	ClientMode     JatMode = "client"
)

func (m JatMode) Valid() bool {
	switch m {
	case StandaloneMode, ServerMode, ClientMode:
		return true
	}
	return false
}

func (m JatMode) String() string {
	return string(m)
}

type ConfigFile struct {
	initialized  bool
	mode         JatMode
	serverURL    *string
	serverPort   *string
	dbUri        *string
	userID       *string
	userToken    *string
	refreshToken *string
	db           *gorm.DB
	Services     *services.ServiceManager
	SecretJWT    *string
}

type rawFile struct {
	Initialized  bool    `json:"initialized"`
	Mode         JatMode `json:"mode"`
	ServerURL    *string `json:"server_url,omitempty"`
	DbUri        *string `json:"db_uri,omitempty"`
	UserToken    *string `json:"user_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	UserID       *string `json:"user_id,omitempty"`
	ServerPort   *string `json:"server_port,omitempty"`
	SecretJWT    *string `json:"secret_jwt,omitempty"`
}

func (c *ConfigFile) getConfigFilePath() (string, error) {
	configFilePath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFilePath = filepath.Join(configFilePath, "jat")
	configFilePath = filepath.Join(configFilePath, "config.json")
	return configFilePath, nil
}

func (c *ConfigFile) exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func (c *ConfigFile) create(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer utils.CloseFile(file)

	raw := rawFile{
		Initialized: false,
		// this will become selectable in the future
		Mode:         StandaloneMode,
		ServerURL:    nil,
		DbUri:        c.dbUri,
		UserToken:    nil,
		RefreshToken: nil,
		UserID:       nil,
		ServerPort:   c.serverPort,
		SecretJWT:    c.SecretJWT,
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	c.initialized = raw.Initialized
	c.mode = raw.Mode
	c.serverURL = raw.ServerURL
	c.dbUri = raw.DbUri
	c.userToken = raw.UserToken
	c.refreshToken = raw.RefreshToken
	return nil
}

func (c *ConfigFile) load(path string) error {
	file, err := os.ReadFile(path)
	if len(file) == 0 {
		err = c.create(path)
		if err != nil {
			return err
		}
		file, err = os.ReadFile(path)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	raw := rawFile{}
	err = json.Unmarshal(file, &raw)
	if err != nil {
		return err
	}
	c.initialized = raw.Initialized
	c.mode = raw.Mode
	c.serverURL = raw.ServerURL
	c.dbUri = raw.DbUri
	c.userToken = raw.UserToken
	c.refreshToken = raw.RefreshToken
	c.userID = raw.UserID
	c.serverPort = raw.ServerPort
	c.SecretJWT = raw.SecretJWT
	return nil
}

func (c *ConfigFile) LoadData() error {
	path, err := c.getConfigFilePath()
	if err != nil {
		return err
	}
	// check if the file exists
	if !c.exists(path) {
		err = c.create(path)
		if err != nil {
			return err
		}
	} else {
		err = c.load(path)
		if err != nil {
			return err
		}
	}
	if !c.mode.Valid() {
		return fmt.Errorf("invalid mode was read in the config file: %s", c.mode)
	}
	return nil
}

func (c *ConfigFile) IsInitialized() bool {

	return c.initialized

}

func (c *ConfigFile) IsServer() bool {
	return c.mode == ServerMode
}

func (c *ConfigFile) IsClient() bool {
	return c.mode == ClientMode
}

func (c *ConfigFile) IsStandalone() bool {
	return c.mode == StandaloneMode
}

func (c *ConfigFile) Update() error {
	err := c.writeToDisk()
	if err != nil {
		return err
	}
	path, err := c.getConfigFilePath()
	if err != nil {
		return err
	}
	err = c.load(path)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigFile) SetInitialized(initialized bool) {
	if c.initialized {
		return
	}
	c.initialized = initialized
}

func (c *ConfigFile) SetDbUri(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	// TODO: This might require a refactor
	// TODO: in the future in case users want to use a ssl encrypted DB
	q := u.Query()
	if q.Get("sslmode") == "" {
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()
	}
	parsed := u.String()
	c.dbUri = &parsed
	return nil
}

func (c *ConfigFile) UnsetDbUri() {
	c.dbUri = nil
}

func (c *ConfigFile) UnsetServerUrl() {
	c.serverURL = nil
}

func (c *ConfigFile) writeToDisk() error {
	path, err := c.getConfigFilePath()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer utils.CloseFile(file)

	raw := rawFile{
		Initialized:  c.initialized,
		Mode:         c.mode,
		ServerURL:    c.serverURL,
		DbUri:        c.dbUri,
		UserToken:    c.userToken,
		RefreshToken: c.refreshToken,
		UserID:       c.userID,
		ServerPort:   c.serverPort,
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	nb, err := file.Write(data)
	if len(data) != nb {
		return fmt.Errorf("could not write entire file to disk")
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigFile) Mode() JatMode {
	return c.mode
}

func (c *ConfigFile) DbUri() string {
	if c.dbUri != nil {
		return *c.dbUri
	}
	return ""
}

func (c *ConfigFile) ServerURL() string {
	if c.serverURL != nil {
		return *c.serverURL
	}
	return ""
}

func (c *ConfigFile) SetServerURL(serverURL string) {
	c.serverURL = &serverURL
}

func (c *ConfigFile) SetMode(mode JatMode) error {
	if !mode.Valid() {
		return fmt.Errorf("invalid mode: %s", mode)
	}
	if c.mode == mode {
		return nil
	}
	c.mode = mode
	switch c.mode {
	case ClientMode:
		c.dbUri = nil
	case StandaloneMode:
		c.serverURL = nil
	}
	return nil
}

func (c *ConfigFile) GetDB() *gorm.DB {
	return c.db
}

func (c *ConfigFile) SetDB(db *gorm.DB) {
	c.db = db
}

func (c *ConfigFile) GetUserID() (string, error) {
	if c.userID == nil {
		return "", errors.New("user is not initialized")
	}
	return *c.userID, nil
}

func (c *ConfigFile) SetUserID(userID string) {
	c.userID = &userID
}

func (c *ConfigFile) LoadFromEnv() error {
	dbUri := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	secretJWT := os.Getenv("SECRET_JWT")

	if dbUri == "" || port == "" || secretJWT == "" {
		return fmt.Errorf("DATABASE_URL, SECRET_JWT and PORT must be set")
	}
	c.dbUri = &dbUri
	c.serverPort = &port
	c.SecretJWT = &secretJWT
	devMode := os.Getenv("JAT_DEV") != ""
	db, err := database.ConnectDb(*c.dbUri, devMode)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *ConfigFile) GetPort() (string, error) {
	if c.serverPort == nil {
		return "", errors.New("port is not initialized")
	}
	return *c.serverPort, nil
}

func New() *ConfigFile {
	return &ConfigFile{}
}

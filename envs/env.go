package envs

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Env struct {
	jWTToken        string
	jWTRefreshToken string
	dbHost          string
	dbPassword      string
	dbUser          string
	dbPort          string
	dbName          string
	redisHost       string
	redisPort       string
}

func newEnv() *Env {
	loadEnv()
	return &Env{
		jWTToken:        os.Getenv("JWT_SECRET"),
		jWTRefreshToken: os.Getenv("JWT_REFRESH_SECRET"),
		dbHost:          os.Getenv("DB_HOST"),
		dbPassword:      os.Getenv("DB_PASSWORD"),
		dbUser:          os.Getenv("DB_USERNAME"),
		dbPort:          os.Getenv("DB_PORT"),
		dbName:          os.Getenv("DB_DATABASE"),
		redisHost:       os.Getenv("REDIS_HOST"),
		redisPort:       os.Getenv("REDIS_PORT"),
	}
}

func (env *Env) GetJWTToken() string {
	return env.jWTToken
}

func (env *Env) GetJWTRefreshToken() string {
	return env.jWTRefreshToken
}

func (env *Env) GetDbHost() string {
	return env.dbHost
}

func (env *Env) GetDbPassword() string {
	return env.dbPassword
}

func (env *Env) GetDbUser() string {
	return env.dbUser
}

func (env *Env) GetDbPort() string {
	return env.dbPort
}

func (env *Env) GetDbName() string {
	return env.dbName
}

func (env *Env) GetRedisHost() string {
	return env.redisHost
}

func (env *Env) GetRedisPort() string {
	return env.redisPort
}

var instance *Env

func GetInstance() *Env {
	if instance == nil {
		instance = newEnv()
	}
	return instance
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		execPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		envPath := filepath.Join(filepath.Dir(execPath), ".env")
		err = godotenv.Load(envPath)
		if err != nil {
			panic(err)
		}
	}
}

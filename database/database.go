package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"shopping_cart/dtobjects"

	"gorm.io/gorm/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func Connection() *gorm.DB {

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	// DB_USER := os.Getenv("DB_USER")
	// DB_PASS := os.Getenv("DB_PASS")

	dbCreds := Conn()
	// DB_NAME := dbCreds.BbInstanceIdentifier
	DB_USER := dbCreds.Username
	DB_PASS := dbCreds.Password
	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println("dns==========>", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("#%v", db.Logger)
	//db.Debug()
	DB = db
	return DB
}

func GetConnection() *gorm.DB {
	return DB
}

func Conn() dtobjects.DBCredentials {
	secretName := os.Getenv("SecretName")
	region := "ap-south-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	var dbCredentials dtobjects.DBCredentials

	err = json.Unmarshal([]byte(secretString), &dbCredentials)
	if err != nil {
		log.Panic("parse error", err)
	}
	log.Printf("users creds===>%#v", dbCredentials)
	return dbCredentials

	// Your code goes here.
}

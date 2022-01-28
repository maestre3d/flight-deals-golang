package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
	"sync"
)

var (
	snsClientOnce = sync.Once{}
	snsClient *sns.Client
)

func init() {
	snsClientOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Print(err.Error())
			return
		}

		snsClient = sns.NewFromConfig(cfg)
	})
}

type SmsNotifier interface {
	SendSMS(message string) error
}

type AmazonSmsNotifier struct {
	client *sns.Client
}

var _ SmsNotifier = AmazonSmsNotifier{}

func (a AmazonSmsNotifier) SendSMS(message string) error {
	_, err := snsClient.Publish(context.Background(), &sns.PublishInput{
		Message:                aws.String(message),
		MessageAttributes:      nil,
		MessageDeduplicationId: nil,
		MessageGroupId:         nil,
		MessageStructure:       nil,
		PhoneNumber:            nil,
		Subject:                nil,
		TargetArn:              nil,
		TopicArn:               aws.String("arn:aws:sns:us-east-1:228850758643:flight-deal-found"),
	})
	return err
}

func notifyIfLowerPriceFound(notifier SmsNotifier, currentPrice, newPrice float64, destination string) {
	if currentPrice < newPrice {
		return
	}

	message := fmt.Sprintf("New lower price found (from %f to %f) for destionation %s",
		currentPrice, newPrice, destination)
	log.Print(message)

	if err := notifier.SendSMS(message); err != nil {
		log.Print(err.Error())
		return
	}
}

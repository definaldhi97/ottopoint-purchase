package publisher

import (
	"errors"
	"fmt"

	"github.com/benmanns/goworker"
	jsoniter "github.com/json-iterator/go"
)

/*
Client Worker Redis
*/
func sendRedisQueNameJobs(jobname, quename string, args ...interface{}) error {

	//1501575271.516704 [0 [::1]:63077] "RPUSH" "resque:queue:que1" "{\"class\":\"job1\",\"args\":[\"DMG ke [0] time [2017-08-01 15:14:31.515741102 +0700 WIB]\"]}"
	//jobName := jobname
	//queueName := quename
	settings := goworker.WorkerSettings{
		URI:            redis.GetRedisUri(), // "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{quename, "json", "queues"}, //[]string{quename, "delimiter", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
	err := goworker.Enqueue(&goworker.Job{
		Queue: quename,
		Payload: goworker.Payload{
			Class: jobname,
			Args:  args, //[]interface{}{userId, filename},
		},
	})
	if err != nil {

		fmt.Println("[ERROR-QUEUE]")
		fmt.Println("[publisher]")
		fmt.Println("[sendRedisQueNameJobs]")
		fmt.Println(fmt.Sprintf("Error while enqueue %s", err))

	} else {

		fmt.Println("[REDIS-SEND-SUCCESS]")
		fmt.Println("[publisher]")
		fmt.Println("[sendRedisQueNameJobs]")
		fmt.Println(fmt.Sprintf("Sukses send Redis args %s \n", args))

	}
	return err
}

// Pub ..
func Pub(data interface{}, jobname, quename string) error {

	fmt.Println("[Redis Queue Publish]")
	fmt.Println("[publisher]")
	fmt.Println("[Pub]")
	fmt.Println(fmt.Sprintf("Publish to redis queue: %v - %v ", jobname, quename))

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := json.Marshal(data)
	if err := sendRedisQueNameJobs(jobname, quename, string(b)); err != nil {
		//if err := services.SendRedisQueNameJobs(jobname, quename, fmt.Sprintf("DMG ke [%d] time [%v]", i, time.Now())); err != nil {

		fmt.Println("[Error]")
		fmt.Println("[publisher]")
		fmt.Println("[Pub]")
		fmt.Println(fmt.Sprintf("Error : %v", err))

		return errors.New("Error publish message")
	}

	fmt.Println("[Send Queue]")
	fmt.Println("[publisher]")
	fmt.Println("[Pub]")
	fmt.Println("Send queue No Error")

	return nil
}

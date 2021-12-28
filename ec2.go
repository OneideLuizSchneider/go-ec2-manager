package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2_machine struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Ec2VmsHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, GetVMs())
}

func Ec2VmsStartHandler(w http.ResponseWriter, r *http.Request) {
	vmid, ok := r.URL.Query()["vm-id"]
	if !ok || len(vmid[0]) < 1 {
		log.Println("Url Param 'vm-id' is missing")
		respondWithError(w, 500, "Url Param 'vm-id' is missing")
		return
	}

	vm_id := vmid[0]

	log.Println("Starting Session...")
	sess, err := session.NewSession(&aws.Config{})

	svc := ec2.New(sess)
	log.Println("Session Started...")

	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(vm_id),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = svc.StartInstances(input)
		if err != nil {
			log.Println("Error", err)
			respondWithError(w, 500, "error starting 2")
			return
		} else {
			log.Println("Success", result.StartingInstances)
		}
	} else {
		log.Println("Error", err)
		respondWithError(w, 500, "error starting 1")
		return
	}

	respondWithJson(w, 200, "Starting...")
}

func Ec2VmsStopHandler(w http.ResponseWriter, r *http.Request) {
	vmid, ok := r.URL.Query()["vm-id"]
	if !ok || len(vmid[0]) < 1 {
		log.Println("Url Param 'vm-id' is missing")
		respondWithError(w, 500, "Url Param 'vm-id' is missing")
		return
	}

	vm_id := vmid[0]

	log.Println("Starting Session...")
	sess, err := session.NewSession(&aws.Config{})
	svc := ec2.New(sess)
	log.Println("Session Started...")

	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(vm_id),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = svc.StopInstances(input)
		if err != nil {
			log.Println("Error", err)
			respondWithError(w, 500, "error stopping 2")
			return
		} else {
			log.Println("Success", result.StoppingInstances)
		}
	} else {
		log.Println("Error", err)
		respondWithError(w, 500, "error stopping 1")
		return
	}

	respondWithJson(w, 200, "Stopping...")
}

func GetVMs() []ec2_machine {
	log.Println("Starting Session...GetVMs")
	sess, err := session.NewSession(&aws.Config{
		//Region: aws.String("sa-east-1"),
		//Credentials: credentials.NewStaticCredentials(aws_user, aws_pass, ""),
	})
	client := ec2.New(sess)
	log.Println("Session Started...")

	log.Println("Getting VMs...")
	result, err := client.DescribeInstances(nil)
	if err != nil {
		log.Println("ERROR Getting VMs...")
		return []ec2_machine{}
	}

	log.Println("Loading VMs into Structs...")
	listVms := []ec2_machine{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			p := ec2_machine{
				ID:     *instance.InstanceId,
				Name:   *instance.Tags[0].Value,
				Status: *instance.State.Name,
			}
			listVms = append(listVms, p)
		}
	}
	log.Println("Returning VMs...")
	return listVms
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, "ec2 manager :D")
}

func main() {
	http.HandleFunc("/", DefaultHandler)
	http.HandleFunc("/ec2-vms", Ec2VmsHandler)
	http.HandleFunc("/ec2-vms/stop", Ec2VmsStopHandler)
	http.HandleFunc("/ec2-vms/start", Ec2VmsStartHandler)
	http.ListenAndServe(":8080", nil)
}

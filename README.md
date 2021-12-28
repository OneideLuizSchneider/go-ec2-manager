## GO AWS EC2 Manager

### Steps

- Export the AWS env variables: 
  - `export AWS_REGION=sa-east-1`
  - `export AWS_ACCESS_KEY_ID=....`
  - `export AWS_SECRET_ACCESS_KEY=....`
  - More details here -> <https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html>
- To run it: 
  - With go: 
    - `go get`
    - `go run ec2.go`
  - With docker: `docker-compose up --build`

### Endpoints

- `/ec2-vms` return all the ec2 VMs
- `/ec2-vms/stop?vm-id=123` stop an ec2 VM and return the current status on AWS
- `/ec2-vms/start?vm-id=123` start an ec2 VM and return the current status on AWS

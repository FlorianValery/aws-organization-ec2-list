// Cross account logic, forked from https://maori.geek.nz/assuming-roles-in-aws-with-go-aeeb28fab418

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Struct to store the session with custom parameters
type Clients struct {
	session *session.Session
	configs map[string]*aws.Config
}

// Func to start a session
func (c Clients) Session() *session.Session {
	if c.session != nil {
		return c.session
	}
	sess := session.Must(session.NewSession())
	c.session = sess
	return sess
}

// Custom config func
func (c Clients) Config(
	region *string,
	account_id *string,
	role *string) *aws.Config {

	// return no config for nil inputs
	if account_id == nil || region == nil || role == nil {
		return nil
	}
	arn := fmt.Sprintf(
		"arn:aws:iam::%v:role/%v",
		*account_id,
		*role,
	)
	// include region in cache key otherwise concurrency errors
	key := fmt.Sprintf("%v::%v", *region, arn)

	// check for cached config
	if c.configs != nil && c.configs[key] != nil {
		return c.configs[key]
	}
	// new creds
	creds := stscreds.NewCredentials(c.Session(), arn)
	// new config
	config := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(*region).
		WithMaxRetries(10)
	if c.configs == nil {
		c.configs = map[string]*aws.Config{}
	}
	c.configs[key] = config
	return config
}

// Create EC2 client
func (c *Clients) EC2(
	region string,
	account_id string,
	role string) *ec2.EC2 {
	return ec2.New(c.Session(), c.Config(&region, &account_id, &role))
}

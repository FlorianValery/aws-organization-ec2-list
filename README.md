# aws-organization-ec2-list

Custom script that allows you to query all the EC2 instances within an AWS organization.

## Requirements
In order to make the API calls in each account, you will need a role that your user can assume deployed at the organization level -- e.g. OrganizationEc2ReadRole. 
This role must have read permissions for EC2:
* AWS Managed policy
```
arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess
```
* Custom policy
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "ec2:Describe*",
            "Resource": "*"
        },
    ]
}
```
Additionally, this same role will need organization permissions in the master Organization account, so we are able to automatically retrieve all the accounts IDs within our Org. This will be useful as our organization might grow in the future. You can use the AWS managed policy:
```
arn:aws:iam::aws:policy/AWSOrganizationsReadOnlyAccess
```

## Usage
* Clone the repo locally
```
git clone https://github.com/FlorianValery/aws-organization-ec2-list.git
```
* Update the config/default.json file with the proper region and cross-accounts role that your user can assume
```
{
  "Region": "us-east-1",
  "Organization_Role": "OrganizationEc2ReadRole"
}
```
* Export your AWS credentials using the CLI or tools like Awsume
```
awsume master-role
```
* Build the script package and run it
```
go build -o app
./app
```

**Output example**
```
Account Name,Account ID,Instance Name,Instance Size,Instance ID,Image ID,Platform,Private IP,State,Timestamp
account-prod,000000000000,awesome_app,t3.micro,i-00000000000000aa,ami-000000000000aa,linux,10.0.0.1,running,2021-01-18 00:00:00 +0000 UTC
account-prod,000000000000,awesome_app,t3.micro,i-11111111111111bb,ami-11111111111111bb,linux,10.100.0.0,running,2021-12-18 100:00:00 +0000 UTC
account-staging,111111111111,awesome_db,t3.micro,i-11111111111111cc,ami-11111111111111cc,linux,10.10.0.1,running,2021-01-18 00:00:00 +0000 UTC
[...]
```
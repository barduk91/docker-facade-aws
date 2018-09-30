## Docker Facade AWS

Provide a fully functional local AWS stack using docker. Replicate all your cloud environment in your local machine. Provide your CloudFormation file with all your infrastructure and it will be recreate in local. It based on  [localstack](https://github.com/localstack/localstack "localstack") project.

### Prerequisites
* Docker should be installed on your local (because it's the core of this project, your AWS will be on a docker container)

*  [AWS Cli](https://aws.amazon.com/cli/?nc1=h_ls  "AWS Cli") must be installed and configured (not be afraid, none service will create anything in cloud)

### Develop

To compile project in local, also you need to install [dep](https://github.com/golang/dep  "dep") to ensure dependencies. 
After that you only need to run:
```
dep ensure
```

 and all dependencies should be downloaded. 


### How to start using it?
You just need execute
```
./main -f your-cloudformation-file.yaml
```

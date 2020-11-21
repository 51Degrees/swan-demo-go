# This is needed because the windows zip process used by EB will not enable
# the executable to be run on linux.
# https://forums.aws.amazon.com/message.jspa?messageID=825738#825738
# go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip

# Set up Go for AWS Elastic Beanstalk (EB) build
set GOPATH=%CD%
set GOARCH=amd64
set GOOS=linux

# Build the application
go build -o ./application ./src/server.go

# Create a zip file with the application and the settings file
bin\build-lambda-zip.exe -o aws-eb-swan-demo.zip application appsettings.json Procfile .ebextensions/.config images/190811762.jpeg images/221406343.jpeg images/234657570.jpeg
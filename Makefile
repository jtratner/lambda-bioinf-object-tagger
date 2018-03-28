NAME:=hpc-object-tagger
HANDLER_NAME?=$(NAME)
REGION?=region
ACCOUNT_ID?=whatever
ROLE?=arn:aws:iam::$(ACCOUNT_ID):service-role/$(NAME)
package:
	GOOS=linux GOARCH=amd64 go build -o $(NAME) -ldflags "-X main.Version=$$(git describe) -X main.GitCommit=$$(git rev-parse HEAD | cut -c 1-8)" .
	zip $(NAME).zip $(NAME)
	# unclear if handler is necessary
publish: package
	aws lambda create-function \
	--region $(REGION) \
	--function-name $(HANDLER_NAME) \
	--role $(ROLE) \
	--runtime go1.x \
	--zip-file fileb://$(PWD)/$(NAME).zip \
	--handler $(HANDLER_NAME)

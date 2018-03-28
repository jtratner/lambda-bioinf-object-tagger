NAME:=$(shell basename $(PWD))
HANDLER_NAME?=$(NAME)
REGION?=region
ACCOUNT_ID?=whatever
ROLE?=arn:aws:iam::$(ACCOUNT_ID):role/execution_role
package:
	GOOS=linux go build -ldflags "-X main.Version=$$(git describe) -X main.GitCommit=$$(git rev-parse HEAD | cut -c 1-8)" .
	zip $(NAME).zip $(NAME)
	# unclear if handler is necessary
	echo aws lambda create-function \
	--region $(REGION) \
	--function-name $(HANDLER_NAME) \
	--memory 128 \
	--role $(ROLE) \
	--runtime go1.x \
	--zip-file fileb://$(PWD)/$(NAME).zip \
	--handler $(HANDLER_NAME)

# Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GORUN=$(GOCMD) run
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get

    BINARY_NAME=websearch
    BINARY_UNIX=$(BINARY_NAME)_unix

    OUTPUT_FILE=output.xml
    XML_FILE=small.xml
    MAIN_NAME=engine.go
all:
	$(GORUN) $(MAIN_NAME) $(XML_FILE)
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(OUTPUT_FILE)
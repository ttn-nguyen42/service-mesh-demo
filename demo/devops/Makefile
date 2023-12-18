all: weather locations dashboard

weather: NAME=weather
weather: TAG=nguyentrantrung
weather: SOURCE_PATH="../weather"
weather: HELPER_PATH=../helper
weather: DOCKERFILE=Dockerfile
weather: 
	git pull
	cp $(DOCKERFILE) $(SOURCE_PATH)
	cd $(SOURCE_PATH) && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main main.go && \
	docker build -t $(TAG)/$(NAME):latest --platform linux/amd64 -f $(DOCKERFILE) . && \
	rm -rf main $(DOCKERFILE)
	docker push $(TAG)/$(NAME)

locations: NAME=locations
locations: TAG=nguyentrantrung
locations: SOURCE_PATH="../locations"
locations: HELPER_PATH="../helper"
locations: DOCKERFILE=Dockerfile
locations:
	git pull
	cp $(DOCKERFILE) $(SOURCE_PATH)
	cd $(SOURCE_PATH) && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main main.go && \
	docker build -t $(TAG)/$(NAME):latest --platform linux/amd64 -f $(DOCKERFILE) . && \
	rm -rf main $(DOCKERFILE)
	docker push $(TAG)/$(NAME)

dashboard: NAME=dashboard
dashboard: TAG=nguyentrantrung
dashboard: SOURCE_PATH="../dashboard"
dashboard: HELPER_PATH="../helper"
dashboard: DOCKERFILE=htmx.Dockerfile
dashboard:
	git pull
	cp $(DOCKERFILE) $(SOURCE_PATH)
	cd $(SOURCE_PATH) && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main main.go && \
	docker build -t $(TAG)/$(NAME):latest --platform linux/amd64 -f $(DOCKERFILE) . && \
	rm -rf main $(DOCKERFILE)
	docker push $(TAG)/$(NAME)
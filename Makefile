# Define the target executable
TARGET = remem_app

build:
	go build -o $(TARGET)

clean:
	rm -f $(TARGET)

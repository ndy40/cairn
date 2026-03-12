BINARY_NAME := cairn
INSTALL_DIR := /usr/local/bin
GO := /usr/local/go/bin/go

.PHONY: build install clean extension

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/cairn

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

clean:
	rm -f $(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

extension:
	cd vicinae-extension && npm run build

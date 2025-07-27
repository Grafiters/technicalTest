#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect the operating system
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "MacOS"
    elif [[ "$OSTYPE" == "cygwin" ]]; then
        echo "Cygwin"
    elif [[ "$OSTYPE" == "msys" ]]; then
        echo "Git Bash"
    elif [[ "$OSTYPE" == "win32" ]]; then
        echo "Windows"
    else
        echo "Unknown"
    fi
}

# Install Scoop on Windows
install_scoop() {
    if ! command_exists scoop; then
        echo "Installing Scoop..."
        powershell -Command "Set-ExecutionPolicy RemoteSigned -scope CurrentUser"
        powershell -Command "iwr -useb get.scoop.sh | iex"
    else
        echo "Scoop is already installed."
    fi

    echo "Updating Scoop..."
    scoop update

    echo "Installing Golang 1.20..."
    scoop install go@1.20

    echo "Setting Golang 1.20 as the current version..."
    scoop reset go@1.20
}

# Install Go and Migrate on Linux
install_linux_tools() {
    echo "Installing Golang 1.20 on Linux..."

    if ! command_exists go; then
        wget https://golang.org/dl/go1.20.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
        source ~/.profile
    else
        echo "Go is already installed."
    fi
}

# Install Go on MacOS
install_go_macos() {
    echo "Installing Golang 1.20 on macOS..."

    if ! command_exists go; then
        wget https://golang.org/dl/go1.20.darwin-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.20.darwin-amd64.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.zshrc
        source ~/.zshrc
    else
        echo "Go is already installed."
    fi
}

# Install Migrate tool using Go
install_migrate() {
    echo "Installing Migrate tool..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

    echo "Installation completed."

    # Adding Go bin directory to the PATH
    export PATH=$PATH:$(go env GOPATH)/bin
    echo "Migrate tool installed at $(go env GOPATH)/bin/migrate"
}

install_swag() {
    echo "Installing Swag for Swagger documentation..."

    # Determine the platform
    platform=$(uname)

    # Installing Swag for each platform
    if [[ "$platform" == "Linux" ]]; then
        echo "Detected Linux OS"
        go install github.com/swaggo/swag/cmd/swag@latest
        echo "Swag installed successfully for Linux"
    elif [[ "$platform" == "Darwin" ]]; then
        echo "Detected macOS"
        go install github.com/swaggo/swag/cmd/swag@latest
        echo "Swag installed successfully for macOS"
    elif [[ "$platform" == "MINGW"* ]] || [[ "$platform" == "CYGWIN"* ]] || [[ "$platform" == "MSYS"* ]]; then
        echo "Detected Windows OS (MINGW, CYGWIN, or MSYS)"
        go install github.com/swaggo/swag/cmd/swag@latest
        echo "Swag installed successfully for Windows"
    else
        echo "Unsupported platform: $platform"
        return 1
    fi

    # Verifying installation
    swag_version=$(swag --version 2>/dev/null)

    if [ $? -eq 0 ]; then
        echo "Swag version: $swag_version"
    else
        echo "Swag installation failed or not found in PATH."
        echo "Make sure your Go binary path is correctly set up in your environment variables."
    fi
}

# Main script
os_type=$(detect_os)

if [[ "$os_type" == "Windows" || "$os_type" == "Git Bash" || "$os_type" == "Cygwin" ]]; then
    install_scoop
    install_migrate
    install_swag
elif [[ "$os_type" == "Linux" ]]; then
    install_linux_tools
    install_migrate
    install_swag
elif [[ "$os_type" == "MacOS" ]]; then
    install_go_macos
    install_migrate
    install_swag
else
    echo "Unsupported OS: $os_type"
fi

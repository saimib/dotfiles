#!/bin/bash
# macOS setup script

echo "Starting macOS setup..."

# Check if Homebrew is already installed
if command -v brew &>/dev/null; then
    echo "Homebrew is already installed. Updating..."
    brew update
else
    echo "Installing Homebrew..."
    # Official Homebrew installation command
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Add Homebrew to PATH for the current session if needed
    if [[ $(uname -m) == "arm64" ]]; then
        # For Apple Silicon Macs
        echo "Configuring Homebrew for Apple Silicon Mac..."
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
        eval "$(/opt/homebrew/bin/brew shellenv)"
    else
        # For Intel Macs
        echo "Configuring Homebrew for Intel Mac..."
        echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
        eval "$(/usr/local/bin/brew shellenv)"
    fi
    
    # Verify installation
    brew doctor
fi

echo "Homebrew installation completed!"

# Check if Git is already installed
if command -v git &>/dev/null; then
    echo "Git is already installed. Version: $(git --version)"
else
    echo "Installing Git using Homebrew..."
    brew install git
    
    # Verify Git installation
    if command -v git &>/dev/null; then
        echo "Git installation completed. Version: $(git --version)"
    else
        echo "Failed to install Git. Please install it manually."
        exit 1
    fi
fi

# Install pyenv

if command -v pyenv &>/dev/null; then
    echo "pyenv is already installed. Version: $(pyenv --version)"
else
    echo "Installing pyenv using Homebrew..."
    brew install pyenv
    
    # Verify pyenv installation
    if command -v pyenv &>/dev/null; then
        echo "pyenv installation completed. Version: $(pyenv --version)"
    else
        echo "Failed to install pyenv. Please install it manually."
        exit 1
    fi
fi


# Install Gcloud CLI
# Download Google Cloud SDK
echo "Downloading Google Cloud SDK..."
curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-darwin-arm.tar.gz

# Extract the archive
echo "Extracting Google Cloud SDK..."
tar -xzf google-cloud-cli-darwin-arm.tar.gz

# Install the SDK
echo "Installing Google Cloud SDK..."
./google-cloud-sdk/install.sh

# Clean up
rm google-cloud-cli-darwin-arm.tar.gz



# Source .zshrc to apply changes
echo "Sourcing .zshrc to apply changes..."
source ~/.zshrc
echo "macOS setup completed successfully!"


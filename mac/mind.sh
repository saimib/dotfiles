#!/bin/bash
# This script installs the necessary dependencies for the mind repository.

echo "Installing mind dependencies..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Homebrew is not installed. Please run setup.sh first."
    exit 1
fi

if ! command -v pyenv &> /dev/null; then
    echo "pyenv is not installed. Please run setup.sh first."
    exit 1
fi

# Install gsutil
pyenv local 3.12.x
pip install gsutil


if ! command -v gcloud &> /dev/null; then
    echo "Google Cloud SDK is not installed. Please run setup.sh first."
else
    echo "Google Cloud SDK is already installed."
    gcloud init
fi


# Setup mind repository
git clone git@github.com:saimib/mind.git

cd mind && mkdir _portal

cp -p .hooks/pre-commit .git/hooks/pre-commit
cp -p .hooks/post-merge .git/hooks/post-merge


# setup gsutil
gsutil -m rsync gs://mib-personal-bucket/_portal _portal



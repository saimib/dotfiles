# set execution policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

Write-Host "Locating windows profile ... $PROFILE"

# Check for winprofile in repo
$winProfilePath = Join-Path -Path (Get-Location) -ChildPath "./winprofile.ps1"
if (-Not (Test-Path $winProfilePath)) {
  Write-Host "winprofile missing at $winProfilePath"
  exit 1
}

# Check for Profile, if it exists create a backup before removing it
if (Test-Path $PROFILE) {
  $backupPath = Join-Path -Path (Split-Path $PROFILE) -ChildPath "Microsoft.PowerShell_profile_backup.ps1"
  Copy-Item -Path $PROFILE -Destination $backupPath -Force
  Remove-Item -Path $PROFILE -Force
  Write-Host "Existing profile backed up to: $backupPath" -ForegroundColor Yellow
}

# Create a new profile
New-Item -ItemType File -Path $PROFILE -Force | Out-Null

# Copy contents from winprofile
if (Test-Path $winProfilePath){
  Get-Content $winProfilePath | Set-Content $PROFILE
  Write-Host "Profile successfully created"
} else {
  Write-Host "Winprofile not found at $winProfilePath"
}


Write-Host "Installing Scoop..."
if (Get-Command scoop -ErrorAction SilentlyContinue) {
  Write-Host "Scoop already installed"
} else {
  # Install scoop
  Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
}


Write-Host "Install basic tools..."
# Install basic stuff
scoop install neovim git ripgrep wget fd unzip gzip mingw make win32yank

# Install nerd fonts
scoop bucket add nerd-fonts
scoop install FiraCode-NF NerdFontsSymbolsOnly

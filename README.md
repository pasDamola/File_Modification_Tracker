# File Modification Tracker

## Overview

File Modification Tracker is a macOS application that monitors specified directories for file modifications and logs changes.

## Features

- Monitors specified directories for file changes.
- Displays logs in real-time.
- Start and stop monitoring through a user-friendly interface.

## Installation

1. Download the latest release from [here](https://github.com/pasDamola/File_Modification_Tracker/blob/main/File-Tracker.pkg).
2. Open the downloaded `.pkg` file.
3. Follow the installation instructions.

## Usage

1. Launch the application from `/Applications/FileModificationTracker`.
2. Click **Start Service** to begin monitoring.
3. Click **Stop Service** to halt monitoring.
4. Logs will be displayed in the application window.

## Requirements

- macOS version X or later.
- Go installed (for development).

## Testing

To run tests, navigate to the project directory and execute:

```bash
go test ./...
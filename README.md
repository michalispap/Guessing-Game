# Guessing Game 🎲

A terminal-based number guessing game.

## About

Guessing Game is an interactive command-line game where you try to guess a randomly generated number between 1 and 100. You have **10 attempts**.

## Getting Started

### Prerequisites

- Go 1.25.5 or higher

### Installation

1. Clone the repository:
```
git clone https://github.com/michalispap/Guessing-Game. git
cd Guessing-Game
```

2. Install dependencies:
```
go mod download
```

3. Run the game:
```
go run main.go
```

#### Docker (Optional)

In case you prefer to play the game using Docker, run the below:

```
docker pull michalispap99/guessing-game:latest # Pull the image from Docker Hub
docker run -it --rm guessing-game:latest  # Run the game
```


## Libraries Used

This project leverages the [Charm](https://charm.sh/) ecosystem:

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**
- **[Bubbles](https://github.com/charmbracelet/bubbles)**
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**

# Actor Model Benchmark

This benchmark application tests the performance of the actor model implementation with mock renderer and physics systems.

## Usage

```bash
# Build the benchmark
go build -o bench main.go

# Run with default settings (10 seconds, 64 Hz, 100 players)
./bench

# Run with custom duration (30 seconds)
./bench -duration=30s

# Run with 128 Hz tick rate
./bench -tickrate=128

# Run with 500 players
./bench -players=500

# Combine multiple options
./bench -duration=60s -tickrate=128 -players=1000
```

## Flags

- `-duration`: Benchmark duration (default: 10s)
- `-tickrate`: Tick rate in Hz, must be 64 or 128 (default: 64)
- `-players`: Number of players to simulate (default: 100)

## What it does

1. Creates a mock actor system with:
   - MockRenderer: Counts render update messages
   - MockPhysics: Counts physics update messages and simulates processing
   - MockCamera: Provides camera data
   - BenchmarkPlayers: Generate random movement data on each tick

2. Runs a tick loop that broadcasts movement events to all players

3. Measures and reports:
   - Total ticks processed
   - Total messages sent
   - Messages per second
   - Ticks per second
   - Messages per tick

## Example Output

```
=== BENCHMARK RESULTS ===
Duration: 10.123s
Tick Rate: 64 Hz
Number of Players: 100
Total Ticks: 648
Total Messages: 64800
Messages per second: 6398.45
Ticks per second: 63.98
Messages per tick: 100.00
========================
```

This benchmark helps evaluate the actor model's ability to handle high-frequency message broadcasting and processing. 
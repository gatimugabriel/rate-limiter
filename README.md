# Rate Limiter

A Go-based rate limiter library designed to help control traffic for APIs and web applications by limiting the number of requests per client over a set period. Each branch in this repository features a different rate-limiting algorithm, making it easy to choose and implement the best approach for your use case.

## Repository Structure

This repository is organized by branches, each containing the source code and specific README documentation for a different rate-limiting algorithm:

- **Main Branch** (`main`): The central branch with general information about the project and an overview of the available algorithms.
- **Token-Bucket Branch** (`token-bucket`): Implements the Token-Bucket rate-limiting algorithm, with detailed documentation on usage, setup, and testing.

### Available Algorithms (Branch Structure)
Each branch includes:
- **Source Code**: Implementation of the specific rate-limiting algorithm.
- **README**: Detailed documentation covering the algorithm’s functionality, usage, configuration, and testing.

#### Current Algorithms
1. **Token-Bucket Algorithm** ([`token-bucket` branch](https://github.com/gabrielgatimu/rate-limiter/tree/token-bucket))
    - Limits client requests based on a token-bucket strategy.
    - [Checkout the branch and view the README for more details](https://github.com/gabrielgatimu/rate-limiter/tree/token-bucket).

More algorithms, such as sliding-window and leaky-bucket rate limiting, will be added in their respective branches in the future.

## Installation

To add this library to your Go project:
```bash
go get github.com/gabrielgatimu/rate-limiter

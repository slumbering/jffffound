# jffffound.com


This go project aim to parse a [til repository*](https://github.com/slumbering/til) in order to generate a static website.

*til repository is a collection of "Today I Found" posts, usually written in markdown format.

## Installation

> Requirements:
- Go 1.16 or higher
- Air (for live reloading during development)
- Need to set a GITHUB_TOKEN environment variable for accessing GitHub API. [More info here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

To install the project, you need to have Go installed on your machine. Then you can run the following command to download the dependencies:

```bash
go mod tidy
```
To install Air for live reloading, you can use the following command:

```bash
go install github.com/air-verse/air@latest
```

To set the GITHUB_TOKEN environment variable, you can use the following command in your terminal:

```bash
export GITHUB_TOKEN=your_github_token_here
```

## Usage
To run the project, you can use the following command:

```bash
air
```

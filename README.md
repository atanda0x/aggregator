# Aggregator

Aggregator is a Go-based application designed to scrape RSS feeds and store the data in a PostgreSQL database using SQLC for type-safe queries. The project leverages Gin for building a robust API and includes features for user management, feed handling, and feed follow management.

## Features

- User management (create, read)
- RSS feed management (create, read)
- Feed follow management (create, read, delete)
- Scheduled scraping of RSS feeds
- Endpoint for getting posts for users

## Getting Started

### Prerequisites

- Go 1.16+
- PostgreSQL
- Docker (optional, for containerization)

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/atanda0x/aggregator.git
    cd aggregator
    ```

2. **Load environment variables:**
    Create a `.env` file in the project root and set the following variables:
    ```sh
    SERVER_ADDRESS=<your-port>
    DB_DRIVER=<your-driver(postgres)>
    DB_SOURCE=<your-database-url>
    ```

3. **Install dependencies:**
    ```sh
    go mod tidy
    ```

4. **Run the application:**
    ```sh
    go run main.go
    ```

## Usage

### API Endpoints

- **Health Check:**
    ```sh
    GET /healthz
    ```
- **User Management:**
    ```sh
    POST /users
    GET /users
    ```
- **Feed Management:**
    ```sh
    POST /feeds
    GET /feeds
    ```
- **Feed Follow Management:**
    ```sh
    POST /feed_follows
    GET /feed_follows
    DELETE /feed_follows/:feedFollowID
    ```
- **Posts for Users:**
    ```sh
    GET /posts
    ```

### Running the Scraper

The RSS feed scraper runs as a goroutine within the main function. It scrapes feeds every specified duration and processes them concurrently.

### Configuration

The configuration is loaded from the `.env` file. Ensure you have the necessary environment variables set for the application to function correctly.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for review.

## License

This project is licensed under the MIT License.

## Contact

For any inquiries or support, please contact [atanda0x](mailto:atanda0x@gmail.com).

# Shodan Proxy

[中文版 README](README_CN.md)

Shodan Proxy is a Go-based proxy server for Shodan API requests with additional features like IP filtering, path blocking, and an admin panel.

## Features

- Proxy Shodan API requests
- IP whitelisting
- Path blocking
- Admin panel for configuration
- Multiple Shodan API key management with round-robin rotation


## Installation

1. Clone the repository
2. Build the Docker image:
   ```
   docker compose build
   ```
3. Start the container:
   ```
   docker compose up -d
   ```

## Configuration

Edit the `config/config.yaml` file to set up:

- Blocked paths
- Allowed IPs
- Trusted proxies
- Admin credentials

## Usage

Access the admin panel at `http://localhost:8080/admin` to manage settings and API keys.

Default admin credentials:
- Username: admin
- Password: shodanproxy

**Note:** For security reasons, it is strongly recommended to change the default password after your first login.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

This project uses Go modules for dependency management. The `go.mod` file is generated during the Docker build process to ensure the latest compatible versions of dependencies are used. If you're developing outside of Docker, you can generate the `go.mod` file by running:

```
go mod init shodan-proxy
go mod tidy
```

## Contact

If you have any questions or suggestions, please open an issue or contact the project maintainer directly.

Thank you for your interest in the Shodan Proxy project!
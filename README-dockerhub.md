# Shodan Proxy Docker Image

Shodan Proxy is a Go-based proxy server for Shodan API requests with additional features like IP filtering, path blocking, and an admin panel.

## Features

- Proxy Shodan API requests
- IP whitelisting
- Path blocking
- Admin panel for configuration
- Multiple Shodan API key management with round-robin rotation
- Secure admin panel with authentication

## Usage

### Running the Container

To run the Shodan Proxy container:

```bash
docker run -p 8080:8080 liuweitao/shodan-proxy:latest
```

Access the admin panel at `http://localhost:8080/admin` to manage settings and API keys.

Default admin credentials:
- Username: admin
- Password: shodanproxy

## Configuration

The container is configured using a `config.yaml` file. To configure the container, you need to mount a volume containing your `config.yaml` file.

Example of mounting a configuration file:

```bash
docker run -p 8080:8080 -v /path/to/your/config:/app/config liuweitao/shodan-proxy:latest
```

Note: Ensure that your `config.yaml` file is properly secured, especially if it contains sensitive information like API keys.

## Docker Compose

To use Docker Compose for easy deployment, follow these steps:

1. Create and enter the project directory:
   ```bash
   mkdir shodan-proxy && cd shodan-proxy
   ```

2. Download the `compose.yaml` file:
   ```bash
   curl -O https://raw.githubusercontent.com/liuweitao/shodan-proxy/main/compose.yaml
   ```

3. Run the container:
   ```bash
   docker compose up -d
   ```

This will set up and start the Shodan Proxy container using the configuration specified in the `compose.yaml` file.

## API Call Examples

Here are some examples of Shodan API calls, comparing the official API and this proxy's usage:

1. Search for host information

   Official API:
   ```
   https://api.shodan.io/shodan/host/search?key=YOUR_API_KEY&query=apache
   ```

   This proxy:
   ```
   http://localhost:8080/shodan/host/search?query=apache
   ```

2. Get information for a specific IP

   Official API:
   ```
   https://api.shodan.io/shodan/host/1.1.1.1?key=YOUR_API_KEY
   ```

   This proxy:
   ```
   http://localhost:8080/shodan/host/1.1.1.1
   ```

3. Get information about the current API plan

   Official API:
   ```
   https://api.shodan.io/api-info?key=YOUR_API_KEY
   ```

   This proxy:
   ```
   http://localhost:8080/api-info
   ```

Note:
1. When using this proxy, you typically don't need to include the API key in each request. The proxy will automatically manage and rotate the configured API keys.
2. If a key parameter is passed in the call, the proxy behaves as follows:
   - If the key value is "shodanproxy" or no key is provided, the proxy will use its own configured API keys.
   - If another key value is provided (e.g., `http://localhost:8080/api-info?key=YOUR_API_KEY`), the proxy will use the key provided by the caller.
   This flexibility allows users to use their own API keys when needed, while also leveraging the proxy's key management functionality.
3. For security reasons, it's recommended to use this proxy server within a controlled environment and not expose it directly to the public internet.

### Shodan CLI Integration

To use the Shodan Proxy with the Shodan CLI, you can modify the API address in the Shodan client library. Execute the following command:

```bash
sed -i 's|https://api.shodan.io|http://your-shodan-proxy-address:port|g' ~/.local/lib/python3.9/site-packages/shodan/client.py
```

Note:
- Replace "your-shodan-proxy-address:port" with the actual address and port where you've deployed the Shodan proxy. For example: `http://localhost:8080` or `http://192.168.1.100:8080`.
- The path `~/.local/lib/python3.9/site-packages/shodan/client.py` may differ depending on your Python installation. Make sure to use the correct path.
- After executing this command, the Shodan CLI will use your Shodan proxy server instead of the official API.

After making this modification, you can use the Shodan CLI as usual, but all requests will be routed through your proxy server.



## GitHub Repository

For more information, contributions, or issues, please visit our [GitHub repository](https://github.com/liuweitao/shodan-proxy).

## License

This project is licensed under the MIT License. For full details, please refer to the [LICENSE](https://github.com/liuweitao/shodan-proxy/blob/main/LICENSE) file in our GitHub repository.

## Support

If you encounter any issues or have questions, please open an issue on our [GitHub repository](https://github.com/liuweitao/shodan-proxy/issues).

## Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](https://github.com/liuweitao/shodan-proxy/blob/main/CONTRIBUTING.md) file for guidelines on how to contribute to this project.

## Security Considerations

- Always use strong, unique passwords for the admin panel.
- Regularly update the Docker image to ensure you have the latest security patches.
- Be cautious when exposing the proxy server to the internet. It's recommended to use it within a controlled network environment.
- Regularly review and update your IP whitelist and path blocking rules.
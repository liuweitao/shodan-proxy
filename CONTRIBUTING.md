# Contributing to Shodan Proxy

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## We Develop with Github
We use github to host code, to track issues and feature requests, as well as accept pull requests.

## Configuration Files
The project includes configuration files (`config/config.yaml` and `config/shodan_keys.yaml`) which are committed to the repository. These files contain default settings and placeholder values.

**Important:** While these files are tracked in the repository, contributors should be cautious about committing changes to them. Only commit changes to these files if you are updating the default settings or structure that should apply to all users.

For your personal use:
1. Make your personal changes directly in these files.
2. Be careful not to accidentally commit these personal changes.

We recommend using environment variables or other methods to override sensitive settings locally without modifying the tracked files.

## We Use [Github Flow](https://guides.github.com/introduction/flow/index.html), So All Code Changes Happen Through Pull Requests
Pull requests are the best way to propose changes to the codebase. We actively welcome your pull requests:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issues](https://github.com/liuweitao/shodan-proxy/issues)
We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/liuweitao/shodan-proxy/issues/new); it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

People *love* thorough bug reports. I'm not even kidding.

## Use a Consistent Coding Style

* We encourage using consistent indentation and formatting throughout the project.
* While we don't strictly enforce a specific style, we recommend following Go's standard formatting conventions where possible.
* Before submitting a pull request, consider running `go fmt` on your code to ensure basic formatting consistency.

## License
By contributing, you agree that your contributions will be licensed under its MIT License.

## References
This document was adapted from the open-source contribution guidelines for [Facebook's Draft.js](https://github.com/facebook/draft-js/blob/main/CONTRIBUTING.md)
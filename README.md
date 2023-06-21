# Sailfish REST API

Sailfish is a RESTful API written in Golang designed to manage notification by 62teknologi.com.

## Table of Contents

1. [Introduction](#introduction)
2. [Catalog Description](#catalog-description)
3. [Features](#features)
4. [Installation](#installation)
5. [API Endpoints](#api-endpoints)
6. [Usage Examples](#usage-examples)
7. [Contributing](#contributing)
8. [License](#license)

## Introduction

Sailfish REST API is built with a focus on simplicity, reliability, and extensibility. With this API, users can manage notification with long pooling, pub/sub, email, and more.

###  Sailfish notification behaviors
- can be created
- can be retrieved
- can be updated
- can be deleted
- can be pushed

## Features

- Easy-to-use RESTful API
- Easy to setup
- Easy to Customizable
- Written in Golang for high performance and concurrency
- Robust data validation and error handling
- Well-documented API endpoints

## Installation

To install and run Sailfish REST API on your local machine, follow these steps:

1. Clone the repository:

   ```git clone https://gitlab.62teknologi.com/62teknologi/sailfish-be-golang.git```


2. Change directory to the cloned repository:

   ```cd sailfish-be-golang```


3. Clone sailfish submodule:

   ```git submodule update --init```


4. Setup config with your own credentials:

   ```cp app.env.example app.env```


4. Build the application:

   ```go build```


5. Run the server:

   ```./sailfish```

The API server will start running at `{{HTTP_SERVER_ADDRESS}}` based on the configuration in app.env. 
You can now interact with the API using your preferred API client or through the command line with `curl`.

## API Endpoints

| Method | Endpoint                   | Description                                                   |
| - |----------------------------|---------------------------------------------------------------|
| GET | /api/v1/notifications    | Retrieve a list of all notifications                          |
| GET | /api/v1/notifications/:id | Retrieve a specific notification by ID                        |
| POST | /api/v1/notifications    | Add a new notification                                        |
| PUT | /api/v1/notifications/:id | Update information for a specific field of notification by ID |
| DELETE | /api/v1/notifications/:id | Delete a specific notification by ID                          |

For more detailed information about each endpoint, including request and response format, please refer to the [API documentation](./API_DOCUMENTATION.md). (WIP)

### Usage Examples  (WIP)

Here are some examples of how to interact with the Sailfish REST API using `curl`:

1. Get a list of all notifications:
2. Get a specific notification by ID:


## Contributing

If you'd like to contribute to the development of the Sailfish REST API, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Commit your changes to the branch
4. Create a pull request, describing the changes you've made

We appreciate your contributions and will review your pull request as soon as possible.

## License

This project is licensed under the MIT License. For more information, please see the [LICENSE](./LICENSE) file.
# GreenkmchSever

GreenkmchSever is a web application that provides a platform for reporting and managing environmental problems. It allows users to send reports of various types of environmental issues, and administrators can view and manage these reports.

## Features

- User authentication and authorization
- Report submission and management
- Administrator dashboard for managing reports and users

## Technologies Used

- Go programming language
- Chi router for HTTP routing
- PostgreSQL for database management
- Migrate for database migrations
- Go-chi for building HTTP handlers

## API Endpoints

### App Routes

#### GET /api/mapp/routes_all

- Description: Retrieves a list of all routes.
- Response: A list of `app.Routes` objects.

#### POST /api/mapp/visit_request

- Description: Submits a visit request.
- Request Body: A `app.RequestData` object.
- Response: A JSON object with a success message.

#### GET /api/mapp/visit_request

- Description: Retrieves the status of a visit request.
- Response: A JSON object with the status of the visit request.

#### POST /api/mapp/send_report

- Description: Submits a report.
- Request Body: A `app.ReportData` object.
- Response: A JSON object with a success message.

#### GET /api/mapp/get_all_reports

- Description: Retrieves all reports.
- Response: A list of `app.Report` objects.

### Portal Routes

#### GET /api/portal/get_all_problems

- Description: Retrieves all problems.
- Response: A list of `portal.Problem` objects.

#### PUT /api/portal/update_problem

- Description: Updates a problem.
- Request Body: A `portal.Problem` object.
- Response: A JSON object with a success message.

### Portal Static Routes

#### GET /portal/troubles

- Description: Retrieves the list of troubles.
- Response: An HTML page with the list of troubles.

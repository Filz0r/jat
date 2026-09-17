# Roadmap

The following roadmap is a rough outline of what is planned before the first stable release. It is subject to change as the project evolves.

- [x] **v0.1.0**
  - Initial project scaffold, API server, database schema, and core job application CRUD endpoints.
- [x] **v0.2.0**
  - Implementation of the React frontend.
  - Multiple bug fixes related to the API requests and responses.
  - Embedding of the React frontend into the final binary release.
  - Improvement of the development workflow.
  - Base implementation of the server setup flow.
  - Refining of some database schemas and introduction of new ones that will be implemented in further releases.
  - Implementation of basic statistics for users and companies (this will be improved in the future).
  - Standardization of response and request bodies.
  - Added a generated OpenAPI/Swagger specification that is imported into the frontend for typesafe requests.
  - Added the ability to archive Application Statuses (before you would delete them only).
  - Enforcing deletion of Job Applications when you soft delete an Application Status or Company.
  - Added the ability to restore changes made to a Company (admin only).
  - Improved the traceability of changes to a Company.
  - Improvements to the Roadmap and Readme structure
- [ ] **v0.3.0**
  - Standardization of API errors.
  - Improving how errors are rendered in the Web UI, some actions trigger toasts, others show directly in the forms.
  - Implementation of an Administrative Dashboard on the Web UI such as:
    - user management
    - restoring deleted user data
    - banning users (time based)
    - enabling/disabling users (permanent)
  - Serving the OpenAPI/Swagger docs while in development.
  - Improving the WebUI branding (favicons and similar)
  - Refactoring the logic used for server configurations, basically the current config module is doing too much due to deprecated code, instead make a new logic to handle server configs and implement a different main file to serve the server instead


## Planed Features
- Adding more admin features to the REST API.
- Adding a build pipeline for both server and web ui.
- Implementing the Client TUI interface.
- Initial docker container builds.
- Implementing standalone TUI interface (there's a possibility this wont be implemented).
- Adding a statistics dashboard for the admin interface (web ui only).
- Adding a notification system (probably with websockets, if possible).
- Adding version checking from within the web server (basically triggers a notification warning all admin users that a new version was released).
- Adding a way to import .csv/.xsl(x) (idk if the latter is possible but that's the idea) files into the web ui.
- Adding time based notifications such as after 2 months of a job application not having changes suggesting changing it to a ghosted state.
- Adding proper CLI option support (as in right now manually parsing arguments is a pain and this can easily be improved).
- Documentation for development.
- Documentation for deployments.
- Open up feature requests.
- Open up external contributions.
- Adding support for S3 storage buckets (to store User CV's, profile pictures, etc).
- Adding support for user account verification via SMTP emails.
- Introduction of a tutorial flow to teach the users how to use the web ui (this will also apply to admin users)
- Adding rate limiting to the API
- Adding the ability to hard delete data (admin only)
- Introduction of automated testing (unit tests and integration tests)

## Planed Improvements
- Improve GitHub dependant stuff like enforcing a standard issue format, tags etc.
- Improving the Application State logic for all users (some of the state_kind should only have 1 record per user).
- Improving the server setup flow (as of writing this you only create an account and set the server as initialized in the final step)

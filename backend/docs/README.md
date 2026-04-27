# Generating swagger.yaml

## Prerequisites

- swag CLI is installed (if not, it will be installed automatically by the script)
- API version is updated in `backend/server/routes.go` file (e.g., `@version 1.0`)

## Simple Command
Run the following command from `backend` folder:

```
swag init --ot yaml -g server/routes.go
```

### Run Script For Generating and Formating Swagger File for Better Comparability

Run the following command from the repository root:

```bash
python3 backend/docs/generate-swagger-spec.py
```

## Validation
Check console outputs for warnings or errors. If the script runs successfully, it will regenerate `backend/docs/swagger.yaml` with the latest API documentation based on the Swaggo annotations in `backend/server/routes.go`.
After generating the `swagger.yaml`, you can validate it using command-line tools like `swagger-cli`:

```
docker run -p 80:8080 -e SWAGGER_JSON=/docs/swagger.yaml -v $(pwd)/docs/swagger.yaml:/docs/swagger.yaml docker.swagger.io/swaggerapi/swagger-ui
```

Then open `http://localhost` in your browser to view the Swagger UI and check for any validation errors.

## Version Control
Finally, you have to commit the updated `swagger.yaml` file to the repository to keep the API documentation up to date.
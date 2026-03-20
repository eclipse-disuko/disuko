# Token Management

- **Scope:** You can create API tokens for your projects to enable your software supplier to use the Disclosure Portal API for a specific project. When the token is created on a project group level, then it is valid for all children projects below that project group. If multiple suppliers provide data to a single project, make sure to create different tokens for each of them to ensure traceability of SBOM deliveries.
- **Disclaimer:** The project token has to be treated confidentially and may only be transmitted via encrypted communication channels. For security reasons, do not transmit the project UUID in the same message as the project token.
- **Validity:** The token has a defined lifetime, after which it automatically expires. You can also manually expire an active token. In consequence, an API Token is valid until the expiry date is reached, or until it is explicitly revoked or updated.
- **Transmission:** To transmit an API token to your supplier in a safe way, you can send it as encrypted email in the message body.
- **Security:** The project token has to be treated confidentially and may only be transmitted via encrypted communication channels. For security reasons, do not transmit the project UUID in the same message as the project token.
- **Endpoint:** The endpoint to be used depends on the specific instance of the FOSS Disclosure Portal, on which the token is generated. The endpoint is displayed after creating the token. On this dialog you have an action to copy all information to your clipboard.
- **Automation:** For ease of access to the FOSS Disclosure Portal API a [Command Line Interface Client (CLI)](https://github.com/eclipse-disuko/disuko-cli) is provided on the Open Source Organization of eclipse-disuko on public GitHub.com.

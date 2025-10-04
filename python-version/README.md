# Gmail AI Telegram Agent (Python)

This is a Python translation of the Go Gmail AI Telegram agent. It fetches unread emails from your Gmail account, uses OpenAI to classify them, and then notifies you on Telegram about important emails while automatically trashing the unimportant ones.

## How it works

The agent performs the following steps:

1. Connects to your Gmail account using the Gmail API.
2. Fetches all unread emails from your inbox.
3. For each email, it uses OpenAI's language model to determine if the email is important.
4. If an email is classified as important, a notification is sent to your specified Telegram chat. The email is then marked as read in Gmail.
5. If an email is classified as unimportant, it is moved to the trash in Gmail.

The agent polls your Gmail account for new unread messages at a configurable interval.

## Configuration

To use this agent, you need to configure the following:

1. **Gmail API Credentials**:
   * Go to the [Google Cloud Console](https://console.cloud.google.com/).
   * Create a new project.
   * Enable the "Gmail API".
   * Create an "OAuth 2.0 Client ID" for a "Desktop application".
   * Download the credentials as `gmail-credentials.json` and place it in the `configs` directory of this project.

2. **Configuration Files**:
   * Copy `config.example.yaml` to `config.yaml` and `prompts.example.yaml` to `prompts.yaml`.
   * Edit `config.yaml` and fill in your actual credentials:
       * `openai_api_key`: Your API key from OpenAI (keep empty for local LLM).
       * `telegram_bot_token`: The token for your Telegram bot (you can get this from the "BotFather" on Telegram).
       * `telegram_chat_id`: The ID of the Telegram chat where you want to receive notifications.
   * The `prompts.yaml` file contains the templates for the prompts used by the AI. You can customize these prompts to better suit your needs.
   * You can use any local LLM that is compatible with the OpenAI API by changing the `openai.endpoint` in `config.yaml`.

## How to Run

1. **Install Python**: Make sure you have Python 3.8+ installed on your system.
2. **Install Dependencies**: Open a terminal in the project's root directory and run:
   ```bash
   pip install -r requirements.txt
   ```
3. **First Run (OAuth2 Authentication)**:
   * Run the application for the first time:
       ```bash
       python src/main.py
       ```
   * This will open a new page in your web browser asking you to authorize the application to access your Gmail account.
   * After you grant permission, the application will save a `token.json` file in the project's root directory. This token will be used for future authentications.
4. **Subsequent Runs**:
   * After the initial authentication, you can run the agent with the same command:
       ```bash
       python src/main.py
       ```
   * The agent will start polling your Gmail account for unread messages.

## Running with Docker

You can also run the agent using Docker.

1. **Build the Docker image**:
   ```bash
   docker build -t gmail-ai-agent-python .
   ```

2. **Prepare the configuration directory**:
   * Make sure you have the `gmail-credentials.json`, `config.yaml`, and `prompts.yaml` files in a single directory on your host machine. For this example, we'll assume they are in a directory named `configs` in your project root.

3. **First Run (OAuth2 Authentication)**:
   * Run the Docker container interactively (`-it`) with the configuration directory mounted as a volume. This is necessary to complete the one-time browser-based authentication.
   ```bash
   docker run -it -v $(pwd)/configs:/app/configs gmail-ai-agent-python
   ```
   * Follow the instructions in the console to authorize the application. A `token.json` file will be created in your local `configs` directory, allowing future runs to be non-interactive.

4. **Subsequent Runs**:
   * For all subsequent runs, you can run the container in detached mode (`-d`) since the authentication token is already present:
   ```bash
   docker run -d -v $(pwd)/configs:/app/configs gmail-ai-agent-python
   ```
   The agent will now run in the background and poll your Gmail account.

## Disclaimer

This is a prototype and should be used with care. For production use, consider the following:

* **Error Handling**: The current error handling is basic.
* **Security**: Credentials and tokens should be stored securely.
* **Deployment**: For continuous operation, you should run this as a service (e.g., using Docker or systemd).

import os
import rollbar
from dotenv import load_dotenv


load_dotenv(override=True)


class Logger:
    # Initialize Rollbar
    rollbar.init(
        access_token=os.getenv("ROLLBAR_ACCESS_TOKEN", ""),
        environment=os.getenv("ROLLBAR_ENVIRONMENT", "development"),
        code_version="1.0"
    )

    def payload_handler(payload, **kw):
        payload['data']['custom'] = {
            'service': 'propuestas-api',
            'version': '0.1.0'
        }
        return payload

    rollbar.events.add_payload_handler(payload_handler)

    def info(self, message):
        rollbar.report_message(message, "info")

    def error(self, message):
        rollbar.report_message(message, "error")

    def warning(self, message):
        rollbar.report_message(message, "warning")

    def critical(self, message):
        rollbar.report_message(message, "critical")

    def report_exc_info(self, extra_data=None):
        if extra_data is None:
            extra_data = {}
        rollbar.report_exc_info(extra_data=extra_data)



def DeployLogger() -> dict:

    import requests

    url = "https://api.rollbar.com/api/1/deploy"
    

    payload: dict[str, str] = {
        "environment": os.getenv("ROLLBAR_ENVIRONMENT", "development"),
        "rollbar_username": "osmargm1202",
        "comment": "gcr deployment",
        "revision": os.getenv("APP_REVISION", "")
    }
    headers: dict[str, str] = {
        "accept": "application/json",
        "X-Rollbar-Access-Token": os.getenv("ROLLBAR_ACCESS_TOKEN", ""),
        "content-type": "application/json"
    }

    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        try:
            body = response.json()
        except Exception:
            body = response.text
        return {"status_code": response.status_code, "body": body}
    except Exception as e:
        return {"status_code": 0, "body": str(e)}


logger: Logger = Logger()
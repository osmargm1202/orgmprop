
from google.oauth2 import service_account
from google.auth.transport.requests import Request
from typing import LiteralString
import json
import os
from typing import Any

def gcloud_token(url_api: LiteralString):

    info: Any = json.loads(os.getenv("GOOGLE_CREDENTIALS"))


    id_creds: service_account.IDTokenCredentials = service_account.IDTokenCredentials.from_service_account_info(
        info,
        target_audience=url_api
    )

    # Refrescar (obtener token válido)
    id_creds.refresh(Request())
    return id_creds.token

import os
import redis
from _logs_config import Logger


logger: Logger = Logger()

class RedisClient:
    """
    Cliente para manejar la conexión a Redis.
    """

    url: str = os.getenv("REDIS_URL", "")
    if not url:
        raise ValueError("REDIS_URL no está configurada en settings.")

    def __init__(self):
        try:
            self._client = redis.Redis.from_url(self.url)
            # Test connection
            self._client.ping()
        except Exception as e:
            logger.report_exc_info(extra_data={
                'error_type': 'redis_connection_error',
                'redis_url': self.url,
                'message': 'Error al conectar con Redis'
            })
            raise ValueError(f"Error al conectar con Redis: {e}")

    def get(self, key):
        try:
            return self._client.get(key)
        except Exception as e:
            logger.report_exc_info(extra_data={
                'error_type': 'redis_get_error',
                'key': key,
                'redis_url': self.url,
                'message': 'Error al obtener valor de Redis'
            })
            raise ValueError(f"Error al obtener valor de Redis: {e}")

    def set(self, key, value, ex=None):
        try:
            self._client.set(key, value, ex=ex)
        except Exception as e:
            logger.report_exc_info(extra_data={
                'error_type': 'redis_set_error',
                'key': key,
                'value': str(value),
                'expiration': ex,
                'redis_url': self.url,
                'message': 'Error al guardar valor en Redis'
            })
            raise ValueError(f"Error al guardar valor en Redis: {e}")
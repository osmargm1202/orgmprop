from pathlib import Path
from _logs_config import Logger
from main import s3_client, R2_BUCKET, R2_ENDPOINT_URL


logger: Logger = Logger()


def upload_to_r2(file_path: Path, key: str) -> str:
    """Upload file to Cloudflare R2 and return public URL"""
    try:
        with open(file_path, 'rb') as file:
            s3_client.upload_fileobj(file, R2_BUCKET, key)
        
        # Return the R2 URL
        public_url = f"{R2_ENDPOINT_URL}/{R2_BUCKET}/{key}"
        return public_url
    except Exception as e:
        logger.report_exc_info(extra_data={
            "error_type": "r2_upload_error",
            "file_path": str(file_path),
            "key": key,
            "error": str(e)
        })
        raise


def update_static() -> None:
    """Update static files in R2"""
    try:
        for file in Path("propuesta").iterdir():
            if file.is_file():
                upload_to_r2(file, f"templates/{file.name}")
                print(f"{file.name} uploaded")
            elif file.is_dir():
                for subfile in file.iterdir():
                    upload_to_r2(subfile, f"templates/{file.name}/{subfile.name}")
                    print(f"{file.name}/{subfile.name} uploaded")
            
    except Exception as e:
        logger.report_exc_info(extra_data={
            "error_type": "r2_upload_error",
            "error": str(e)
        })
        raise

if __name__ == "__main__":
    update_static()

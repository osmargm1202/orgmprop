
FROM python:3.13-slim

# Create necessary directories
RUN mkdir -p /app


RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libcairo2 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libgdk-pixbuf-2.0-0 \
    libffi-dev \
    libglib2.0-0 \
    fonts-dejavu-core \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the application
COPY . /app

# Install uv
RUN pip install uv

# Install dependencies using uv
RUN uv sync

# Expose the port
EXPOSE 8000

# Run the application
CMD ["sh", "-c", "uv run uvicorn main:app --host 0.0.0.0 --port 8000"]
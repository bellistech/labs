# Python Project Quick Reference & Docker Setup

## Quick Reference Card

### Virtual Environment

```bash
# Create
python3 -m venv venv

# Activate
source venv/bin/activate      # Linux/macOS
venv\Scripts\activate         # Windows

# Deactivate
deactivate

# Install from requirements
pip install -r requirements.txt

# Freeze current packages
pip freeze > requirements.txt
```

### Common Imports

```python
# Standard Library
import os
import sys
import json
import logging
from pathlib import Path
from typing import List, Dict, Optional, Any, Tuple
from datetime import datetime
from dataclasses import dataclass, field
from enum import Enum
from contextlib import contextmanager
import asyncio

# Third Party
import requests                    # HTTP client
from pydantic import BaseModel     # Data validation
from fastapi import FastAPI        # Web framework
import psycopg2                    # PostgreSQL
from sqlalchemy import create_engine  # ORM
import pytest                      # Testing
```

### Type Hints Cheat Sheet

```python
from typing import List, Dict, Optional, Union, Tuple, Any, Callable

# Basic types
name: str = "value"
count: int = 0
price: float = 9.99
active: bool = True

# Collections
names: List[str] = ["a", "b"]
scores: Dict[str, int] = {"alice": 100}
point: Tuple[int, int] = (10, 20)

# Optional (can be None)
result: Optional[str] = None

# Union (multiple types)
value: Union[int, str] = 42

# Function signatures
def process(items: List[str]) -> Dict[str, int]:
    pass

# Callable
handler: Callable[[str, int], bool]
```

### Logging Setup

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('app.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)
logger.info("Application started")
logger.error("Error occurred", exc_info=True)
```

### Environment Variables

```python
import os
from dotenv import load_dotenv

load_dotenv()  # Load from .env file

DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://localhost/db")
DEBUG = os.getenv("DEBUG", "false").lower() == "true"
PORT = int(os.getenv("PORT", 8000))
```

### FastAPI Quick Start

```python
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel

app = FastAPI(title="My API")

class Item(BaseModel):
    name: str
    price: float

@app.get("/items/{item_id}")
async def get_item(item_id: int):
    return {"item_id": item_id}

@app.post("/items", status_code=201)
async def create_item(item: Item):
    return item

# Run: uvicorn main:app --reload
```

### PostgreSQL with psycopg2

```python
import psycopg2
from contextlib import contextmanager

@contextmanager
def get_connection():
    conn = psycopg2.connect(
        host="localhost",
        database="mydb",
        user="user",
        password="pass"
    )
    try:
        yield conn
    finally:
        conn.close()

with get_connection() as conn:
    with conn.cursor() as cur:
        cur.execute("SELECT * FROM users WHERE id = %s", (user_id,))
        row = cur.fetchone()
```

### SQLAlchemy ORM

```python
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

Base = declarative_base()
engine = create_engine("postgresql://user:pass@localhost/db")
Session = sessionmaker(bind=engine)

class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True)
    name = Column(String(100), nullable=False)

# Create tables
Base.metadata.create_all(engine)

# Query
session = Session()
users = session.query(User).filter(User.name.like("%john%")).all()
```

### Pytest Patterns

```python
import pytest

# Basic test
def test_addition():
    assert 1 + 1 == 2

# Fixture
@pytest.fixture
def sample_data():
    return {"key": "value"}

def test_with_fixture(sample_data):
    assert sample_data["key"] == "value"

# Parametrize
@pytest.mark.parametrize("input,expected", [
    (1, 2),
    (2, 4),
    (3, 6),
])
def test_double(input, expected):
    assert input * 2 == expected

# Exception testing
def test_raises():
    with pytest.raises(ValueError):
        raise ValueError("error")

# Async test
@pytest.mark.asyncio
async def test_async():
    result = await async_function()
    assert result == "expected"
```

---

## Docker Setup for Capstone Project

### Dockerfile

```dockerfile
# Dockerfile
FROM python:3.11-slim

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Set work directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libpq-dev \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY app/ ./app/
COPY alembic/ ./alembic/
COPY alembic.ini .

# Create non-root user
RUN useradd --create-home --shell /bin/bash appuser && \
    chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# Run application
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: network-api
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://netadmin:secretpassword@db:5432/network_inventory
      - LOG_LEVEL=INFO
      - SECRET_KEY=${SECRET_KEY:-changeme}
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  db:
    image: postgres:15-alpine
    container_name: network-db
    environment:
      - POSTGRES_USER=netadmin
      - POSTGRES_PASSWORD=secretpassword
      - POSTGRES_DB=network_inventory
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "5432:5432"
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U netadmin -d network_inventory"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Optional: pgAdmin for database management
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: network-pgadmin
    environment:
      - PGADMIN_DEFAULT_EMAIL=admin@example.com
      - PGADMIN_DEFAULT_PASSWORD=admin
    ports:
      - "5050:80"
    depends_on:
      - db
    networks:
      - app-network
    profiles:
      - tools

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### docker-compose.test.yml

```yaml
# docker-compose.test.yml
version: '3.8'

services:
  test-db:
    image: postgres:15-alpine
    container_name: network-test-db
    environment:
      - POSTGRES_USER=test
      - POSTGRES_PASSWORD=test
      - POSTGRES_DB=test_network_inventory
    ports:
      - "5433:5432"
    tmpfs:
      - /var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U test -d test_network_inventory"]
      interval: 5s
      timeout: 3s
      retries: 5

  tests:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: network-api-tests
    environment:
      - DATABASE_URL=postgresql://test:test@test-db:5432/test_network_inventory
      - TESTING=true
    depends_on:
      test-db:
        condition: service_healthy
    command: >
      sh -c "
        alembic upgrade head &&
        pytest tests/ -v --cov=app --cov-report=xml --cov-report=term
      "
    volumes:
      - ./tests:/app/tests:ro
      - ./coverage:/app/coverage
```

### init.sql (Database Initialization)

```sql
-- init.sql
-- Initial database setup

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types
DO $$ BEGIN
    CREATE TYPE device_status AS ENUM ('up', 'down', 'maintenance', 'unknown');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE device_type AS ENUM ('router', 'switch', 'firewall', 'server', 'other');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE network_inventory TO netadmin;
```

### requirements.txt

```
# Core
fastapi>=0.100.0
uvicorn[standard]>=0.23.0
pydantic>=2.0.0
pydantic-settings>=2.0.0

# Database
psycopg2-binary>=2.9.0
sqlalchemy>=2.0.0
alembic>=1.11.0

# Utilities
python-dotenv>=1.0.0
python-multipart>=0.0.6
httpx>=0.24.0

# Security
python-jose[cryptography]>=3.3.0
passlib[bcrypt]>=1.7.4
```

### requirements-dev.txt

```
# Testing
pytest>=7.4.0
pytest-cov>=4.1.0
pytest-asyncio>=0.21.0
httpx>=0.24.0

# Code quality
black>=23.0.0
isort>=5.12.0
flake8>=6.0.0
mypy>=1.4.0
pre-commit>=3.3.0

# Documentation
mkdocs>=1.5.0
mkdocs-material>=9.0.0
```

### .pre-commit-config.yaml

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files

  - repo: https://github.com/psf/black
    rev: 23.7.0
    hooks:
      - id: black
        language_version: python3.11

  - repo: https://github.com/pycqa/isort
    rev: 5.12.0
    hooks:
      - id: isort
        args: ["--profile", "black"]

  - repo: https://github.com/pycqa/flake8
    rev: 6.0.0
    hooks:
      - id: flake8
        args: ["--max-line-length=100", "--ignore=E501,W503"]

  - repo: https://github.com/pre-commit/mirrors-mypy
    rev: v1.4.1
    hooks:
      - id: mypy
        additional_dependencies: [types-requests]
        args: ["--ignore-missing-imports"]
```

### pyproject.toml

```toml
# pyproject.toml
[build-system]
requires = ["setuptools>=61.0", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "network-api"
version = "1.0.0"
description = "Network Inventory REST API"
readme = "README.md"
requires-python = ">=3.10"
license = {text = "MIT"}
authors = [
    {name = "Your Name", email = "your.email@example.com"}
]
classifiers = [
    "Development Status :: 4 - Beta",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
]

[tool.black]
line-length = 100
target-version = ['py310', 'py311']
include = '\.pyi?$'
exclude = '''
/(
    \.git
    | \.venv
    | venv
    | build
    | dist
)/
'''

[tool.isort]
profile = "black"
line_length = 100
skip = [".venv", "venv"]

[tool.mypy]
python_version = "3.11"
warn_return_any = true
warn_unused_configs = true
ignore_missing_imports = true

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
python_functions = ["test_*"]
asyncio_mode = "auto"
addopts = "-v --tb=short"

[tool.coverage.run]
source = ["app"]
omit = ["tests/*", "venv/*"]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "raise NotImplementedError",
]
```

---

## GitHub Actions Alternative to Jenkins

### .github/workflows/ci.yml

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  PYTHON_VERSION: '3.11'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: ${{ env.PYTHON_VERSION }}
          cache: 'pip'
      
      - name: Install dependencies
        run: |
          pip install --upgrade pip
          pip install black isort flake8 mypy
      
      - name: Run Black
        run: black --check app/ tests/
      
      - name: Run isort
        run: isort --check-only app/ tests/
      
      - name: Run Flake8
        run: flake8 app/ tests/ --max-line-length=100
      
      - name: Run MyPy
        run: mypy app/ --ignore-missing-imports

  test:
    runs-on: ubuntu-latest
    needs: lint
    
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test_db
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: ${{ env.PYTHON_VERSION }}
          cache: 'pip'
      
      - name: Install dependencies
        run: |
          pip install --upgrade pip
          pip install -r requirements.txt
          pip install -r requirements-dev.txt
      
      - name: Run tests
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/test_db
        run: |
          pytest tests/ -v --cov=app --cov-report=xml --cov-report=term
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.xml

  build:
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/main'
    
    permissions:
      contents: read
      packages: write
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}
      
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    environment: production
    
    steps:
      - name: Deploy to server
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.DEPLOY_HOST }}
          username: ${{ secrets.DEPLOY_USER }}
          key: ${{ secrets.DEPLOY_KEY }}
          script: |
            cd /opt/network-api
            docker compose pull
            docker compose up -d
            docker system prune -f
```

---

## Usage Commands

```bash
# Development
docker compose up -d                      # Start all services
docker compose logs -f api                # Follow API logs
docker compose exec api bash              # Shell into API container
docker compose down                       # Stop all services

# Testing
docker compose -f docker-compose.test.yml up --abort-on-container-exit

# Production build
docker build -t network-api:latest .
docker push your-registry.com/network-api:latest

# Database operations
docker compose exec db psql -U netadmin -d network_inventory
docker compose exec api alembic upgrade head
docker compose exec api alembic revision --autogenerate -m "Add new table"

# Pre-commit hooks
pre-commit install
pre-commit run --all-files
```

---

This completes the Docker setup and quick reference for your Python learning path capstone project!

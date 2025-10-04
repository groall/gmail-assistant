FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Copy requirements first for better caching
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY src/ ./src/
COPY configs/ ./configs/

# Make main.py executable
RUN chmod +x src/main.py

# Set Python path
ENV PYTHONPATH=/app

# Run the application
CMD ["python", "src/main.py"]

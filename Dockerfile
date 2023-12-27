FROM python:3.10

WORKDIR /app

# -- Install poetry
RUN pip install poetry

# -- Install requirements
COPY pyproject.toml .
COPY poetry.lock .
ENV POETRY_VIRTUALENVS_CREATE=false
RUN poetry install --no-dev

# -- Copy the app
COPY . .

# -- Run the app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3000"]
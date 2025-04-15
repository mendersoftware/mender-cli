FROM python:3.13-slim

RUN pip3 install requests pytest cryptography boto3 docker

RUN mkdir -p /tests
ENTRYPOINT ["python3", "-m", "pytest", "tests"]

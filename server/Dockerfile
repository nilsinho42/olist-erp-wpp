FROM python:3.9-alpine

WORKDIR /server

COPY auth_server.py /server

RUN pip install flask

EXPOSE 5100

CMD ["python", "auth_server.py"]

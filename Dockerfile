FROM python:3.10

RUN mkdir /app

COPY ./requirements.txt /app/requirements.txt

RUN pip install -r /app/requirements.txt

COPY ./src /app

WORKDIR /app

RUN chmod a+x /app/scripts/*.sh


FROM python:3-alpine

# This config forces that the python output is unbuffered that way we can see
# the application logs right in the terminal
# Reference -> https://docs.python.org/3/using/cmdline.html#envvar-PYTHONUNBUFFERED
ENV PYTHONUNBUFFERED=1

RUN mkdir -p /conf
RUN mkdir -p /app
# WORKDIR /app

RUN apk add build-base

COPY requirements.txt /conf
RUN pip3 install -r /conf/requirements.txt

# COPY . /app

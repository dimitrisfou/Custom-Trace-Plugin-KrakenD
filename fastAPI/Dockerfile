FROM python:3.10.2-buster

WORKDIR /code/app

COPY requirements.txt /code/app

RUN pip3 install -r requirements.txt

COPY app.py /code/app

CMD ["python3" , "-u", "app.py"]
# , "-u",
FROM python:3.6-onbuild

RUN easy_install pip

ADD requirements.txt /src/requirements.txt
RUN cd /src; pip install -r requirements.txt

ADD . /src

EXPOSE 5000

CMD ["python", "/src/application.py"]

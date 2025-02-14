@echo off
cd ..
docker build -f Dockerfile -t johnkoal/artdotech.website:latest .
docker push johnkoal/artdotech.website:latest
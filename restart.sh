docker-compose down
docker volume rm sdv_projet_aws_bdd-data
docker-compose build --no-cache
docker-compose up -d
watch docker ps
version: '3'

services:
  foobar:
    build:
      dockerfile: Dockerfile
      context: .
    env_file:
      - ./env/local.env
    volumes:
      - './env/local.config.yaml:/config.yaml'
    ports:
    - '8090:8090'
    networks:
      scaffold:

networks:
  scaffold: {}



# Lighthouse

A simple and efficient URL shortener that is designed to be built into the Panda-Education ecosystem.

## Overview

This application is designed as a microservice, and should sit behind a series of API gateways (e.g. 
authentication, API, rate limiting gateways). As such, this application is only focused on speed. As 
of now, the application is deployed using docker, hosted with its own dedicated PostGreSQL instance 
running in a sister docker container. Future developments would give way for a centralised RDS solution
for the entire Panda-Education ecosystem; only then will this repo contain code purely for URL shortening.

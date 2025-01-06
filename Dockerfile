FROM debian:stable-slim
COPY budget_buddy /bin/budget_buddy
COPY .env /.env
ENV PORT=8080
CMD [ "/bin/budget_buddy" ] 
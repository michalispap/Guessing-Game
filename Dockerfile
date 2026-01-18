FROM debian:stable-slim
WORKDIR /app
COPY guessing-game .
RUN chmod +x guessing-game
CMD ["./guessing-game"]
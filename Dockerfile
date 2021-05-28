FROM golang:alpine AS builder

WORKDIR /service

COPY app /service

RUN go build .

RUN ls

FROM scratch

WORKDIR /service

COPY --from=builder /service/feedback_bot /service/feedback_bot

ENV FEEDBACK_BOT_DB_FILE_PATH=messages.db

EXPOSE 6654

CMD "./feedback_bot"
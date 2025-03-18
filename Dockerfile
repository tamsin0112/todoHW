# 使用官方的 Golang 映像作為基礎映像
FROM golang:1.23

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 並下載依賴
COPY hw1.go hw1.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN echo "before bulid" && sleep 2
RUN ls -al && sleep 2

# 編譯應用程式
RUN go build -v

RUN echo "after bulid" && sleep 2

RUN ls -al && sleep 2

ENTRYPOINT [ "./hw1" ]
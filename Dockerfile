FROM golang:1.24.1

LABEL authors="KORIEBRUH AKA JAMAL"
# Set working directory
WORKDIR /app

# Copy seluruh isi proyek ke dalam container
COPY . .

# Install dependencies
RUN go mod download

# Compile aplikasi
RUN go build -o main .

# Jalankan aplikasi
CMD ["./main"]
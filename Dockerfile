FROM golang:1.26

WORKDIR /workspace

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Cache go module downloads in their own layer
COPY api/go.mod api/go.sum ./api/
RUN cd api && go mod download

# Source code is bind-mounted by docker-compose for hot reload during
# development; it's copied here too so the image stays runnable on its own.
COPY . .

EXPOSE 8080

CMD ["air"]
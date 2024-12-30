# **YC-W22 Dating App**

Welcome to the **YC-W22 Dating App**, a modern and scalable dating application designed for seamless matchmaking. Built with **Go** (Golang) as the backend, this app is optimized for high performance, security, and ease of use. 🚀

---

## **Features**

- **Authentication**: Secure sign-up and login with hashed passwords.
- **Swiping**: Swipe left (pass) or right (like) on profiles.
- **Daily Swipe Limits**: Enforced limits for non-premium users.
- **Premium Features**:
  - Verified user labels.
  - Unlimited swipes. (Future update)
- **Mutual Matching**: Match users who like each other.
- **Secure Communication**:
  - JWT-based authentication.
  - Encrypted data transfer (HTTPS recommended).

---

## **Project Structure**

```plaintext
yc-w22-dating-app-valdy/
├── cmd/                      # Main entry point
│   └── server/
│       └── main.go           # Starts the server
├── config/                   # Configuration files
│   ├── config.go             # App configuration logic
│   ├── config.yml            # YAML-based settings
├── internal/                 # Core application logic
│   ├── api/                  # Handlers for routes
│   ├── models/               # Database models
│   ├── repository/           # Data access layer
│   ├── services/             # Business logic
├── pkg/                      # Shared utilities and packages
│   ├── db/                   # Database connection
│   ├── logger/               # Logging setup
│   ├── jwt/                  # JWT utilities
├── migrations/               # Database migration scripts
├── go.mod                    # Dependency management
└── README.md                 # Project documentation
```

## **Tech Stack**

- **Backend**: Go (Golang)
- **Database**: PostgreSQL, Redis
- **Authentication**: JWT (JSON Web Tokens)
- **ORM**: GORM
- **API Framework**: Echo
- **Encryption**: AES-GCM, bcrypt

---

## **Getting Started**

### Prerequisites

1. Install [Go](https://go.dev/dl/).
2. Install [PostgreSQL](https://www.postgresql.org/).
3. Install [Redis Server](https://redis.io/)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/vivaldy22/yc-w22-dating-app-valdy.git
   cd yc-w22-dating-app-valdy
   ```
2. Copy `config/config.yml.example` file into `config/config.yml`:
3. Run `go mod download`
4. Start the server:
   ```bash
   make dating-app
   ```
   or
   ```bash
   go run cmd/server/main.go
   ```  

## **API Endpoints**

### **Authentication**
| Method | Endpoint       | Description                |
|--------|----------------|----------------------------|
| POST   | `/v1/auth/signup`      | Create a new user account  |
| POST   | `/v1/auth/login`       | Login and retrieve JWT     |

### **Onboard**
| Method | Endpoint              | Description                        |
|--------|------------------------|------------------------------------|
| GET    | `/v1/onboard/swipe/profiles` | Get profiles available for swiping |
| POST   | `/v1/onboard/swipe/swipe/pass`              | Record a swipe pass        |
| POST   | `/v1/onboard/swipe/swipe/like`              | Record a swipe like        |
| POST   | `/v1/onboard/premium/buy`    | Buy a premium feature    |

---

## **Testing**

Run unit tests using:
```bash
go test ./...
```
or
```bash
make test
```

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**

- Built with ❤️ by [Vivaldy22](https://github.com/vivaldy22).

---

## **Contact**

Feel free to reach out if you have questions or suggestions:

- **Author**: [Vivaldy22](https://github.com/vivaldy22)


# Backer

**Backer** is a crowdfunding platform backend that connects campaign creators with donors. It handles everything needed to run a donation-based platform — from user authentication and campaign management to payment processing and transaction tracking.

## About the Project

Backer allows users to:
- Create an account and log in securely
- Launch and manage fundraising campaigns
- Receive donations from backers/donors
- Process payments through an integrated payment gateway
- Track and review transaction history

Built entirely in **Go**, the project follows a clean, modular architecture separating each domain (auth, campaign, payment, transaction, user) into its own package — making it easy to maintain and extend.

## Tech Stack

- **Language:** Go (Golang)
- **Architecture:** Modular, domain-driven structure

## Project Structure

```
backer/
├── auth/          # Authentication & authorization logic
├── campaign/       # Campaign domain (model, service, repository)
├── config/         # App configuration (database, env, etc.)
├── handler/         # HTTP handlers / controllers
├── helper/           # Utility functions & response formatting
├── payment/          # Payment gateway integration
├── transaction/       # Transaction domain
├── user/               # User domain
├── go.mod
├── go.sum
└── main.go
```

## Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/agusdedi/backer.git
   cd backer
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment configuration**

   Configure your database and environment variables in the `config/` folder.

4. **Run the application**
   ```bash
   go run main.go
   ```

## API Documentation

Full API documentation, including all available endpoints, request/response examples, and authentication details, is published via Postman:

[![Postman Documentation](https://img.shields.io/badge/API%20Docs-Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)](https://documenter.getpostman.com/view/30805799/2sBY4HTiWm)

Or import the collection directly into Postman:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://documenter.getpostman.com/view/30805799/2sBY4HTiWm)

## Contributing

Contributions and suggestions are welcome. For major changes, please open an issue first to discuss what you'd like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with by [agusdedi](https://github.com/agusdedi)

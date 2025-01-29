# Pump.fun Token Scanner

This is a **Solana blockchain scanner** that listens for **transactions involving the Pump.fun token** and decodes relevant trade data. It utilizes the **Solana RPC API** to subscribe to transaction logs, fetch transaction details, and extract buy orders.

## Features

- **Real-time monitoring** of transactions mentioning the Pump.fun program.
- **Transaction decoding** to extract trade actions (Buy, Sell, Create).
- **Balance tracking** before and after transactions.
- **Concurrency support** for efficient transaction handling.

## Technologies Used

- **Golang** – Primary programming language.
- **Solana-Go SDK** – Interact with the Solana blockchain.
- **RPC & WebSockets** – Fetch transaction data in real-time.
- **Rate Limiting** – Controls API request rate to avoid exceeding limits.

## Installation

### Prerequisites
Ensure you have **Go 1.23+** installed.

### Future Improvements
	•	Add database storage for tracking historical trades.
	•	Implement error handling & retry mechanisms.
	•	Optimize multi-threading for better performance.

### License

This project is open-source under the MIT License.
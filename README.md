# WebSocket Bidding System API

## Introduction

This README provides a guide to interacting with the WebSocket server for a bidding system. The server allows users to create rooms for auctions, start auctions, and place bids. This document outlines the API structure, message formats, examples, and error handling to help developers integrate with the server.

## API Documentation

### Endpoints

- **`/room/{id}`**:  
  - **Purpose**: Connect to a specific room for bidding.
  - **Behavior**:
    - If the room does not exist, the server expects an `AuctionStartReq` message to create and start a new auction.
    - If the room exists, the client joins the auction as a bidder.

### Message Formats

#### 1. AuctionStartReq

- **Purpose**: Start a new auction in a room.
- **Format**:
  ```json
  {
    "Item": "string",
    "Price": uint,
    "Seller": "string"
  }
  ```
- **Example**:
  ```json
  {
    "Item": "Laptop",
    "Price": 1000,
    "Seller": "Seller123"
  }
  ```

#### 2. Bid

- **Purpose**: Place a bid in an ongoing auction.
- **Format**:
  ```json
  {
    "Item": "string",
    "Bid": uint,
    "Bidder": "string"
  }
  ```
- **Example**:
  ```json
  {
    "Item": "Laptop",
    "Bid": 1200,
    "Bidder": "Bidder456"
  }
  ```

## Examples

### Creating a New Room and Starting an Auction

1. **Connect to the WebSocket server**:
   ```javascript
   const socket = new WebSocket(`ws://localhost:8080/room/1`);
   ```
2. **Send an AuctionStartReq message**:
   ```javascript
   socket.send(JSON.stringify({
     "Item": "Laptop",
     "Price": 1000,
     "Seller": "Seller123"
   }));
   ```
3. **Receive confirmation**:
   ```javascript
   socket.onmessage = function(event) {
     console.log(event.data); // "your message was received"
   };
   ```

### Joining an Existing Room and Placing a Bid

1. **Connect to the WebSocket server**:
   ```javascript
   const socket = new WebSocket(`ws://localhost:8080/room/1`);
   ```
2. **Send a Bid message**:
   ```javascript
   socket.send(JSON.stringify({
     "Item": "Laptop",
     "Bid": 1200,
     "Bidder": "Bidder456"
   }));
   ```
3. **Receive confirmation**:
   ```javascript
   socket.onmessage = function(event) {
     console.log(event.data); // "your message was received"
   };
   ```

## Error Handling

- **Invalid Room ID**:
  - **Error Message**: "invalid id"
  - **Cause**: The room ID provided is not a valid integer.
  - **Solution**: Ensure the room ID is a valid integer.

- **Room Does Not Exist**:
  - **Behavior**: The server will prompt for an `AuctionStartReq` to create a new room.
  - **Solution**: Provide the necessary `AuctionStartReq` message to start a new auction.

- **JSON Unmarshal Errors**:
  - **Error Message**: Error details will be logged on the server side.
  - **Cause**: The message sent is not in the correct JSON format.
  - **Solution**: Ensure that the messages sent adhere to the specified JSON formats.

## Limitations and Future Work

- **Timer Functionality**: Planned but not currently implemented. Future updates will include timer broadcasting to all bidders.
- **SellerManager**: Currently commented out in the code; planned for future implementation.

## Conclusion

This WebSocket API allows clients to interact with the bidding system by creating and joining rooms, starting auctions, and placing bids. Clients should follow the specified message formats and endpoints to ensure proper communication with the server.
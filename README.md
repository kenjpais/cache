# In-Memory Database/Cache

## Features

### CRUD Operations:
- **SET**: Add a key-value pair to the cache.
- **GET**: Retrieve a value by its key.
- **DEL**: Delete a key-value pair from the cache.

### Eviction Policies:
- **FIFO (First In, First Out)**: Items are evicted in the order they were inserted.  
  (Can be extended to add more policies like):
  - **LRU (Least Recently Used)**
  - **LFU (Least Frequently Used)**

### TTL Support:
- Cache entries can expire after a specified **Time-to-Live (TTL)**.

### Clear Cache:
- Clears all data from the cache, resetting both the data map and the queue.

## Design Patterns

### Singleton Pattern:
- Ensures a single instance of the cache throughout the application's lifecycle.

### Factory Pattern:
- Used for creating eviction policy structs based on the policy type.

### Strategy Pattern:
- Implements eviction policy logic based on the chosen policy type.

## Data Structures

### Doubly Linked List:
- Keeps track of the order in which items are inserted into the database.  
  Example:

  ```plaintext
  +-----------+      +-----------+      +-----------+
  |   key1    | <--> |   key2    | <--> |   key3    |
  +-----------+      +-----------+      +-----------+

### Hashmap:
- Stores reference to the node
